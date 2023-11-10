package kit

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

func GetMessageText(aMessage client.UpdateNewMessage) (*client.MessageSenderUser, *client.MessageText, error) {
	aMessageSenderUser, ok := aMessage.Message.SenderId.(*client.MessageSenderUser)
	if !ok {
		return nil, nil, fmt.Errorf("类型不对")
	}
	aMessageText, ok := aMessage.Message.Content.(*client.MessageText)
	if !ok {
		return nil, nil, fmt.Errorf("类型不对")
	}

	return aMessageSenderUser, aMessageText, nil
}
