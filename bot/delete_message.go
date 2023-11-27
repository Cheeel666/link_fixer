package bot

import (
	"net/url"
	"strconv"
)

func (b *Bot) DeleteMessage(chatId int64, messageId int64) error {
	data := url.Values{}
	data.Add("chat_id", strconv.FormatUint(uint64(chatId), 10))
	data.Add("message_id", strconv.FormatUint(uint64(messageId), 10))

	err := b.httpClient.Post("deleteMessage", data, nil)
	if err != nil {
		return err
	}
	return nil
}
