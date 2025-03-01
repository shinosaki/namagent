package nico

import (
	"context"
	"crypto/ecdh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/shinosaki/namagent/internal/config"
	"github.com/shinosaki/namagent/internal/recorder"
	"github.com/shinosaki/namagent/internal/utils"
	"github.com/shinosaki/namagent/pkg/nico"
	"github.com/shinosaki/webpush-client-go/autopush"
	"github.com/shinosaki/webpush-client-go/rfc8291"
	"github.com/shinosaki/webpush-client-go/sites/nicopush"
	"github.com/spf13/viper"
)

func Alert(
	config *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
) error {
	var (
		wg              = &sync.WaitGroup{}
		activeRecording = &sync.Map{}
		client          = nico.NewSession(config)
	)

	recording := func(programId string, userId string) {
		// If not followed
		if !slices.Contains(config.Following.Nico, userId) {
			return
		}

		// If exists recorded
		if _, exists := activeRecording.Load(programId); exists {
			return
		}

		log.Println("NicoAlert: detected live stream", userId)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer activeRecording.Delete(programId)
			defer log.Println("recorder done")

			activeRecording.Store(programId, struct{}{})

			recorderCtx, recorderCancel := context.WithCancel(context.Background())
			go func() {
				<-ctx.Done()
				recorderCancel()
				log.Println("===recorderCancel() called")
			}()

			if err := recorder.Recorder(programId, config, recorderCtx, recorderCancel); err != nil {
				log.Println("NicoAlert: recorder error", err)
			}
		}()
	}

	// if login
	if nico.IsLogin(client) {
		// nicolive webpush
		log.Println("NicoAlert: try to receive nicopush")
		go func() {
			np, ch, err := nicoPushClient(activeRecording, config, client)
			if err != nil {
				log.Println("NicoPushClient Error:", err)
				return
			}

			for {
				select {
				case <-ctx.Done():
					log.Println("NicoPushClient: receive interrupt")
					return

				case notification, ok := <-ch:
					if !ok {
						return
					}

					payload, err := np.Decrypt(notification)
					if err != nil {
						log.Println("NicoPushClient Error: decrypt error", err)
						continue
					}

					var data nicopush.PushData
					if err := json.Unmarshal(payload.Data, &data); err != nil {
						log.Println("NicoPushClient Error: unmarshal error", err)
						continue
					}

					programId := nico.ExtractId(data.OnClick)
					if programId != "" {
						if program, err := nico.FetchProgramData(programId, client); err == nil {
							recording(program.Program.NicoliveProgramId, program.Program.Supplier.ProgramProviderId)
						}
					}
				}
			}
		}()
	}

	check := func(data []RecentProgram) {
		for _, program := range data {
			recording(program.Id, program.ProgramProvider.Id)
		}
	}

	fetch := func(isBulkFetch bool) {
		data, err := FetchRecentPrograms(isBulkFetch, client)
		if err != nil {
			log.Println("NicoAlert: fetch recent programs error", err)
			return
		}
		log.Printf("NicoAlert: fetch %d programs", len(data))
		check(data)
	}

	// in first time
	fetch(true)

	// monitoring
	lowerLimit := time.Duration(10 * time.Second)
	interval := config.Alert.CheckInterval
	if interval < lowerLimit {
		interval = lowerLimit
	}

	jitter := interval / 5 // 20% of interval
	ticker := utils.NewJitterTicker(interval, jitter)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("NicoAlert: receive interrupt")
			wg.Wait()
			return nil

		case <-ticker.C:
			fetch(false)
		}
	}
}

func nicoPushClient(
	activeRecording *sync.Map,
	config *config.Config,
	client *http.Client,
) (np *nicopush.NicoPushClient, ch chan autopush.Notification, err error) {
	curve := ecdh.P256()
	authSecret, _, privateKey := rfc8291.NewSecrets(curve) // generate webpush secrets

	// authsecret load from config
	if config.WebPush.NicoPush.AuthSecret != "" {
		authSecret, err = base64.RawURLEncoding.DecodeString(config.WebPush.NicoPush.AuthSecret)
		if err != nil {
			return nil, nil, fmt.Errorf("auth_secret base64 decode error %v", err)
		}
	}

	// privatekey load from config
	if config.WebPush.NicoPush.PrivateKey != "" {
		privateKeyBytes, err := base64.RawURLEncoding.DecodeString(config.WebPush.NicoPush.PrivateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("private_key base64 decode error %v", err)
		}
		privateKey, err = curve.NewPrivateKey(privateKeyBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("private_key ecdh load error %v", err)
		}
	}

	np, ch, err = nicopush.NewNicoPushClient(
		config.WebPush.NicoPush.UAID,
		config.WebPush.NicoPush.ChannelIDs,
		authSecret,
		privateKey,
		client,
	)
	if err != nil {
		return nil, nil, err
	}

	uaid, err := np.Handshake()
	if err != nil {
		return nil, nil, err
	}

	// 取得したUAIDを保存
	config.WebPush.NicoPush.UAID = uaid
	viper.WriteConfig()

	if len(config.WebPush.NicoPush.ChannelIDs) == 0 {
		channelID, err := np.Register()
		if err != nil {
			return nil, nil, err
		}

		// 取得したChannelIDを保存
		config.WebPush.NicoPush.ChannelIDs = append(config.WebPush.NicoPush.ChannelIDs, channelID)
		viper.WriteConfig()
	}

	return np, ch, nil
}
