package main

import (
	"fmt"
	"log"

	"github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/bwmarrin/discordgo"
)

type CardInfo struct {
	Name         string
	MultiverseId uint
	ImageUrl     string
	IsJapanese   bool
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

	fmt.Println(cardInfo)

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
	fmt.Println(c)
	ret := findJapaneseCardInfo(c)
	if !ret.IsJapanese {
		fmt.Println("日本語がないセットのカードです")
		// 取得したカード名で日本語の情報があるカードを検索
		cards, err := mtg.NewQuery().Where(mtg.CardName, c.Name).All()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(fmt.Sprintf("同名カードが%d枚存在します", len(cards)))
		if len(cards) > 0 {
			for _, card := range cards {
				ret = findJapaneseCardInfo(card)
				if ret.IsJapanese {
					break
				}
			}

		}
	}

	return ret
}

// 日本語のカード情報を探索
func findJapaneseCardInfo(card *mtg.Card) CardInfo {
	var ret CardInfo
	foreignNames := card.ForeignNames

	// 日本語のカード情報を取得
	var japaneseCard mtg.ForeignCardName
	for _, fn := range foreignNames {
		if fn.Language == "Japanese" {
			japaneseCard = fn
			fmt.Println(japaneseCard)
		}
	}

	if japaneseCard.Name != "" {
		ret.Name = japaneseCard.Name
		ret.MultiverseId = japaneseCard.MultiverseId
		ret.ImageUrl = createCardImageUrl(japaneseCard.MultiverseId)
		ret.IsJapanese = true
	} else {
		// 日本語の情報がない場合、オリジナル(英語)のカード情報を返還
		ret.Name = card.Name
		ret.MultiverseId = uint(card.MultiverseId)
		ret.ImageUrl = card.ImageUrl
		ret.IsJapanese = false
	}

	return ret
}

// カードIDからURL生成
func createCardImageUrl(id uint) string {
	return fmt.Sprintf("http://gatherer.wizards.com/Handlers/Image.ashx?multiverseid=%d&type=card", id)
}
