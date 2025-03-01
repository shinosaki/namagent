package nico

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/namagent"
	"github.com/shinosaki/namagent/pkg/nico/nicows"
)

var Plugin = &namagent.Plugin{
	Client:    Client,
	Session:   NewSession,
	ExtractId: ExtractId,
}

func Client(
	programId string,
	config *config.Config,
	client *http.Client,
	ctx context.Context,
) (chan any, chan namagent.StreamData, error) {
	var (
		chatChan          = make(chan any, 50)
		streamDataChan    = make(chan namagent.StreamData, 1)
		clientCtx, cancel = context.WithCancel(context.Background())
	)

	// Program Data
	program, err := FetchProgramData(programId, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch program data: %v", err)
	}

	log.Println("NicoClient: fetched program data", program.Program.NicoliveProgramId)

	// Connect to websocket server
	if program.Site.Relive.WebSocketUrl == "" {
		return nil, nil, fmt.Errorf("websocket url is empty")
	}

	ws, ch := nicows.NewClient(cancel)
	if err := ws.Connect(program.Site.Relive.WebSocketUrl, 3, 2, false); err != nil {
		return nil, nil, fmt.Errorf("failed to connect websocket server %v", err)
	}

	log.Println("NicoClient: successful connected websocket server")

	// Receive message forever
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("NicoClient: receive interrupt")
				ws.Disconnect(false)
				return

			case m, ok := <-ch:
				if !ok {
					log.Println("NicoClient: message channel is closed")
					return
				}
				log.Println("NicoClient: received message type is", m.Type)

				switch m.Type {
				case nicows.MESSAGE_SERVER:
					if data := utils.Unmarshaller[nicows.MessageServerData](m.Data); data != nil {
						if uri, err := url.Parse(data.ViewURI); err == nil {
							go commentHandler(uri, "now", chatChan, client, clientCtx)
						}
					}

				case nicows.STREAM:
					if data := utils.Unmarshaller[nicows.StreamData](m.Data); data != nil {
						var cookies []*http.Cookie
						for _, c := range data.Cookies {
							cookies = append(cookies, &http.Cookie{
								Name:   c.Name,
								Value:  c.Value,
								Domain: c.Domain,
								Path:   c.Path,
							})
						}

						streamDataChan <- namagent.StreamData{
							URL:     data.URI,
							Cookies: cookies,
							Template: namagent.Template{
								AuthorId:   program.Program.Supplier.ProgramProviderId,
								AuthorName: program.Program.Supplier.Name,

								ProgramId:    program.Program.NicoliveProgramId,
								ProgramTitle: program.Program.Title,

								CreatedAt:  program.Program.OpenTime.Time,
								StartedAt:  program.Program.BeginTime.Time,
								FinishedAt: program.Program.EndTime.Time,
							},
						}

						close(streamDataChan)
					}
				}
			}
		}
	}()

	return chatChan, streamDataChan, nil
}
