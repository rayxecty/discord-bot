package main

import (
	"fmt"
	"log"

	"github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/bwmarrin/discordgo"
)

// ランダムに統率者カードを送信する
func sendRandomCommanderCard(s *discordgo.Session, c *discordgo.Channel) {
	cards, err := mtg.NewQuery().Where(mtg.CardGameFormat, "Commander").Where(mtg.CardType, "Creature").Where(mtg.CardSupertypes, "Legendary").Random(1)
	if err != nil {
		log.Panic(err)
	}

	for _, card := range cards {
		sendMessage(s, c, getJapaneseCardImageUrl(card))
	}
}

// 日本語のカードを送信
func getJapaneseCardImageUrl(card *mtg.Card) string {
	var foreignNames = card.ForeignNames

	// 日本語のカード情報を取得
	var japaneseCardName mtg.ForeignCardName
	for _, fn := range foreignNames {
		if fn.Language == "Japanese" {
			japaneseCardName = fn
			fmt.Println(japaneseCardName)
		}
	}

	var cardImageUrl string
	if japaneseCardName.Name != "" {
		cardImageUrl = createCardImageUrl(japaneseCardName.MultiverseId)
	} else {
		// 日本語のカードがない場合、オリジナル(英語)のカードを送信
		cardImageUrl = card.ImageUrl
	}

	return cardImageUrl
}

// カードIDからURL生成
func createCardImageUrl(id uint) string {
	return fmt.Sprintf("http://gatherer.wizards.com/Handlers/Image.ashx?multiverseid=%d&type=card", id)
}
