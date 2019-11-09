package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/bwmarrin/discordgo"
)

type CardInfo struct {
	Name         string
	MultiverseId uint
	ImageUrl     string
	IsJapanese   bool
}

// ランダムにカードを送信する
func sendRandomCard(s *discordgo.Session, c *discordgo.Channel) {
	// ランダムに表示させるURLをリクエスト
	target_url := "http://gatherer.wizards.com/Pages/Card/Details.aspx?action=random"
	req, _ := http.NewRequest("GET", target_url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)

	// レスポンスからクエリを取得
	raw_query := resp.Request.URL.RawQuery
	defer resp.Body.Close()

	// カードのみを表示するURLを作成
	base_url := "https://gatherer.wizards.com/Handlers/Image.ashx"
	randomcard_url := base_url + "?" + raw_query + "&type=card"
	fmt.Println(randomcard_url)

	// URLを送信
	sendMessage(s, c, randomcard_url)
}

// ランダムに統率者カードを送信する
func sendRandomCommanderCard(s *discordgo.Session, c *discordgo.Channel) {
	cards, err := mtg.NewQuery().Where(mtg.CardGameFormat, "Commander").Where(mtg.CardType, "Creature").Where(mtg.CardSupertypes, "Legendary").Random(1)
	if err != nil {
		log.Panic(err)
	}

	for _, card := range cards {
		sendMessage(s, c, sendJapaneseCard(card))
	}
}

// 日本語のカードを送信
func sendJapaneseCard(c *mtg.Card) string {
	var msg string
	cardInfo := getJapaneseCardInfo(c)

	log.Println(cardInfo)

	msg += cardInfo.Name + "\n"
	if !cardInfo.IsJapanese {
		msg += "日本語のカードがなかったよ\n"
	}
	if cardInfo.MultiverseId == 0 {
		msg += "カード画像がなかったよ\n"
	}
	msg += cardInfo.ImageUrl

	return msg
}

// 日本語のカード画像を取得
func getJapaneseCardInfo(c *mtg.Card) CardInfo {
	log.Println(c)
	ret := findJapaneseCardInfo(c)
	if !ret.IsJapanese {
		// 取得したカード名で日本語の情報があるカードを検索
		cards, err := mtg.NewQuery().Where(mtg.CardName, c.Name).All()
		if err != nil {
			log.Panic(err)
		}
		log.Println(fmt.Sprintf("%s は %d 枚存在します", c.Name, len(cards)))

		if len(cards) > 0 {
			for _, card := range cards {
				if card.Name != c.Name {
					log.Println(card.Name + " is not " + c.Name)
					continue
				}
				ret = findJapaneseCardInfo(card)
				if ret.IsJapanese {
					if ret.MultiverseId != 0 {
						break
					} else {
						continue
					}
				}
			}
		}
	}

	return ret
}

// 日本語のカード情報を探索
func findJapaneseCardInfo(c *mtg.Card) CardInfo {
	var ret CardInfo
	foreignNames := c.ForeignNames

	// 日本語のカード情報を取得
	var japaneseCard mtg.ForeignCardName
	for _, fn := range foreignNames {
		if fn.Language == "Japanese" {
			japaneseCard = fn
			log.Println(japaneseCard)
		}
	}

	if japaneseCard.Name != "" {
		ret.Name = japaneseCard.Name
		ret.MultiverseId = japaneseCard.MultiverseId
		ret.ImageUrl = createCardImageUrl(japaneseCard.MultiverseId)
		ret.IsJapanese = true
	} else {
		// 日本語の情報がない場合、オリジナル(英語)のカード情報を返還
		log.Println(c.SetName + " は日本語版がありません")
		ret.Name = c.Name
		ret.MultiverseId = uint(c.MultiverseId)
		ret.ImageUrl = c.ImageUrl
		ret.IsJapanese = false
	}

	return ret
}

// カードIDからURL生成
func createCardImageUrl(id uint) string {
	return fmt.Sprintf("http://gatherer.wizards.com/Handlers/Image.ashx?multiverseid=%d&type=card", id)
}
