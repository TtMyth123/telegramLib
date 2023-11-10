package HandleListener

type HandleListener struct {
	UpdateUserStatus UpdateUserStatusFunc
	//UpdateNewMessage UpdateNewMessageFunc

	UpdateNewMessageText UpdateNewMessageTextFunc

	UpdateUser UpdateUserFunc

	UpdateChatTitle UpdateChatTitleFunc

	UpdateBasicGroupFullInfo UpdateBasicGroupFullInfoFunc
	UpdateNewChat            UpdateNewChatFunc
	UpdateSupergroupFullInfo UpdateSupergroupFullInfoFunc
}
