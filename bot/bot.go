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

var userRobertJohnson = "262737106319835136"
var userTommyTheBold = "108613322227601408"

var channelUnmoderated = "345397705436168192"
var channelGeneral = "225589027460349953"
var channelScreenShots = "241404784425435148"
var channelTechTalk = "326561998227767299"
var channelSpawn = "343904691480166403"
var channelNether = "236466997158871040"
var channelGoodBoyPoints = "402105717932294144"


var guild24CarrotCraft = "195174072634572800"

//BotID ... is exported var
var BotID string
var goBot *discordgo.Session
var database *sql.DB
var userQuery *sql.Rows

var adminRole = "204660177067048960"
var seniorRole = "289847564583305216"
var modRole = "204281815819747329"
var memberRole = ""

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

	statement, _ := database.Prepare("INSERT INTO chatUser (userID, userXP, lastSpoke, userGrad) VALUES (?, ?, ?, ?)")
	statement.Exec(id, xp, ls, 0)

}

//updateUserGrad updates users graditude returns bool value
func updateUserGrad(id string, amount int) bool {

	if verifyDatabaseUser(id) {

		updateGrad, _ := database.Prepare("UPDATE `chatUser` SET `userGrad` = $amount WHERE _rowid_ = $id")
		updateGrad.Exec(amount, id)
		return true
	}

	return false

}

//updateUser function to update user.
func updateUser(id int, xp int, ls int) {

	statement, _ := database.Prepare("UPDATE `chatUser` SET `userXP`= $xp, lastSpoke = $ls WHERE _rowid_=$id")
	statement.Exec(xp, ls, id)

}

//subGrad subtracts Graditude from user. Returns as bool
func subGrad(id string, grad int, s *discordgo.Session) bool {
	if verifyRole(id, modRole, s) {
		balance := getGradBalance(id)
		if updateUserGrad(id, (balance - grad)) {

			return true

		}
	}

	return false
}

//addGrad adds Graditude to user. Returns as integer Graditude amount.
func addGrad(id string, grad int, s *discordgo.Session) bool {

	if verifyRole(id, modRole, s) {
		balance := getGradBalance(id)
		if updateUserGrad(id, (grad + balance)) {
			return true
		}
		return false
	}
	return false
}

//verifyRole will verify that user does have specified role and return bool value.
func verifyRole(user string, role string, s *discordgo.Session) bool {
	member, error := s.GuildMember(guild24CarrotCraft, user)
	if error != nil {
		println("no member")
		return false
	}
	for index := range member.Roles {
		if member.Roles[index] == role {
			return true
		}
	}
	return false
}

//verifyGradBalance
func verifyGradBalance(id string, subtract int) bool {

	if verifyDatabaseUser(id) {
		balance := getGradBalance(id)
		subtractedBalance := balance - subtract

		if subtractedBalance > 0 || subtractedBalance == 0 {
			return true
		}

	}
	return false
}

//convertMentionToID
func convertMentionToID(id string) string {

	t := strings.Replace(id, "<@", "", -1)
	a := strings.Replace(t, "!", "", -1)
	parsedID := strings.Replace(a, ">", "", -1)
	return parsedID
}

//verifyTarget verfies that a target is real and has permissions.
func verifyTarget(id string, role string, s *discordgo.Session) bool {

	parsedID := convertMentionToID(id)

	if verifyDatabaseUser(parsedID) && verifyRole(parsedID, role, s) {

		return true

	}

	return false
}

//verifyDatabaseUser verifies that user is in database.
func verifyDatabaseUser(id string) bool {
	var userID int
	convertedID, err := strconv.Atoi(id)

	if err != nil {

		return false

	}

	verifyUserQuery := database.QueryRow("SELECT userID FROM chatUser WHERE userID = ?", convertedID)

	verifyUserQuery.Scan(&userID)
	if userID == convertedID {

		return true
	}

	return false
}

//verifiedTime
func verifiedTime() bool {
	return false
}

//getGradBalance Returns balance as a integer.
func getGradBalance(id string) int {
	var balance int

	gradQuery := database.QueryRow("SELECT `userGrad` FROM `chatUser` WHERE `userID` = ?", id)
	gradQuery.Scan(&balance)

	return balance
}

//gradHandler event handler for all graditude events.
func gradHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {

		if verifyRole(m.Author.ID, modRole, s) {

			userCommand := strings.Split(m.Content, " ")

			if strings.ToLower(userCommand[0]) == "`gratitude" {

				//makesure that the full command is valid.

				if len(userCommand) == 3 {
					target := convertMentionToID(userCommand[1])
					//verify mention is valid role ect.
					if verifyTarget(target, modRole, s) && target != m.Author.ID {

						amount, error := strconv.Atoi(userCommand[2])

						if error != nil {

							s.ChannelMessageSend(channelGoodBoyPoints, "Invalid amount")

						} else {

							if amount > 0 {

								if verifyGradBalance(m.Author.ID, amount) {

									addGrad(target, amount, s)
									subGrad(m.Author.ID, amount, s)
									message := "<@!" + m.Author.ID + "> has given " + strconv.Itoa(amount) + " Gratitudes to <@!" + target + ">"

									s.ChannelMessageSend(m.ChannelID, message)
									s.ChannelMessageDelete(m.ChannelID, m.ID)
								} else {
									message := "You do not have enough grad to give this amount. Your balance is " + strconv.Itoa(getGradBalance(m.Author.ID))
									s.ChannelMessageSend(m.ChannelID, message)

								}

							} else {

								s.ChannelMessageSend(m.ChannelID, "you must enter in a valid amount.")

							}

						}

					}

				}

			}

			if strings.ToLower(userCommand[0]) == "`give" {
				println("Given")
				_ = addGrad(m.Author.ID, int(200), s)
			}

			if strings.ToLower(userCommand[0]) == "`balance" {
				var message string

				balance := getGradBalance(m.Author.ID)

				message = "<@!" + m.Author.ID + "> you have " + strconv.Itoa(balance) + " Gratitudes!"
				s.ChannelMessageSend(m.ChannelID, message)
			}

		} else {

			//s.ChannelMessageSend(m.ChannelID, "You do not have the permissions for this command.")

		}

	}

}

func expHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if verifyRole(m.Author.ID, modRole, s) {
		if strings.HasPrefix(m.Content, config.BotPrefix) {

			userCommand := strings.Split(m.Content, " ")
			if strings.ToLower(userCommand[0]) == "`good_boy_points" {

				userID, _ := strconv.Atoi(m.Author.ID)
				var userXP string
				var lastSpoke string

				var err = database.QueryRow("SELECT userID, userXP, lastSpoke FROM chatUser WHERE userID = ?", userID).Scan(&userID, &userXP, &lastSpoke)

				if err != nil {

					createUser(userID, 0, 0)
					uid := strconv.Itoa(userID)

					println("Created user " + uid)
					s.ChannelMessageSend(channelGoodBoyPoints, "Your good boy points are 0!")

				} else {

					message := "Your total good boy points are " + userXP + "!"
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

//Start connects bot to discord and starts stream.
func Start() {
	database, _ = sql.Open("sqlite3", "./dootdoot.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS chatUser (userID INTEGER PRIMARY KEY, userXP INTEGER, userGrad INTEGER, lastSpoke INTEGER)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS botStats (goodBoy INTEGER, badBoy INTEGER)")
	statement.Exec()

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
	goBot.AddHandler(gradHandler)
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
