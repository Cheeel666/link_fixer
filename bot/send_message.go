package bot

import (
	"fmt"
	"net/url"
	"strconv"
)

func (b *Bot) SendMessage(chatId int64, replyID int64, text string) error {
	data := url.Values{}
	data.Add("chat_id", strconv.FormatUint(uint64(chatId), 10))
	data.Add("text", text)
	data.Add("reply_to_message_id", strconv.FormatUint(uint64(replyID), 10))

	fmt.Println(data)
	err := b.httpClient.Post("sendMessage", data, nil)
	if err != nil {
		return err
	}
	return nil
}
