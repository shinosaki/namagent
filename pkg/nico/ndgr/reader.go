package ndgr

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shinosaki/nicolive-comment-protobuf/proto/dwango/nicolive/chat/service/edge"
	"google.golang.org/protobuf/proto"
)

// require http2 client
func Reader(url string, client *http.Client, ctx context.Context) (<-chan any, error) {
	log.Println("NDGR get from", url)
	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch ndgr server failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("ndgr http error: %d %s", res.StatusCode, res.Status)
	}

	ch := make(chan any)

	go func() {
		defer res.Body.Close()
		defer close(ch)

		reader := NewProtobufStreamReader()
		buffer := make([]byte, 8192)

		for {
			select {
			case <-ctx.Done():
				log.Println("ndgr received interrupt")
				return
			default:
				n, err := res.Body.Read(buffer)
				if err != nil {
					if err != io.EOF {
						log.Println("ndgr read error:", err)
					}
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
