// Package main implements discord bot to store meme image into elasticsearch backend.
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
	"github.com/zaebee/covalent-discord-bot/bot/model"
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
	author := model.Author{
		ID:       m.Author.ID,
		Avatar:   m.Author.Avatar,
		Username: m.Author.Username,
	}
	if len(m.Attachments) > 0 {
		for _, a := range m.Attachments {
			created, err := m.Timestamp.Parse()
			if err != nil {
				log.Printf("unable convert timestamp %v: %v", m.Timestamp, err)
			}
			answer := fmt.Sprintf("Jaja! Your meme was saved `%v`\nCheck it https://cqt-memes.herokuapp.com/", a.URL)
			if err := b.saveMeme(m.ID, a.URL, author, created, []string{}); err != nil {
				answer = err.Error()
			}
			s.ChannelMessageSend(m.ChannelID, answer)
		}
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	// If the message is "!sync" collect memes from history."
	if m.Content == "!sync" {
		messages := b.readHistory(s, m.ChannelID, m.ID)
		memes := b.saveHistoryMemes(messages)
		res := fmt.Sprintf("History synced. saved %d memes", len(memes))
		s.ChannelMessageSend(m.ChannelID, res)
	}
	if m.Content == "!top" {
		top, err := b.topAuthors()
		res := "Uh.. We don't have memes :cry: Try to upload something?"
		if err == nil && top != "" {
			res = fmt.Sprintf("Top meme makers:```%s```", top)
		}
		s.ChannelMessageSend(m.ChannelID, res)
	}
}

// This function will be called (due to AddHandler above) every time a new
// reaction is created on any message that the authenticated bot has access to.
func (b *Bot) messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	data := fmt.Sprintf(`{
    "script": {
		    "source": "if (!ctx._source.Reactions.contains(params.reaction)) { ctx._source.Reactions.add(params.reaction); ctx._source.Rating += 1 }",
		    "lang": "painless",
		    "params": {
		      "reaction": "%s"
        }
      }
  }`, m.UserID)
	res, _ := b.elasticClient.Update(
		"memes",
		m.MessageID,
		strings.NewReader(data),
		b.elasticClient.Update.WithPretty())
	log.Printf("saved user reaction %+v for meme %v: %v\n", m.UserID, m.MessageID, res)
}

func (b *Bot) topAuthors() (string, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"byAuthor": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "Author.Username.keyword",
				},
			},
		},
	}
	res, err := b.elasticClient.Search(
		b.elasticClient.Search.WithIndex("memes"),
		b.elasticClient.Search.WithBody(esutil.NewJSONReader(&query)),
		b.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var sr model.SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return "nil", err
	}
	m, ok := sr.Aggregations.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("want type map[string]interface{};  got %T", sr.Aggregations)
	}
	var result []string
	if byAuthor, ok := m["byAuthor"]; ok {
		if m, ok := byAuthor.(map[string]interface{}); ok {
			if buckets, ok := m["buckets"]; ok {
				if authors, ok := buckets.([]interface{}); ok {
					for i, a := range authors {
						if m, ok := a.(map[string]interface{}); ok {
							username := m["key"]
							count := m["doc_count"]
							msg := fmt.Sprintf("%d @%v: (%v memes)", i+1, username, count)
							result = append(result, msg)
						}
					}
				}
			}
		}
	}
	return strings.Join(result[:], "\n"), nil
}

func (b *Bot) readHistory(s *discordgo.Session, channelID, msgID string) []*discordgo.Message {
	var result []*discordgo.Message
	var fetch func(beforeMsgID string)
	fetch = func(beforeMsgID string) {
		messages, err := s.ChannelMessages(channelID, 100, beforeMsgID, "", "")
		if err != nil {
			log.Printf("while read history: %v", err)
			return
		}
		if len(messages) == 0 {
			return
		}
		result = append(result, messages...)
		lastMsg := messages[len(messages)-1]
		fetch(lastMsg.ID)
	}
	fetch(msgID)
	return result
}

func (b *Bot) saveHistoryMemes(messages []*discordgo.Message) []model.Meme {
	var memes []model.Meme
	for _, m := range messages {
		author := model.Author{
			ID:       m.Author.ID,
			Avatar:   m.Author.Avatar,
			Username: m.Author.Username,
		}
		var reactions []string
		for _, r := range m.Reactions {
			reactions = append(reactions, r.Emoji.Name)
		}
		if len(m.Attachments) == 1 {
			attach := m.Attachments[0]
			created, _ := m.Timestamp.Parse()
			go b.saveMeme(m.ID, attach.URL, author, created, reactions)
			memes = append(memes, model.Meme{ID: m.ID})
		}
	}
	fmt.Printf("got %d memes from history\n", len(memes))
	return memes
}

func (b *Bot) saveMeme(msgID, url string, author model.Author, ts time.Time, reactions []string) error {
	meme := model.Meme{
		ID:        msgID,
		Url:       url,
		Author:    author,
		Timestamp: ts,
		Reactions: reactions,
		Rating:    len(reactions),
	}
	if hash := imageHash(url); hash != "" {
		meme.ImageHash = hash
		similar, err := b.duplicates(hash)
		if err != nil {
			log.Printf("request find duplicates failed: %v", err)
		}
		if len(similar) > 0 {
			msg := ":japanese_ogre: hey jojo! we already have had this meme :japanese_ogre: "
			log.Println(msg)
			return fmt.Errorf(msg)
		}
	}
	res, err := b.elasticClient.Index(
		"memes",
		esutil.NewJSONReader(&meme),
		b.elasticClient.Index.WithDocumentID(msgID),
		b.elasticClient.Index.WithRefresh("true"))
	log.Printf("saved meme %v: %v\n", meme.Url, res)
	return err
}

func (b *Bot) duplicates(hash string) ([]*model.SearchHit, error) {
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
	var sr model.SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return sr.Hits.Hits, nil
}

func imageHash(u string) string {
	apiUrl := fmt.Sprintf("%v/hash?", os.Getenv("ELASTICSEARCH_URL"))
	params := url.Values{}
	params.Add("url", u)
	resp, err := http.Get(apiUrl + params.Encode())
	if err != nil {
		log.Printf("request meme content for url %q failed: %v", u, err)
		return ""
	}
	defer resp.Body.Close()
	var out model.ImageHash
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ""
	}
	return out.Hash
}
