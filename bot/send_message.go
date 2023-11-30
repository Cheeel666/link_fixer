package bot

import (
	"net/url"
	"strconv"
)

func (b *Bot) SendMessage(chatId int64, replyID int64, text string) error {
	data := url.Values{}
	data.Add("chat_id", strconv.FormatUint(uint64(chatId), 10))
	data.Add("text", text)
	data.Add("reply_to_msg_id", strconv.FormatUint(uint64(replyID), 10))

	err := b.httpClient.Post("sendMessage", data, nil)
	if err != nil {
		return err
	}
	return nil
}
