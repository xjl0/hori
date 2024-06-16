package calendar

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
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

	req, err := http.NewRequestWithContext(ctxHttp, "GET", "https://kakoysegodnyaprazdnik.com/", nil)
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
	div := doc.Find("div.boxed-text").First()
	div.Find("li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		result += ":small_orange_diamond: " + text + "\n"
	})

	result = "**Праздники сегодня kakoysegodnyaprazdnik.com**\n" + result
	return result, nil
}
