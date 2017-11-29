package bot

import (
	"fmt"
	"strings"

	"../config"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

//BotID ... is exported var
var BotID string
var goBot *discordgo.Session

var help_menu = []string{
	"`help                - Displays command list to user",
	"`tweet24cc [message] - tweets message from 24Carrotcraft Twitter account.",
	"`setstatus [message] - Sets message to send when pinged from other user",
	"`adamsandlermovies   - displays all adam sandler movies",
}

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	//Adds Message Handler
	goBot.AddHandler(messageHandler)

	//Open Connection.
	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := color.New(color.FgYellow)

	if strings.HasPrefix(m.Content, config.BotPrefix) {

		userCommand := strings.Split(m.Content, " ")

		if m.Author.ID == BotID {
			return
		}

		if strings.ToLower(userCommand[0]) == "`help" {

			var message = "```--Here is a list of my commands!--\n"

			for index := range help_menu {

				message += help_menu[index] + "\n"

			}

			message += " ```"
			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		}

		if strings.ToLower(userCommand[0]) == "`adamsandlermovies" {

			_, _ = s.ChannelMessageSend(m.ChannelID, "robocop")

		}

		if strings.ToLower(userCommand[0]) == "`tweet24cc" {

			if m.Author.ID == "151081244824698880" || m.Author.ID == "179255730937790464" {

				var message string
				for i := 1; i < len(userCommand); i++ {
					message += userCommand[i] + " "
				}
				var tweetURL = Tweet(message)
				_, _ = s.ChannelMessageSend(m.ChannelID, tweetURL)
				_, _ = s.ChannelMessageSend("225589027460349953", tweetURL)

			} else {
				_, _ = s.ChannelMessageSend(m.ChannelID, "you do not have permission")
			}

		}
	}

	/*
	   if m.Author.ID == BotID {
	     return
	   }

	   if m.Content == "hi" {
	     var message = fmt.Sprintf("Hello %s!", m.Author.Mention())
	     _, _ = s.ChannelMessageSend(m.ChannelID, message)
	   }

	   if m.Content == "bye" {
	     var message = fmt.Sprintf("Good bye %s!", m.Author.Mention())
	     _, _ = s.ChannelMessageSend(m.ChannelID, message)
	   }

	   if m.Content == "Good Morning!" {

	     var message = fmt.Sprintf("Good Morning %s!", m.Author.Mention())
	     _, _ = s.ChannelMessageSend(m.ChannelID, message)

	   }
	*/
	fmt.Print(m.Author.Username)
	c.Print(" > ")
	fmt.Println(m.Content)

}
