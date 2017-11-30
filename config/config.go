package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	//Token ... contains oath2 string
	Token string
	//BotPrefix cotains command prefex for commands.
	BotPrefix string
	//TwitterConsumerSecret ... contains ConsumerSecret token
	TwitterConsumerSecret string
	//TwitterConsumerKey contains ConsumerKey token
	TwitterConsumerKey string
	//TwitterAccessToken contains Twitter Access Token
	TwitterAccessToken string
	//TwitterAccessTokenSecret Contains Twitter Access Token Secret
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
