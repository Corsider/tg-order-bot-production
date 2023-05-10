package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const APIURL = "https://api.telegram.org/bot"

type Bot struct {
	Token  string
	Me     User
	Url    string
	Client *http.Client
}

func CreateBot(token string) (*Bot, error) {
	bot := &Bot{
		Token:  token,
		Url:    APIURL,
		Client: &http.Client{},
	}
	// getting User

	usr, err := bot.SimpleRawRequest("getMe", nil)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal(usr.Result, &user)
	bot.Me = user

	return bot, err
}

func (b *Bot) SimpleRawRequest(method string, data map[string]string) (*Response, error) {
	requstURL := b.Url + b.Token + "/" + method

	var params url.Values
	if data == nil {
		params = url.Values{}
	} else {
		params = url.Values{}
		for key, value := range data {
			params.Set(key, value)
		}
	}

	req, err := http.NewRequest(http.MethodPost, requstURL, strings.NewReader(params.Encode()))
	if err != nil {
		return &Response{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var output Response
	err = json.NewDecoder(response.Body).Decode(&output)
	return &output, err
}

func (b *Bot) sendAnswerCallbackQuery(callbackQueryID string) (bool, error) {
	param := map[string]string{
		"callback_query_id": callbackQueryID,
	}
	_, err := b.SimpleRawRequest("answerCallbackQuery", param)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *Bot) sendMessage(chatID int64, text string, replyMarkup *InlineKeyboardMarkup) (Message, error) {
	var params map[string]string
	if replyMarkup == nil {
		params = map[string]string{
			"chat_id": strconv.FormatInt(chatID, 10),
			"text":    text,
		}
	} else {
		data, err := json.Marshal(replyMarkup)
		if err != nil {
			return Message{}, err
		}
		params = map[string]string{
			"chat_id":      strconv.FormatInt(chatID, 10),
			"text":         text,
			"reply_markup": string(data),
		}
	}

	resp, err := b.SimpleRawRequest("sendMessage", params)
	if err != nil {
		return Message{}, err
	}
	var msg Message
	err = json.Unmarshal(resp.Result, &msg)
	return msg, err
}

func (b *Bot) sendEditMessageText(chatID int64, messageID int, text string, replyMarkup *InlineKeyboardMarkup) (Message, error) {
	var params map[string]string
	if replyMarkup == nil {
		params = map[string]string{
			"chat_id":    strconv.FormatInt(chatID, 10),
			"text":       text,
			"message_id": strconv.Itoa(messageID),
		}
	} else {
		data, err := json.Marshal(replyMarkup)
		if err != nil {
			return Message{}, err
		}
		params = map[string]string{
			"chat_id":      strconv.FormatInt(chatID, 10),
			"text":         text,
			"message_id":   strconv.Itoa(messageID),
			"reply_markup": string(data),
		}
	}
	resp, err := b.SimpleRawRequest("editMessageText", params)
	if err != nil {
		return Message{}, err
	}
	var msg Message
	err = json.Unmarshal(resp.Result, &msg)
	return msg, err
}

func (b *Bot) sendEditMessageMarkup(chatID int64, messageID int, replyMarkup *InlineKeyboardMarkup) (Message, error) {
	data, err := json.Marshal(replyMarkup)
	if err != nil {
		return Message{}, err
	}
	params := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"message_id":   strconv.Itoa(messageID),
		"reply_markup": string(data),
	}
	resp, err := b.SimpleRawRequest("editMessageReplyMarkup", params)
	if err != nil {
		return Message{}, err
	}
	var msg Message
	err = json.Unmarshal(resp.Result, &msg)
	return msg, err
}

func (b *Bot) sendEditMessageTextAndMarkup(chatID int64, messageID int, text string, replyMarkup *InlineKeyboardMarkup) (Message, error) {
	data, err := json.Marshal(replyMarkup)
	if err != nil {
		return Message{}, err
	}
	params := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"message_id":   strconv.Itoa(messageID),
		"text":         text,
		"reply_markup": string(data),
	}
	_, err = b.SimpleRawRequest("editMessageText", params)
	if err != nil {
		return Message{}, err
	}
	resp, err := b.SimpleRawRequest("editMessageReplyMarkup", params)
	if err != nil {
		return Message{}, err
	}
	var msg Message
	err = json.Unmarshal(resp.Result, &msg)
	return msg, err
}

func (b *Bot) sendDeleteMessage(chatID int64, messageID int) (bool, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.Itoa(messageID),
	}
	_, err := b.SimpleRawRequest("deleteMessage", params)
	if err != nil {
		return false, err
	}
	return true, err
}

func (b *Bot) UpdatesChannel(params UpdateParams) UpdateChannel {
	channel := make(chan Update, 100)
	go func() {
		for {
			paramTo := make(map[string]string)
			paramTo["offset"] = strconv.Itoa(params.Offset)
			paramTo["limit"] = strconv.Itoa(params.Limit)
			paramTo["timeout"] = strconv.Itoa(params.Timeout)
			var updates []Update
			response, err := b.SimpleRawRequest("getUpdates", paramTo)
			if err != nil {
				updates = []Update{}
			}
			err = json.Unmarshal(response.Result, &updates)
			if err != nil {
				log.Fatal(err)
			}
			for _, update := range updates {
				if update.UpdateId >= params.Offset {
					params.Offset = update.UpdateId + 1
					channel <- update
				}
			}
		}
	}()
	return channel
}

func CreateInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton
	keyboard = append(keyboard, rows...)
	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func CreateInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton
	row = append(row, buttons...)
	return row
}

func CreateInlineKeyboardButton(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: &data,
	}
}

func (m *Message) CheckIfCommand() bool {
	if m.Entities == nil {
		return false
	}
	t := m.Entities[0]
	return t.Type == "bot_command"
}
