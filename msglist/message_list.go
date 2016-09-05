package msglist

import (
	"encoding/json"
	"html"
	"log"
	"strings"
)

type CommonMessageInfo struct {
	Id       int `json:"id"`
	Datetime int `json:"datetime"`
}

type AppMessageExtensionInfo struct {
	Title      string `json:"title"`
	ContentUrl string `json:"content_url"`
}

type MessageInfo struct {
	CommMsgInfo   CommonMessageInfo       `json:"comm_msg_info"`
	AppMsgExtInfo AppMessageExtensionInfo `json:"app_msg_ext_info"`
}

type MessageList struct {
	List []*MessageInfo `json:"list"`
}

func ExtractFromJSON(data []byte) *MessageList {
	var msgList MessageList
	err := json.Unmarshal(data, &msgList)
	if err != nil {
		log.Fatalf("Failed to parse JSON (reason: %s) from extracted msgList string: %s\n", err, string(data))
	}
	return &msgList
}

func (msgList *MessageList) Last() *MessageInfo {
	size := len(msgList.List)
	if size > 0 {
		return msgList.List[size-1]
	} else {
		return nil
	}
}

func (msgList *MessageList) UnescapeSlashFromURL() {
	for _, info := range msgList.List {
		info.AppMsgExtInfo.ContentUrl = strings.Replace(info.AppMsgExtInfo.ContentUrl, "\\/", "/", -1)
	}
}

func (msgList *MessageList) UnescapeHTMLFromURL() {
	for _, info := range msgList.List {
		info.AppMsgExtInfo.ContentUrl = html.UnescapeString(info.AppMsgExtInfo.ContentUrl)
	}
}
