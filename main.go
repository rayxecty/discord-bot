package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	stopBot   = make(chan bool)
	vcsession *discordgo.VoiceConnection

	appConfig         AppConfig
	RandomMessageList []string
	logfilePath       = "./log/discord-bot.log"
)

// AppConfig 設定秘匿情報
type AppConfig struct {
	DiscordBotToken       string `json:"bot_token"`
	BotName               string `json:"bot_name"`
	SpreadsheetId         string `json:"spreadsheet_id"`
	SpreadsheetSecretPath string `json:"spreadsheet_secret_path"`
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

func httpClient(credentialFilePath string) (*http.Client, error) {
	data, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	return conf.Client(oauth2.NoContext), nil
}

func main() {
	f, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open " + logfilePath + ":" + err.Error())
	}
	defer f.Close()

	log.SetOutput(f)

	discord, err := discordgo.New()
	discord.Token = appConfig.DiscordBotToken
	if err != nil {
		log.Println("Error logging in")
		log.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	// websocket
	err = discord.Open()
	if err != nil {
		log.Println(err)
	}

	log.Println("Listening...")
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
