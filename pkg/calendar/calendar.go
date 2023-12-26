package calendar

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

func CalendarReq() (string, error) {
	// Создаем прокси-диалог
	ctxHttp, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Создаем клиента с настроенным транспортом
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctxHttp, "GET", "https://kakoysegodnyaprazdnik.ru/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return res.Status, nil
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	var result string
	count := 0
	doc.Find(".listing .listing_wr .main").Each(func(i int, s *goquery.Selection) {
		if count < 20 {
			text := s.Find("span").First().Text()
			text = strings.Replace(text, "США", ":flag_um:", 1)
			text = strings.Replace(text, "Япония", ":flag_jp:", 1)
			result += ":small_blue_diamond: " + text + "\n"
		}
		count++
	})
	result = "**Праздники сегодня kakoysegodnyaprazdnik.ru**\n" + result
	return result, nil
}
