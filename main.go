package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type Author struct {
	ID       string
	Username string
	Avatar   string
}

type Meme struct {
	ID        string
	Url       string
	Author    Author
	Timestamp time.Time
	Reactions []string
}

type Bot struct {
	discordClient *discordgo.Session
	elasticClient *elasticsearch.Client
}

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	})
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	bot := Bot{dg, es}
	// Register the messageCreate func as a callback for MessageCreate events.
	bot.discordClient.AddHandler(bot.messageCreate)
	// Register the messageReactionAdd func as a callback for MessageReactionAdd events.
	bot.discordClient.AddHandler(bot.messageReactionAdd)
	// Open a websocket connection to Discord and begin listening.
	if err := dg.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	author := Author{
		ID:       m.Author.ID,
		Username: m.Author.Username,
		Avatar:   m.Author.Avatar,
	}
	if len(m.Attachments) > 0 {
		fmt.Printf("New meme from author: %+v\n", m.Author)
		for _, a := range m.Attachments {
			created, err := m.Timestamp.Parse()
			if err != nil {
				fmt.Printf("unable convert timestamp %v: %v", m.Timestamp, err)
			}
			if err := b.saveMeme(m.ID, a.URL, author, created); err == nil {
				answer := fmt.Sprintf("Jaja! Your meme was saved `%v`\nCheck it https://cqt-memes.herokuapp.com/", a.URL)
				s.ChannelMessageSend(m.ChannelID, answer)
			}
		}
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

// This function will be called (due to AddHandler above) every time a new
// reaction is created on any message that the authenticated bot has access to.
func (b *Bot) messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	data := fmt.Sprintf(`{
    "script": {
		    "source": "if (!ctx._source.Reactions.contains(params.reaction)) { ctx._source.Reactions.add(params.reaction) }",
		    "lang": "painless",
		    "params": {
		      "reaction": "%s"
        }
      }
  }`, m.UserID)
	res, _ := b.elasticClient.Update("memes", m.MessageID, strings.NewReader(data), b.elasticClient.Update.WithPretty())
	fmt.Printf("saved user reaction %+v for meme %v: %v\n", m.UserID, m.MessageID, res)
}

func (b *Bot) saveMeme(msgID, url string, author Author, ts time.Time) error {
	meme := Meme{
		ID:        msgID,
		Url:       url,
		Author:    author,
		Timestamp: ts,
		Reactions: []string{},
	}
	res, err := b.elasticClient.Index("memes", esutil.NewJSONReader(&meme), b.elasticClient.Index.WithDocumentID(msgID))
	fmt.Printf("saved meme %+v: %v\n", meme, res)
	return err
}
