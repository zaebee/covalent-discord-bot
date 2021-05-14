package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

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
		log.Println("error creating Discord session,", err)
		return
	}
	bot := Bot{dg, es}
	// Register the messageCreate func as a callback for MessageCreate events.
	bot.discordClient.AddHandler(bot.messageCreate)
	// Register the messageReactionAdd func as a callback for MessageReactionAdd events.
	bot.discordClient.AddHandler(bot.messageReactionAdd)
	// Open a websocket connection to Discord and begin listening.
	if err := dg.Open(); err != nil {
		log.Println("error opening connection,", err)
		return
	}
	log.Println("Bot is now running.  Press CTRL-C to exit.")
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
		log.Printf("New meme from author: %+v\n", m.Author)
		for _, a := range m.Attachments {
			created, err := m.Timestamp.Parse()
			if err != nil {
				log.Printf("unable convert timestamp %v: %v", m.Timestamp, err)
			}
			if err := b.saveMeme(m.ID, a.URL, author, created); err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			answer := fmt.Sprintf("Jaja! Your meme was saved `%v`\nCheck it https://cqt-memes.herokuapp.com/", a.URL)
			s.ChannelMessageSend(m.ChannelID, answer)
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
	log.Printf("saved user reaction %+v for meme %v: %v\n", m.UserID, m.MessageID, res)
}

func (b *Bot) saveMeme(msgID, url string, author Author, ts time.Time) error {
	meme := Meme{
		ID:        msgID,
		Url:       url,
		Author:    author,
		Timestamp: ts,
		Reactions: []string{},
	}
	if hash, _ := imageHash(url); hash != "" {
		meme.ImageHash = hash
		similar, err := b.duplicates(hash)
		if err != nil {
			log.Printf("request find duplicates failed: %v", err)
		}
		if len(similar) > 0 {
			return fmt.Errorf(":japanese_ogre: hey jojo! we already have had this meme :japanese_ogre: ")
		}
	}
	res, err := b.elasticClient.Index("memes", esutil.NewJSONReader(&meme), b.elasticClient.Index.WithDocumentID(msgID))
	log.Printf("saved meme %+v: %v\n", meme, res)
	return err
}

func (b *Bot) duplicates(hash string) ([]*SearchHit, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"ImageHash": hash,
			},
		},
	}
	res, err := b.elasticClient.Search(
		b.elasticClient.Search.WithIndex("memes"),
		b.elasticClient.Search.WithBody(esutil.NewJSONReader(&query)),
		b.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var sr SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return sr.Hits.Hits, nil
}

func imageHash(u string) (string, error) {
	apiUrl := fmt.Sprintf("%v/hash?", os.Getenv("ELASTICSEARCH_URL"))
	params := url.Values{}
	params.Add("url", u)
	resp, err := http.Get(apiUrl + params.Encode())
	if err != nil {
		log.Printf("request hash url failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	var out ImageHash
	json.NewDecoder(resp.Body).Decode(&out)
	return out.Hash, nil
}
