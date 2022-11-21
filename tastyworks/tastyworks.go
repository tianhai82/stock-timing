package tastyworks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/tianhai82/stock-timing/firebase"
	"github.com/tianhai82/stock-timing/model"
)

const (
	apiUrl = "https://api.tastyworks.com/%s"
)

func getCachedSessionToken() (string, error) {
	doc, err := firebase.FirestoreClient.Collection("tastyworks").Doc("token").Get(context.Background())
	if err != nil {
		return "", err
	}
	var token map[string]string
	err = doc.DataTo(&token)
	if err != nil {
		return "", err
	}
	return token["token"], nil
}

func validateSessionToken(token string) (bool, error) {
	url := fmt.Sprintf(apiUrl, "sessions/validate")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 201 {
		return true, nil
	}
	return false, nil
}

func login() (string, error) {
	login, ok := os.LookupEnv("tastyworks_login")
	if !ok {
		return "", errors.New("tastyworks_login not found")
	}
	password, ok := os.LookupEnv("tastyworks_password")
	if !ok {
		return "", errors.New("tastyworks_password not found")
	}
	url := fmt.Sprintf(apiUrl, "sessions")
	body := map[string]string{
		"login":    login,
		"password": password,
	}
	s, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(s))
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("fail to login: %s", resp.Status)
	}
	var out struct {
		Data struct {
			SessionToken string `json:"session-token"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return "", err
	}
	token := out.Data.SessionToken
	return token, nil
}

func saveToken(token string) error {
	o := map[string]string{
		"token": token,
	}
	_, err := firebase.FirestoreClient.Collection("tastyworks").Doc("token").Set(context.Background(), o)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveHistory(instrument model.InstrumentDisplayData, period int) ([]model.Candle, error) {
	return nil, errors.New("not implemented")
}
