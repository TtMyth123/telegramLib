package main

import (
	"fmt"
	"github.com/TtMyth123/telegramLib"
	"github.com/astaxie/beego"
	"github.com/zelenin/go-tdlib/client"
)

func main() {
	apiHash := beego.AppConfig.String("Telegram::apiHash")
	apiId64, _ := beego.AppConfig.Int64("Telegram::apiId")
	apiId := int32(apiId64)

	aTDClient := telegramLib.NewTDClient(apiId, apiHash, "+8613707720054", "", true)
	aTDClient.HandleListener.UpdateUserStatus = UpdateUserStatus
	aTDClient.HandleListener.UpdateNewMessage = UpdateNewMessage

	var phoneNumber string
	fmt.Scanln(&phoneNumber)

}

func UpdateUserStatus(data client.UpdateUserStatus) {
	fmt.Println(data)
}

func UpdateNewMessage(data client.UpdateNewMessage) {
	fmt.Println(data.Message)
}
