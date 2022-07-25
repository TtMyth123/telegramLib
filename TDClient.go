package telegramLib

import (
	"encoding/json"
	"github.com/TtMyth123/telegramLib/ClientAuthorizer"
	"github.com/TtMyth123/telegramLib/HandleListener"
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

type TDClient struct {
	tdlibClient *client.Client
	isOk        bool

	HandleListener.HandleListener
}

func NewTDClient(apiId int32, apiHash, phoneNumber, Password string, isC bool) *TDClient {
	aTDClient := &TDClient{}
	aTDClient.init(apiId, apiHash, phoneNumber, Password, isC)
	return aTDClient
}

func (t *TDClient) init(apiId int32, apiHash, phoneNumber, Password string, isC bool) {
	authorizer := ClientAuthorizer.ClientAuthorizer(phoneNumber, Password, isC)
	go authorizer.CliInteractor()

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})

	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}
	t.isOk = false
	go func(t *TDClient) {
		t.tdlibClient, err = client.NewClient(authorizer)
		if err != nil {
			log.Fatalf("NewClient error: %s", err)
		} else {
			t.isOk = true

			optionValue, err := t.tdlibClient.GetOption(&client.GetOptionRequest{
				Name: "version",
			})
			if err != nil {
				log.Fatalf("GetOption error: %s", err)
			}
			log.Printf("TDLib version: %s", optionValue.(*client.OptionValueString).Value)
			t.goListener()
		}
	}(t)
}
func (t *TDClient) goListener() {
	aListener := t.tdlibClient.GetListener()
	if aListener.IsActive() {
		go func(aa *client.Listener) {
			for {
				select {
				case bb := <-aa.Updates:
					GetType := bb.GetType()
					switch GetType {
					case client.TypeUpdateUserStatus:
						if t.HandleListener.UpdateUserStatus != nil {
							aData := &client.UpdateUserStatus{}
							strBb, _ := json.Marshal(bb)
							ee := json.Unmarshal(strBb, aData)
							if ee == nil {
								t.HandleListener.UpdateUserStatus(*aData)
							}
						}

					case client.TypeUpdateNewMessage:
						if t.HandleListener.UpdateNewMessage != nil {
							aData := &client.UpdateNewMessage{}
							strBb, _ := json.Marshal(bb)
							ee := json.Unmarshal(strBb, aData)
							if ee == nil {
								t.HandleListener.UpdateNewMessage(*aData)
							}
						}
					}
				}

			}
		}(aListener)

	}

}

func (t *TDClient) GetState() bool {
	return t.isOk
}
