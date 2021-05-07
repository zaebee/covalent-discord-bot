package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/gateway"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var client *gateway.Session

type Author struct {
	ID       string
	Username string
	Avatar   string
}

type Meme struct {
	Url       string
	Author    Author
	Timestamp time.Time
}

func main() {
	fmt.Println("Starting...")
	esHost := os.Getenv("ELASTICSEARCH_URL")
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esHost},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}
	client = goscord.New(&gateway.Options{
		Token: os.Getenv("TOKEN"),
	})

	client.On("ready", func() {
		fmt.Println("Logged in as " + client.Me().Tag())
	})
	client.On("message", func(msg *discord.Message) {
		if len(msg.Attachments) > 0 {
			author := Author{
				ID:       msg.Author.Id,
				Username: msg.Author.Username,
				Avatar:   msg.Author.Avatar,
			}
			fmt.Printf("New meme from author: %+v\n", author)
			for _, a := range msg.Attachments {
				meme := Meme{
					Url:       a.URL,
					Author:    author,
					Timestamp: msg.Timestamp,
				}
				res, _ := es.Index("memes", esutil.NewJSONReader(&meme))
				fmt.Println(res)
				answer := fmt.Sprintf("Jaja! Your meme was saved `%v\nCheck it https://cqt-memes.herokuapp.com/`", a.URL)
				client.Channel.Send(msg.ChannelId, answer)
			}
		}
		if msg.Content == "ping" {
			client.Channel.Send(msg.ChannelId, "Pong ! 🏓")
		}
	})

	client.Login()

	select {}
}
