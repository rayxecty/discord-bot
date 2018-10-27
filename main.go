package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	
	"github.com/bwmarrin/discordgo"
)
var(
    stopBot = make(chan bool)
	vcsession *discordgo.VoiceConnection
	
	appConfig AppConfig
	RandomMessageList []string
)

// AppConfig 設定秘匿情報
type AppConfig struct {
	DiscordBotToken	string `json:"bot_token"`
	BotName			string `json:"bot_name"`
}

func init() {
	settingInit()
}

func settingInit() error {
	appConfigJSON, err := ioutil.ReadFile("./configs/config.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	json.Unmarshal(appConfigJSON, &appConfig)

	fmt.Println(appConfig)
	
	return nil
}

func main() {
	discord, err := discordgo.New()
	discord.Token = appConfig.DiscordBotToken
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	// websocket
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	<-stopBot

	return
}

// メッセージ作成
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	checkMessage(s, m, c)
}

