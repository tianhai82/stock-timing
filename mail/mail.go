package mail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Email struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
type Payload struct {
	HTMLContent string  `json:"htmlContent"`
	Sender      Email   `json:"sender"`
	ReplyTo     Email   `json:"replyTo"`
	TemplateID  int     `json:"templateId"`
	Params      gin.H   `json:"params"`
	To          []Email `json:"to"`
}
type TemplateID int

const (
	SURVEY_ACTIVATED TemplateID = 1 + iota
	REQUEST_FOR_ACTIVATION
	SURVEY_REJECTION
	SHARE_GENERAL_SURVEY
	SHARE_CORE_COMP_SURVEY
)

type sendInBlueResp struct {
	MessageId string `json:"messageId"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

const url = "https://api.sendinblue.com/v3/smtp/email"

func Sendmail(apiKey string, templateId TemplateID, params gin.H, recipients []Email) error {
	payload := Payload{
		HTMLContent: "-",
		Params:      params,
		Sender: Email{
			Name:  "Time to Trade",
			Email: "robohuat82@gmail.com",
		},
		To:         recipients,
		TemplateID: int(templateId),
		ReplyTo: Email{
			Name:  "Time to Trade",
			Email: "robohuat82@gmail.com",
		},
	}
	s, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Add("api-key", apiKey)
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var respBody sendInBlueResp
	err = dec.Decode(&respBody)
	if err != nil {
		return err
	}
	if respBody.MessageId == "" {
		return fmt.Errorf(respBody.Message)
	}
	return nil
}
