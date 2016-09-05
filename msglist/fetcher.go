package msglist

import (
	"encoding/json"
	"log"
	urlutils "net/url"
	"strconv"

	"github.com/bachue/wechat-official-account-crawler/httputils"
)

type ListFetcher struct {
	BaseURL    *urlutils.URL
	Biz        string
	Uin        string
	Key        string
	PassTicket string
}

func NewFetcher(urlstr string) *ListFetcher {
	var fetcher ListFetcher
	url, err := urlutils.Parse(urlstr)
	if err != nil {
		log.Fatalf("Failed to parse URL %s: %s\n", urlstr, err)
	}
	queryValues, err := urlutils.ParseQuery(url.RawQuery)
	if err != nil {
		log.Fatalf("Failed to parse Query from URL %s: %s\n", urlstr, err)
	}
	url.RawQuery = ""
	url.Fragment = ""
	fetcher.BaseURL = url
	fetcher.Biz = queryValues.Get("__biz")
	fetcher.Uin = queryValues.Get("uin")
	fetcher.Key = queryValues.Get("key")
	fetcher.PassTicket = queryValues.Get("pass_ticket")

	return &fetcher
}

type nextPageBody struct {
	GeneralMsgList string `json:"general_msg_list"`
}

func (fetcher *ListFetcher) FetchNextList(fromId int, count int) *MessageList {
	var nextPage nextPageBody
	body := fetcher.fetch(fromId, count)
	err := json.Unmarshal(body, &nextPage)
	if err != nil {
		log.Fatalf("Failed to parse JSON (reason: %s) from Next Page: %s\n", err, body)
	}
	if len(nextPage.GeneralMsgList) > 0 {
		return ExtractFromJSON([]byte(nextPage.GeneralMsgList))
	} else {
		return nil
	}
}

func (fetcher *ListFetcher) fetch(fromId int, count int) []byte {
	url := fetcher.fetchURL(fromId, count)
	return httputils.Get(url)
}

func (fetcher *ListFetcher) fetchURL(fromId int, count int) string {
	var url urlutils.URL = *fetcher.BaseURL
	query := urlutils.Values{}
	query.Set("__biz", fetcher.Biz)
	query.Set("uin", fetcher.Uin)
	query.Set("key", fetcher.Key)
	query.Set("pass_ticket", fetcher.PassTicket)
	query.Set("f", "json")
	query.Set("frommsgid", strconv.Itoa(fromId))
	query.Set("count", strconv.Itoa(count))
	url.RawQuery = query.Encode()
	return url.String()
}
