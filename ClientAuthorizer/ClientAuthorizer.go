package ClientAuthorizer

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

type TtClientAuthorizer struct {
	client.AuthorizationStateHandler
	TdlibParameters chan *client.SetTdlibParametersRequest
	PhoneNumber     chan string
	Code            chan string
	State           chan AuthorizationState
	Password        chan string
	isC             bool
	Ok              bool

	initPhoneNumber string
	initPassword    string
}

func ClientAuthorizer(PhoneNumber, Password string, isC bool) *TtClientAuthorizer {
	aTtClientAuthorizer := &TtClientAuthorizer{
		TdlibParameters: make(chan *client.SetTdlibParametersRequest, 1),
		PhoneNumber:     make(chan string, 1),
		Code:            make(chan string, 1),
		State:           make(chan AuthorizationState, 10),
		Password:        make(chan string, 1),
	}
	aTtClientAuthorizer.isC = isC
	aTtClientAuthorizer.initPassword = Password
	aTtClientAuthorizer.initPhoneNumber = PhoneNumber
	return aTtClientAuthorizer
}

type AuthorizationState interface {
	AuthorizationStateType() string
}

func (stateHandler *TtClientAuthorizer) Handle(client11 *client.Client, state client.AuthorizationState) error {
	stateHandler.State <- state

	switch state.AuthorizationStateType() {
	case client.TypeAuthorizationStateWaitTdlibParameters:
		_, err := client11.SetTdlibParameters(<-stateHandler.TdlibParameters)
		return err

	//case client.TypeAuthorizationStateWaitEncryptionKey:
	//	_, err := client11.CheckDatabaseEncryptionKey(&client.CheckDatabaseEncryptionKeyRequest{})
	//	return err

	case client.TypeAuthorizationStateWaitPhoneNumber:
		_, err := client11.SetAuthenticationPhoneNumber(&client.SetAuthenticationPhoneNumberRequest{
			PhoneNumber: <-stateHandler.PhoneNumber,
			Settings: &client.PhoneNumberAuthenticationSettings{
				AllowFlashCall:       false,
				IsCurrentPhoneNumber: false,
				AllowSmsRetrieverApi: false,
			},
		})
		return err

	case client.TypeAuthorizationStateWaitCode:
		_, err := client11.CheckAuthenticationCode(&client.CheckAuthenticationCodeRequest{
			Code: <-stateHandler.Code,
		})
		return err

	case client.TypeAuthorizationStateWaitRegistration:
		return client.ErrNotSupportedAuthorizationState

	case client.TypeAuthorizationStateWaitPassword:
		_, err := client11.CheckAuthenticationPassword(&client.CheckAuthenticationPasswordRequest{
			Password: <-stateHandler.Password,
		})
		return err

	case client.TypeAuthorizationStateReady:
		return nil

	case client.TypeAuthorizationStateLoggingOut:
		return client.ErrNotSupportedAuthorizationState

	case client.TypeAuthorizationStateClosing:
		return nil

	case client.TypeAuthorizationStateClosed:
		return nil
	}

	return client.ErrNotSupportedAuthorizationState
}

func (stateHandler *TtClientAuthorizer) Close() {
	close(stateHandler.TdlibParameters)
	close(stateHandler.State)
}

func (stateHandler *TtClientAuthorizer) CliInteractor() {
	for {
		select {
		case state, ok := <-stateHandler.State:
			if !ok {
				return
			}

			switch state.AuthorizationStateType() {
			case client.TypeAuthorizationStateWaitPhoneNumber:
				if stateHandler.initPhoneNumber != "" {
					stateHandler.PhoneNumber <- stateHandler.initPhoneNumber
					stateHandler.initPhoneNumber = ""
					break
				}
				if !stateHandler.isC {
					break
				}

				//fmt.Println("Enter phone number: ")
				fmt.Println("输入手机号(如:+8613700001234):")
				var phoneNumber string
				fmt.Scanln(&phoneNumber)
				stateHandler.PhoneNumber <- phoneNumber

			case client.TypeAuthorizationStateWaitCode:
				if !stateHandler.isC {
					break
				}

				//fmt.Println("Enter code: ")
				fmt.Println("输入验证码:")
				var code string
				fmt.Scanln(&code)
				stateHandler.Code <- code

			case client.TypeAuthorizationStateWaitPassword:
				if stateHandler.initPassword != "" {
					stateHandler.Password <- stateHandler.initPassword
					stateHandler.initPassword = ""
					break
				}
				if !stateHandler.isC {
					break
				}

				//fmt.Println("Enter password: ")
				fmt.Println("输入密码:")
				var password string
				fmt.Scanln(&password)

				stateHandler.Password <- password
			case client.TypeAuthorizationStateReady:
				return
			}
		}
	}
}
