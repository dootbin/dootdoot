package bot

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"../config"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

//BotID ... is exported var
var BotID string
var goBot *discordgo.Session

var database *sql.DB
var userQuery *sql.Rows

var adminRole = "204660177067048960"
var seniorRole = "289847564583305216"
var modRole = "204281815819747329"

//var database, _ = sql.Open("sqlite3", "./dootdoot.db")
//var statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS chatUser (id INTEGER PRIMARY KEY, userID INTEGER, userXP INTEGER, lastSpoke INTEGER)")
var helpMenu = []string{
	"`help                - Displays command list to user",
	"`tweet24cc [message] - tweets message from 24Carrotcraft Twitter account.",
	"`setstatus [message] - Sets message to send when pinged from other user",
	"`adamsandlermovies   - displays all adam sandler movies",
}

func returnUserEXP(userID int) (int, int, int, error) {

	userQuery, err := database.Query("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID = '?'", userID)
	var userXP int
	var lastSpoke int

	userQuery.Scan(&userID, &userXP, &lastSpoke)

	return userID, userXP, lastSpoke, err
}

func returnBotStats() (int, int) {

	var goodBoy int
	var badBoy int
	botQuery, err := database.Query("SELECT goodBoy, badBoy FROM botStats")

	if err != nil {
		goodBoy = 0
		badBoy = 0
		return goodBoy, badBoy
	}
	botQuery.Scan(&goodBoy, &badBoy)
	return goodBoy, badBoy
}

//createUser creates a new user
func createUser(id int, xp int, ls int) {

	statement, _ := database.Prepare("INSERT INTO chatUser (userID, userXP, lastSpoke) VALUES (?, ?, ?)")
	statement.Exec(id, xp, ls)

}

//updateUser function to update user.
func updateUser(id int, xp int, ls int) {

	statement, _ := database.Prepare("UPDATE `chatUser` SET `userXP`= $xp, lastSpoke = $ls WHERE _rowid_=$id")
	statement.Exec(xp, ls, id)

}

//verifiedTime
func verifiedTime() {

}

//Start connects bot to discord and starts stream.
func Start() {
	database, _ = sql.Open("sqlite3", "./dootdoot.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS chatUser (userID INTEGER PRIMARY KEY, userXP INTEGER, lastSpoke INTEGER)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS botStats (goodBoy INTEGER, badBoy INTEGER)")
	statement.Exec()
	//updateUser(179255730937790464, 21, 120)
	TestDB()

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
	goBot.AddHandler(expHandler)

	//Open Connection.
	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		goBot.ChannelMessageSend("225589027460349953", scanner.Text())

		if scanner.Text() == ".quit" {

			break

		}

	}

}

func expHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	member, error := s.GuildMember("195174072634572800", m.Author.ID)
	if error != nil {
		println("errorretrieveiasd")
	} else {
		for index := range member.Roles {

			if member.Roles[index] == modRole {
				if strings.HasPrefix(m.Content, config.BotPrefix) {

					userCommand := strings.Split(m.Content, " ")
					if strings.ToLower(userCommand[0]) == "`modpoints" {

						userID, _ := strconv.Atoi(m.Author.ID)
						var userXP string
						var lastSpoke string

						var err = database.QueryRow("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID = ?", userID).Scan(&userID, &userXP, &lastSpoke)

						if err != nil {

							createUser(userID, 0, 0)
							uid := strconv.Itoa(userID)

							println("Created user " + uid)

						} else {

							message := "Your total modpoints are " + userXP + "!"
							_, _ = s.ChannelMessageSend("402105717932294144", message)

						}

					}

				}

				if m.ChannelID == "225589027460349953" {

					userID, _ := strconv.Atoi(m.Author.ID)
					var userXP string
					var lastSpoke string
					var err = database.QueryRow("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID = ?", userID).Scan(&userID, &userXP, &lastSpoke)

					now := time.Now()
					secs := now.Unix()

					if err != nil {

						createUser(userID, 0, 0)
						uid := strconv.Itoa(userID)

						println("Created user " + uid)

					} else {

						var lastTime int

						lastTime, _ = strconv.Atoi(lastSpoke)

						if int(secs) > (lastTime + 60) {

							var newXP int
							oldXP, _ := strconv.Atoi(userXP)
							newXP = oldXP + 10
							updateUser(userID, newXP, int(secs))
							txp := strconv.Itoa(newXP)
							println("Added 10 XP total " + txp)

						}

					}

				}

			}

		}

	}

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	c := color.New(color.FgYellow)

	if strings.HasPrefix(m.Content, config.BotPrefix) {

		userCommand := strings.Split(m.Content, " ")

		if m.Author.ID == BotID {
			return
		}

		if strings.ToLower(userCommand[0]) == "`help" {

			message := "```--Here is a list of my commands!--\n"

			for index := range helpMenu {

				message += helpMenu[index] + "\n"

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
				tweetURL := Tweet(message)
				_, _ = s.ChannelMessageSend(m.ChannelID, tweetURL)
				_, _ = s.ChannelMessageSend("225589027460349953", tweetURL)

			} else {
				_, _ = s.ChannelMessageSend(m.ChannelID, "you do not have permission")
			}

		}

		if strings.ToLower(userCommand[0]) == "`exp" {

			//CheckPermissionlevelmod

		}

	}

	fmt.Print(m.Author.Username)
	c.Print(" > ")
	fmt.Println(m.Content)

}

//TestDB tests the database
func TestDB() {

	var userID string
	var userXP string
	var lastSpoke string
	var err = database.QueryRow("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID = ?", "179255730937790464").Scan(&userID, &userXP, &lastSpoke)
	if err != nil {

	}
	fmt.Println(userID, userXP, lastSpoke)

	// rows, err := database.QueryRow("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID == ?", "179255730937790464").scan(&id)
	// var userID string
	// var userXP string
	// var lastSpoke string
	// if err == nil {
	// 	for rows.Next() {
	// 		rows.Scan(&userID, &userXP, &lastSpoke)
	// 		fmt.Println(userID + ": " + userXP + " " + lastSpoke + " end")
	// 	}

	// }

}
