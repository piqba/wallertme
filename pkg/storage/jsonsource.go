package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// WalletsFromJsonToMAP read json filename and path
func WalletsFromJsonToMAP(path, filename string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Wallet struct {
	Address         string            `json:"address"`
	Symbol          string            `json:"symbol"`
	IsActive        bool              `json:"is_active"`
	NotifierService []NotifierService `json:"notifier_service"`
	NetworkType     string            `json:"network_type"`
}

type NotifierService struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

// WalletsFromJsonToStruct read json filename and path
func WalletsFromJsonToStruct(path, filename string) ([]Wallet, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []Wallet
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
