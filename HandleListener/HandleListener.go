package HandleListener

type HandleListener struct {
	UpdateUserStatus UpdateUserStatusFunc
	UpdateNewMessage UpdateNewMessageFunc
}
