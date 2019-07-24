package bot

import (
	"log"

	"github.com/dootbin/dootdoot/config"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

//Tweet takes message as string and tweets string.
func Tweet(message string) string {

	var consumerKey = config.TwitterConsumerKey
	var consumerSecret = config.TwitterConsumerSecret
	var accessToken = config.TwitterAccessToken
	var accessSecret = config.TwitterAccessTokenSecret

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	//Update (POST!) Tweet (uncomment to run)
	tweet, _, _ := client.Statuses.Update(message, nil)

	var tweetURL = "https://twitter.com/24CarrotCraft/status/" + tweet.IDStr

	return tweetURL
}
