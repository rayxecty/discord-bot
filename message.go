package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	sheets "google.golang.org/api/sheets/v4"
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
	case strings.HasPrefix(strings.ToLower(m.Content), fmt.Sprintf("%s %s", appConfig.BotName, "mtg")):
		// sendRandomCommanderCard(s, c)
		sendRandomCard(s, c)
	default:
		sendGSSMessage(s, m, c)
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

// GSSを参照してメッセージを送信する
func sendGSSMessage(s *discordgo.Session, m *discordgo.MessageCreate, c *discordgo.Channel) {
	spreadsheetId := appConfig.SpreadsheetId
	credentialFilePath := appConfig.SpreadsheetSecretPath

	client, err := httpClient(credentialFilePath)
	if err != nil {
		log.Fatal(err)
	}

	sheetService, err := sheets.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// _, err = sheetService.Spreadsheets.Get(spreadsheetId).Do()
	sheet, err := sheetService.Spreadsheets.Get(spreadsheetId).Ranges("A2:C50").IncludeGridData(true).Do()
	if err != nil {
		log.Fatalf("Unable to get Spreadsheets. %v", err)
	}

	// 一致したか判定
	hit := false
	rowCount := len(sheet.Sheets[0].Data[0].RowData)
	for i := 0; i < rowCount; i++ {
		rowData := sheet.Sheets[0].Data[0].RowData[i]
		reqMsg := rowData.Values[0].FormattedValue

		if strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", appConfig.BotName, reqMsg)) {
			hit = true
			resMsg := rowData.Values[1].FormattedValue

			// コマンドを変換
			// <username> → ユーザー名
			resMsg = strings.Replace(resMsg, "<username>", m.Author.Username, -1)

			// メンション設定
			var b bool
			b, _ = strconv.ParseBool(rowData.Values[2].FormattedValue)
			if b == true {
				resMsg = fmt.Sprintf("%s %s", createMention(m.Author), resMsg)
			}
			sendMessage(s, c, resMsg)
			break
		}
	}

	// 一致しない場合、既存のランダムメッセージを送信
	if !hit {
		fmt.Printf("NO HIT!\n")
		sendRandomMessage(s, c)
	}
}

// ランダムにメッセージを送信する
func sendRandomMessage(s *discordgo.Session, c *discordgo.Channel) {
	l, err := createRandomMessageList()
	if err != nil {
		log.Println(err)
	}
	rand.Seed(time.Now().UnixNano())
	len := len(l)
	msg := l[rand.Intn(len)]

	sendMessage(s, c, msg)
}

// メンションを作成する
func createMention(user *discordgo.User) string {
	return fmt.Sprintf("<@%s>", user.ID)
}

// ランダムに送信するメッセージの配列を作成
// ※ randomMessages.txt を別途用意する
func createRandomMessageList() ([]string, error) {
	msgList := []string{}

	f, err := os.Open("./configs/randomMessages.txt")
	if err != nil {
		log.Println(err)
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
		log.Println(s.Err())
	}

	return msgList, nil
}
