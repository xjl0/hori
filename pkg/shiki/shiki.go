package shiki

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Topic struct {
	Id         int    `json:"id"`
	TopicTitle string `json:"topic_title"`
	Body       string `json:"body"`
	Forum      Forum  `json:"forum"`
	CreatedAt  string `json:"created_at"`
	Linked     Linked `json:"linked"`
}

type Forum struct {
	Id        int    `json:"id"`
	Position  int    `json:"position"`
	Name      string `json:"name"`
	Permalink string `json:"permalink"`
	Url       string `json:"url"`
}
type Linked struct {
	Id    int   `json:"id"`
	Image Image `json:"image"`
}

type Image struct {
	Original string `json:"original"`
	Preview  string `json:"preview"`
}

func ShikiGetTopics(target interface{}) error {
	spaceClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, `https://shikimori.one/api/topics?limit=1&forum=news&page=1`, nil)
	if err != nil {
		return err
	}
	res, err := spaceClient.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &target)
	if err != nil {
		return err
	}
	
	return nil
}
