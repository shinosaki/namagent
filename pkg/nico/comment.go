package nico

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/shinosaki/namagent/pkg/nico/ndgr"
	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/data"
	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/service/edge"
)

var (
	once            sync.Once
	fetchedSegments sync.Map
)

func commentHandler(
	viewUri *url.URL,
	at string,
	chatChan chan any,
	client *http.Client,
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			log.Println("NicoCommentHandler: receive interrupt")
			once.Do(func() {
				close(chatChan)
			})
			return

		default:
			if _, exists := fetchedSegments.Load(viewUri.String()); exists {
				log.Println("NDGR: this uri is exists")
				return
			}
			fetchedSegments.Store(viewUri.String(), struct{}{})

			// Set URL
			if at != "" {
				q := viewUri.Query()
				q.Set("at", at)
				viewUri.RawQuery = q.Encode()
			}

			// Connect to NDGR
			ch, err := ndgr.Reader(viewUri.String(), client, ctx)
			if err != nil {
				log.Println("NDGR read error:", err)
				time.Sleep(1 * time.Second) // いらないかも
				continue
			}

			// Receive message
			for chunk := range ch {
				switch e := chunk.(type) {
				// mpn.live.nicovideo.jp/api/view
				case *edge.ChunkedEntry:
					switch entry := e.Entry.(type) {
					case *edge.ChunkedEntry_Next:
						at = strconv.FormatInt(entry.Next.At, 10)

					case *edge.ChunkedEntry_Segment:
						if uri, err := url.Parse(entry.Segment.Uri); err == nil {
							go commentHandler(uri, "", chatChan, client, ctx)
						}
					}

				// mpn.live.nicovideo.jp/api/segment
				case *edge.ChunkedMessage:
					switch p := e.Payload.(type) {
					case *edge.ChunkedMessage_Message:
						if m, ok := p.Message.Data.(*data.NicoliveMessage_Chat); ok {
							chatChan <- *m.Chat
						}
					}
				}
			}
		}
	}
}
