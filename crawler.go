package main

import (
	"html"
	"log"
	"os"
	"regexp"

	"github.com/bachue/wechat-official-account-crawler/httputils"
	"github.com/bachue/wechat-official-account-crawler/msglist"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s URL FILENAME\n", os.Args[0])
	}
	url := os.Args[1]
	outputFilePath := os.Args[2]
	msglistFetcher := msglist.NewFetcher(url)
	body := httputils.Get(url)
	msgList := extractMsgList(string(body))
	msgList.UnescapeSlashFromURL()
	msgList.ConvertToPDF()
	for {
		last := msgList.Last()
		if last == nil {
			break
		}
		lastId := last.CommMsgInfo.Id
		msgList = msglistFetcher.FetchNextList(lastId, 10)
		if msgList == nil {
			break
		}
		msgList.UnescapeHTMLFromURL()
		msgList.ConvertToPDF()
	}
	msgList.ConcatPDFs(outputFilePath)
}

func extractMsgList(body string) *msglist.MessageList {
	re := regexp.MustCompile("msgList\\s+=\\s+'([^']+)'")
	results := re.FindStringSubmatch(body)
	if len(results) != 2 {
		log.Fatalf("Failed to extract msgList from request body: %s\n", body)
	}
	msgListStr := html.UnescapeString(results[1])
	msgListStr = html.UnescapeString(msgListStr)
	return msglist.ExtractFromJSON([]byte(msgListStr))
}
