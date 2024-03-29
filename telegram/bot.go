package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramMsg struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}
type Message struct {
	Datetime string
	Level    string
	Msg      string
}

var baseTelegramUrl string = "https://api.telegram.org/bot"

const TOKEN = "1743013035:AAF43wU6BX4UOcHwL-vX2OGcM1xMhBoe0Ug"
const CHAT_ID = "21450012"

func SendMessage(msg, chatId, token string) error {
	tMsg := TelegramMsg{
		chatId, msg, "HTML",
	}
	jsonValue, _ := json.Marshal(tMsg)
	resp, err := http.Post(baseTelegramUrl+token+"/sendmessage", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Unable to send telegram msg", msg)
		return err
	}
	var out map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&out)
	if out["ok"].(bool) == false {
		fmt.Printf("%v - %v", out["error_code"], out["description"])
		return fmt.Errorf("%v - %v", out["error_code"], out["description"])
	}
	return nil
}
