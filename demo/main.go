package main

import (
	"fmt"
	"github.com/TtMyth123/telegramLib"
	"github.com/TtMyth123/telegramLib/libType"
	"github.com/astaxie/beego"
	"github.com/zelenin/go-tdlib/client"
)

func main() {
	apiHash := beego.AppConfig.String("Telegram::apiHash")
	apiId64, _ := beego.AppConfig.Int64("Telegram::apiId")
	apiId := int32(apiId64)

	aTDClient := telegramLib.NewTDClient(apiId, apiHash, "+8613707720054", "", true)
	aTDClient.HandleListener.UpdateUserStatus = UpdateUserStatus
	aTDClient.HandleListener.UpdateNewMessageText = UpdateNewMessageText

	var phoneNumber string
	fmt.Scanln(&phoneNumber)

}

func UpdateUserStatus(data client.UpdateUserStatus) (map[string]interface{}, error) {
	fmt.Println(data)
	return nil, nil
}

func UpdateNewMessageText(data libType.Message) (map[string]interface{}, error) {
	fmt.Println(data)
	return nil, nil
}
