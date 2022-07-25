package HandleListener

import "github.com/zelenin/go-tdlib/client"

type UpdateUserStatusFunc func(data client.UpdateUserStatus)
type UpdateNewMessageFunc func(data client.UpdateNewMessage)
