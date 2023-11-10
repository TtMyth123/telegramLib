package HandleListener

import (
	"github.com/TtMyth123/telegramLib/libType"
	"github.com/zelenin/go-tdlib/client"
)

type UpdateUserStatusFunc func(data client.UpdateUserStatus) (map[string]interface{}, error)
type UpdateNewMessageFunc func(data client.UpdateNewMessage) (map[string]interface{}, error)
type UpdateNewMessageTextFunc func(data libType.Message) (map[string]interface{}, error)
type UpdateUserFunc func(data client.UpdateUser) (map[string]interface{}, error)
type UpdateChatTitleFunc func(data client.UpdateChatTitle) (map[string]interface{}, error)
type UpdateBasicGroupFullInfoFunc func(data client.UpdateBasicGroupFullInfo) (map[string]interface{}, error)

type UpdateNewChatFunc func(data client.UpdateNewChat) (map[string]interface{}, error)
type UpdateSupergroupFullInfoFunc func(data client.UpdateSupergroupFullInfo) (map[string]interface{}, error)
