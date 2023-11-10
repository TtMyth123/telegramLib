package telegramLib

import (
	"fmt"
	"github.com/TtMyth123/kit/stringKit"
	"github.com/TtMyth123/telegramLib/ClientAuthorizer"
	"github.com/TtMyth123/telegramLib/HandleListener"
	"github.com/TtMyth123/telegramLib/kit"
	"github.com/TtMyth123/telegramLib/libType"
	"github.com/zelenin/go-tdlib/client"
	"log"
	"path/filepath"
)

var MpI map[string]bool

func init() {
	MpI = make(map[string]bool)

	MpI["updateChatReadInbox"] = true
	MpI["updateDeleteMessages"] = true
	MpI["updateHavePendingNotifications"] = true
	MpI["updateChatReadOutbox"] = true
	MpI["updateChatLastMessage"] = true
	MpI["updateConnectionState"] = true
	MpI["updateChatReadInbox"] = true
	MpI["updateChatReadInbox"] = true
	MpI["updateChatReadInbox"] = true
	MpI["updateChatReadInbox"] = true
}

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
func (t *TDClient) GetClient() *client.Client {
	return t.tdlibClient
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
					if MpI[GetType] {
						continue
					}
					switch GetType {
					case client.TypeUpdateUserStatus:
						if t.HandleListener.UpdateUserStatus != nil {
							aInfo, ok := bb.(*client.UpdateUserStatus)
							if ok {
								t.HandleListener.UpdateUserStatus(*aInfo)
							}
							//aData := &client.UpdateUserStatus{}
							//strBb, _ := json.Marshal(bb)
							//ee := json.Unmarshal(strBb, aData)
							//if ee == nil {
							//	t.HandleListener.UpdateUserStatus(*aData)
							//}
						}

					case client.TypeUpdateNewMessage:
						aInfo, ok := bb.(*client.UpdateNewMessage)
						if ok {
							t.updateNewMessage(*aInfo)
						}
						//aData := &client.UpdateNewMessage{}
						//strBb, _ := json.Marshal(bb)
						//ee := json.Unmarshal(strBb, aData)
						//if ee == nil {
						//	t.updateNewMessage(*aData)
						//}

					case client.TypeUpdateChatTitle:
						aInfo, ok := bb.(*client.UpdateChatTitle)
						if ok {
							t.updateChatTitle(*aInfo)
						}
					case client.TypeUpdateNewChat:
						aInfo, ok := bb.(*client.UpdateNewChat)
						if ok {
							t.UpdateNewChat(*aInfo)
						}

						//aData := &client.UpdateChatTitle{}
						//strBb, _ := json.Marshal(bb)
						//ee := json.Unmarshal(strBb, aData)
						//if ee == nil {
						//	t.updateChatTitle(*aData)
						//}
					case client.TypeUpdateUser:
						aInfo, ok := bb.(*client.UpdateUser)
						if ok {
							t.updateUser(*aInfo)
						}

						//aData := &client.UpdateUser{}
						//strBb, _ := json.Marshal(bb)
						//ee := json.Unmarshal(strBb, aData)
						//if ee == nil {
						//	t.updateUser(*aData)
						//}

					case client.TypeUpdateBasicGroupFullInfo:
						aInfo, ok := bb.(*client.UpdateBasicGroupFullInfo)
						if ok {
							t.updateBasicGroupFullInfo(*aInfo)
						}
					case client.TypeUpdateSupergroupFullInfo:
						aInfo, ok := bb.(*client.UpdateSupergroupFullInfo)
						if ok {
							t.updateSupergroupFullInfo(*aInfo)
						}
					default:
						strBb := stringKit.GetJsonStr(bb)
						fmt.Println("GetType:", GetType, "data:", strBb)
					}
				}
			}
		}(aListener)

	}

}

func (t *TDClient) updateNewMessage(aMessage client.UpdateNewMessage) error {
	MessageSenderType := aMessage.Message.SenderId.MessageSenderType()
	MessageContentType := aMessage.Message.Content.MessageContentType()

	switch MessageContentType {
	case client.TypeMessageText:
		if MessageSenderType == client.TypeMessageSenderUser {
			aMessageSenderUser, ok := aMessage.Message.SenderId.(*client.MessageSenderUser)
			if !ok {
				return fmt.Errorf("类型不对")
			}
			aMessageText, ok := aMessage.Message.Content.(*client.MessageText)
			if !ok {
				return fmt.Errorf("类型不对")
			}

			if t.HandleListener.UpdateNewMessageText != nil {
				data := libType.Message{
					Id:           aMessage.Message.Id,
					ChatId:       aMessage.Message.ChatId,
					SenderUserId: aMessageSenderUser.UserId,
					Text:         kit.GetFormattedText(aMessageText.Text),
					ContentType:  MessageSenderType,
				}
				t.HandleListener.UpdateNewMessageText(data)
			}
		}
	case client.TypeMessageAnimation:
	case client.TypeMessageAudio:
	case client.TypeMessagePhoto:
	case client.TypeMessageVideo:
	case client.TypeMessageContact:

	}

	return nil
}
func (t *TDClient) updateBasicGroupFullInfo(data client.UpdateBasicGroupFullInfo) error {
	if t.HandleListener.UpdateBasicGroupFullInfo == nil {
		return nil
	}

	_, e := t.HandleListener.UpdateBasicGroupFullInfo(data)
	return e
}
func (t *TDClient) updateChatTitle(data client.UpdateChatTitle) error {
	if t.HandleListener.UpdateChatTitle == nil {
		return nil
	}

	_, e := t.HandleListener.UpdateChatTitle(data)
	return e
}
func (t *TDClient) updateNewChat(data client.UpdateNewChat) error {
	if t.HandleListener.UpdateNewChat == nil {
		return nil
	}
	_, e := t.HandleListener.UpdateNewChat(data)
	return e
}
func (t *TDClient) updateSupergroupFullInfo(data client.UpdateSupergroupFullInfo) error {
	if t.HandleListener.UpdateSupergroupFullInfo == nil {
		return nil
	}
	_, e := t.HandleListener.UpdateSupergroupFullInfo(data)
	return e
}

