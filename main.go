package main

import (
	"fmt"

	"github.com/dootbin/dootdoot/bot"
	"github.com/dootbin/dootdoot/config"
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
