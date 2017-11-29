package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	//Token ... contains oath2 string
	Token string
	//BotPrefex ... cotains command prefex for commands.
	BotPrefix                string
	TwitterConsumerSecret    string
	TwitterConsumerKey       string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string

	config *configStruct
)

type configStruct struct {
	Token                    string `json:"Token"`
	BotPrefix                string `json:"BotPrefix"`
	TwitterConsumerSecret    string `json:"TwitterConsumerSecret"`
	TwitterConsumerKey       string `json:"TwitterConsumerKey"`
	TwitterAccessToken       string `json:"TwitterAccessToken"`
	TwitterAccessTokenSecret string `json:"TwitterAccessTokenSecret"`
}

//ReadConfig reads config.json file.
func ReadConfig() error {
	fmt.Println("Reading from the config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	TwitterConsumerKey = config.TwitterConsumerKey
	TwitterConsumerSecret = config.TwitterConsumerSecret
	TwitterAccessToken = config.TwitterAccessToken
	TwitterAccessTokenSecret = config.TwitterAccessTokenSecret
	return nil
}