func (t *TDClient) updateUser(data client.UpdateUser) error {
	if t.HandleListener.UpdateUser == nil {
		return nil
	}

	_, e := t.HandleListener.UpdateUser(data)
	return e
}

func (t *TDClient) GetState() bool {
	return t.isOk
}

func (t *TDClient) SendMessage(ChatId, ReplyToMessageId int64, Text string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId
	aInputMessage1 := &client.InputMessageText{}
	aInputMessage1.Text = kit.FormattedText(Text)
	req.InputMessageContent = aInputMessage1
	_, e := t.tdlibClient.SendMessage(req)
	return e
}

func (t *TDClient) SendMessageContact(ChatId, ReplyToMessageId int64, PhoneNumber, FirstName string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId
	aInputMessage1 := &client.InputMessageContact{}
	aInputMessage1.Contact = &client.Contact{}
	aInputMessage1.Contact.FirstName = FirstName
	aInputMessage1.Contact.PhoneNumber = PhoneNumber
	req.InputMessageContent = aInputMessage1

	_, e := t.tdlibClient.SendMessage(req)
	return e
}
func (t *TDClient) SendMessagePhotoFileLocal(ChatId, ReplyToMessageId int64, Path, Caption string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId

	aInputFileLocal := client.InputFileLocal{}
	aInputFileLocal.Path = Path

	aInputMessage1 := &client.InputMessagePhoto{}
	aInputMessage1.Photo = &aInputFileLocal
	aInputMessage1.Caption = kit.FormattedText(Caption)

	req.InputMessageContent = aInputMessage1
	_, e := t.tdlibClient.SendMessage(req)
	return e
}

func (t *TDClient) SendMessageAudioFileLocal(ChatId, ReplyToMessageId int64, Path, Caption string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId

	aInputFileLocal := client.InputFileLocal{}
	aInputFileLocal.Path = Path

	aInputMessage1 := &client.InputMessageAudio{}
	aInputMessage1.Audio = &aInputFileLocal
	aInputMessage1.Caption = kit.FormattedText(Caption)

	req.InputMessageContent = aInputMessage1
	_, e := t.tdlibClient.SendMessage(req)
	return e
}

func (t *TDClient) SendMessageVideoFileLocal(ChatId, ReplyToMessageId int64, Path, Caption string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId

	aInputFileLocal := client.InputFileLocal{}
	aInputFileLocal.Path = Path

	aInputMessage1 := &client.InputMessageVideo{}
	aInputMessage1.Video = &aInputFileLocal
	aInputMessage1.Caption = kit.FormattedText(Caption)

	req.InputMessageContent = aInputMessage1
	_, e := t.tdlibClient.SendMessage(req)
	return e
}

func (t *TDClient) SendMessageAnimationFileLocal(ChatId, ReplyToMessageId int64, Path, Caption string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	req := &client.SendMessageRequest{}
	req.ChatId = ChatId
	req.ReplyToMessageId = ReplyToMessageId

	aInputFileLocal := client.InputFileLocal{}
	aInputFileLocal.Path = Path

	aInputMessage1 := &client.InputMessageAnimation{}
	aInputMessage1.Animation = &aInputFileLocal
	aInputMessage1.Caption = kit.FormattedText(Caption)

	req.InputMessageContent = aInputMessage1
	_, e := t.tdlibClient.SendMessage(req)
	return e
}

func (t *TDClient) SendMessageFileLocalEx(ChatId, ReplyToMessageId int64, Path, Caption string) error {
	if t.tdlibClient == nil {
		return fmt.Errorf("没成功初始化Lib")
	}
	fileType, strT := kit.GetFileType(Path)
	switch fileType {
	case kit.FT_mp3:
		return t.SendMessageAudioFileLocal(ChatId, ReplyToMessageId, Path, Caption)
	case kit.FT_mp4:
		return t.SendMessageVideoFileLocal(ChatId, ReplyToMessageId, Path, Caption)
	case kit.FT_gif:
		return t.SendMessageAnimationFileLocal(ChatId, ReplyToMessageId, Path, Caption)
	case kit.FT_png:
		return t.SendMessagePhotoFileLocal(ChatId, ReplyToMessageId, Path, Caption)

	}
	return fmt.Errorf("没有对应的类型：%s", strT)
}
