package bot

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// TODO: add full struct
type sendMessageResp struct {
	OK bool `json:"ok"`
}

func (b *Bot) SendMessage(chatId int64, replyID int64, text string) error {
	data := url.Values{}
	data.Add("chat_id", strconv.FormatUint(uint64(chatId), 10))
	data.Add("text", text)
	data.Add("reply_to_message_id", strconv.FormatUint(uint64(replyID), 10))

	resp := &sendMessageResp{}
	err := b.httpClient.Post("sendMessage", data, resp)
	if err != nil {
		return err
	}
	if !resp.OK {
		return errors.New(fmt.Sprintf("response returned not OK, err: %s", err))
	}
	return nil
}
