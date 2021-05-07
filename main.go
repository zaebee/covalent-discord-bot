package main

import (
	"fmt"
	"os"

	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/gateway"
)

var client *gateway.Session

func main() {
	fmt.Println("Starting...")

	client = goscord.New(&gateway.Options{
		Token: os.Getenv("TOKEN"),
	})

	client.On("ready", func() {
		fmt.Println("Logged in as " + client.Me().Tag())
	})

	client.On("message", func(msg *discord.Message) {
		fmt.Printf("message author: %+v\n", msg.Author)
		if len(msg.Attachments) > 0 {
			var urls []string
			for _, a := range msg.Attachments {
				urls = append(urls, a.URL)
			}
			fmt.Printf("attaches urls: %+v\n", urls)
			answer := fmt.Sprintf("Jaja! I saved your meme %v", urls)
			client.Channel.Send(msg.ChannelId, answer)
		}
		if msg.Content == "ping" {
			client.Channel.Send(msg.ChannelId, "Pong ! 🏓")
		}
	})

	client.Login()

	select {}
}
