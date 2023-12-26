package discord

import (
	"discordbotgo/pkg/calendar"
	"discordbotgo/pkg/shiki"
	"encoding/json"
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"github.com/bwmarrin/discordgo"
	"github.com/pemistahl/lingua-go"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func oconHours(i int) string {
	switch i {
	case 1, 21:
		return "час"
	case 2, 3, 4, 22, 23:
		return "часа"
	}
	return "часов"
}

func oconMinutes(i int) string {
	switch i {
	case 1, 21:
		return "минуту"
	case 2, 3, 4, 22, 23, 24:
		return "минуты"
	}
	return "минут"
}

func translateText(targetLanguage, text string) (string, error) {
	return gt.Translate(text, targetLanguage, `ru`)
}

// Проверка состояния бота
func isALife(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "бот жив?" {
		stat, err := json.Marshal(m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```json\n%s\n```", string(stat)))
	}
	if m.Content == "шики" {
		sendNews(s)
	}
}

// Когда на сервере создаётся мероприятие, ссылка на событие отправляется в чат
func eventCreate(s *discordgo.Session, e *discordgo.GuildScheduledEventCreate) {
	testChannel := viper.GetString("discord.testchannel")
	mainChannel := viper.GetString("discord.mainchannel")
	codeInvite := viper.GetString("discord.codeinvite")

	if e.Name == "Test" {
		s.ChannelMessageSend(testChannel, `https://discord.gg/`+codeInvite+`?event=`+e.ID)
		s.ChannelMessageSend(testChannel, "Событие с названием Test. Инвайт код "+codeInvite+" Инвайт ID "+e.ID)
		return
	}
	s.ChannelMessageSend(mainChannel, `https://discord.gg/EFyjYbqn7E?event=`+e.ID)
}

// Основной обработчик сообщений
func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	helloEmote := viper.GetString("discord.helloemote")
	mainChannel := viper.GetString("discord.mainchannel")
	mediaChannel := viper.GetString("discord.mediachannel")
	testChannel := viper.GetString("discord.testchannel")

	if m.Author.ID == s.State.User.ID {
		return
	}

	//Ставит эмоцию приветствия
	if m.Content == "<:"+helloEmote+">" {
		s.MessageReactionAdd(m.ChannelID, m.ID, helloEmote)
	}

	//Под ссылки ставить две эмоции
	match, _ := regexp.MatchString(`https\:\/\/youtu\.be\/.*`, m.Content)
	match2, _ := regexp.MatchString(`https\:\/\/www\.youtube\.com\/watch.*`, m.Content)
	match3, _ := regexp.MatchString(`https\:\/\/coub\.com\/view\/.*`, m.Content)
	if match || match2 || match3 || m.ChannelID == mediaChannel {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
		s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
	}

	//Текущие праздники
	if m.Content == "календарь" {
		a, err := calendar.CalendarReq()
		if err != nil {
			fmt.Println(err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, a)
	}

	//Вычисляет время по количеству серий
	reg := `(?i)Сколько по времени (\d+) сер[А-Яа-я]{2}\?`
	ref, _ := regexp.MatchString(reg, m.Content)
	if ref {
		d, _ := regexp.Compile(reg)
		sch, _ := strconv.Atoi(d.FindStringSubmatch(m.Content)[1])
		allMin := sch * 24
		hour := allMin / 60
		minute := allMin - hour*60
		future := time.Now()
		r := future.Add(time.Minute * time.Duration(allMin))
		var dHour, dMinute string
		if r.Hour() < 10 {
			dHour = `0` + strconv.Itoa(r.Hour())
		} else {
			dHour = strconv.Itoa(r.Hour())
		}
		if r.Minute() < 10 {
			dMinute = `0` + strconv.Itoa(r.Minute())
		} else {
			dMinute = strconv.Itoa(r.Minute())
		}
		zone, _ := r.Zone()
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`%d %s %d %s (1 серия 24 минуты). Если начать сейчас, то закончим в %s:%s %s`, hour, oconHours(hour), minute, oconMinutes(minute), dHour, dMinute, zone))
	}

	//Автоматом переводит предложения на русский язык
	if m.ChannelID == mainChannel || m.ChannelID == testChannel {
		if strings.Contains(m.Content, `:`) || strings.Contains(m.Content, `>`) || strings.Contains(m.Content, `<`) || strings.Contains(m.Content, `/`) || strings.Contains(m.Content, `@`) {
			return
		}
		languages := []lingua.Language{
			lingua.Azerbaijani,
			lingua.Japanese,
			lingua.English,
		}
		detector := lingua.NewLanguageDetectorBuilder().FromLanguages(languages...).Build()
		language, exists := detector.DetectLanguageOf(m.Content)
		if language.IsoCode639_1().String() == "EN" {
			return
		}
		if exists {
			readyText, err := translateText(language.IsoCode639_1().String(), m.Content)
			if err != nil || len(readyText) == 0 {
				return
			}
			s.ChannelMessageSendReply(m.ChannelID, `**`+language.String()+`**: `+readyText, &discordgo.MessageReference{
				MessageID: m.Message.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			})
		}
	}
}

// Запрос и отправка последней новости с Шикимори
func sendNews(s *discordgo.Session) {
	lasted := viper.GetInt("discord.lasted")
	newsChannel := viper.GetString("discord.newschannel")
	var res []shiki.Topic
	//Получим новость
	err := shiki.ShikiGetTopics(&res)
	if err != nil {
		log.Println(err)
	}
	if lasted == res[0].Id {
		return
	}
	//Получим и запишем ID последней новости
	viper.Set("discord.lasted", res[0].Id)
	viper.WriteConfig()
	embed := discordgo.MessageEmbed{
		URL:         `https://shikimori.one` + res[0].Forum.Url + "/" + strconv.Itoa(res[0].Id),
		Type:        "rich",
		Title:       res[0].TopicTitle,
		Description: res[0].Body,
		Timestamp:   res[0].CreatedAt,
		Color:       123222,
		Footer: &discordgo.MessageEmbedFooter{
			Text: res[0].Forum.Name,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://kawai.shikimori.one" + res[0].Linked.Image.Original,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://kawai.shikimori.one" + res[0].Linked.Image.Preview,
		},
		Video:    nil,
		Provider: nil,
		Author:   nil,
		Fields:   nil,
	}
	s.ChannelMessageSendEmbed(newsChannel, &embed)
}
