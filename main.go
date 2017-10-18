package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/matteocontrini/locuspocusbot/tg"
	// "github.com/pelletier/go-toml"
)

var bot *tg.Bot

func main() {
	var err error
	bot, err = tg.NewBot("")

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account @%s", bot.Me.Username)

	updates := make(chan tg.Update, 100)
	bot.GetUpdates(updates, 10*time.Second)

	for update := range updates {
		if update.Message != nil {
			// Ignore non-text messages
			if update.Message.Text != "" {
				handleMessage(update.Message)
			}
		}
	}
}

func handleMessage(message *tg.Message) {
	log.Printf("<%d> %s", message.Chat.ID, message.Text)

	msg := tg.MessageRequest{
		ChatID: message.Chat.ID,
		Text:   isDed(message.Text),
	}

	bot.Send(&msg)
}

func isDed(name string) string {
	url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&rvsection=0&titles=%s&format=json", name)
	resp, _ := http.Get(url)

	defer resp.Body.Close()
	json, _ := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile("death_date[ ]+= {{(.+?)\\|([0-9]{4}\\|[0-9]{1,2}\\|[0-9]{1,2})\\|(.+?)}}")
	match := re.FindStringSubmatch(string(json))

	fmt.Println(match)

	return match[2]
}
