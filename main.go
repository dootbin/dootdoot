package main

import (
	"fmt"

	"./bot"
	"./config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return

}
