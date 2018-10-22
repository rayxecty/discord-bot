package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	
	"github.com/bwmarrin/discordgo"
)

// 投稿されたメッセージのチェック
func checkMessage(s *discordgo.Session, m *discordgo.MessageCreate, c *discordgo.Channel) {
	//  botのメッセージには反応しない
	if m.Author.Bot {
		return
	}

	// botへのメンションがついていないメッセージには反応しない
	if strings.Index(m.Content, appConfig.BotName) == -1 {
		return
	}

	switch {
		case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", appConfig.BotName, "おはよう")):
			sendMessage(s, c, fmt.Sprintf("%s おはよう、%s。", createMention(m.Author), m.Author.Username))
		default :
			sendRandomMessage(s, c)
	}
}

//メッセージを送信
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

// ランダムにメッセージを送信する
func sendRandomMessage(s *discordgo.Session, c *discordgo.Channel) {
	l, err := createRandomMessageList()
	if err != nil {
		fmt.Println(err)
	}
	rand.Seed(time.Now().UnixNano())
	len := len(l)
	msg := l[rand.Intn(len)]

	sendMessage(s, c, msg)
}

// メンションを作成する
func createMention(user *discordgo.User) (string) {
	return fmt.Sprintf("<@%s>", user.ID)
}

// ランダムに送信するメッセージの配列を作成
// ※ randomMessages.txt を別途用意する
func createRandomMessageList() ([]string, error) {
	msgList := []string{}

	f, err := os.Open("./configs/randomMessages.txt")
	if err != nil {
		fmt.Println(err)
		return msgList, err
	}

	// テキスト内の各行の文字列を単一のメッセージとする
	s := bufio.NewScanner(f)
    for s.Scan() {
		log.Println(s.Text())
		msgList = append(msgList, s.Text())
	}
	
    if s.Err() != nil {
        // non-EOF error.
        fmt.Println(s.Err())
    }

	return msgList, nil
}

