package ndgr

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/service/edge"
	"google.golang.org/protobuf/proto"
)

func Reader(url string, client *http.Client, sc *utils.SignalContext) (<-chan interface{}, error) {
	log.Println("NDGR Reader: reading from", url)

	if client == nil {
		client = utils.NewHttp2Client()
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("NDGR Reader: failed to fetch %v", err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("NDGR Reader: invalid http status %d %s", res.StatusCode, res.Status)
	}

	ch := make(chan interface{})

	go func() {
		defer res.Body.Close()
		defer close(ch)

		reader := NewProtobufStreamReader()
		buffer := make([]byte, 4096)

		for {
			select {
			case <-sc.Context().Done():
				log.Println("NDGR Reader: receive interrupt...")
				return
			default:
				n, err := res.Body.Read(buffer)
				if err != nil {
					if err == io.EOF {
						return
					}
					log.Println("NDGR Reader: read error", err)
					return
				}

				reader.AddNewChunk(buffer[:n])
				for {
					chunk, ok := reader.UnshiftChunk()
					if !ok {
						break
					}

					// If ChunkedEntry
					entry := &edge.ChunkedEntry{}
					if err := proto.Unmarshal(chunk, entry); err == nil {
						ch <- entry
						continue
					}

					// If ChunkedMessage
					message := &edge.ChunkedMessage{}
					if err := proto.Unmarshal(chunk, message); err == nil {
						ch <- message
						continue
					}
				}

			}
		}
	}()

	return ch, nil
}
