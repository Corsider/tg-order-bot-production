package main

import (
	"encoding/json"
)

// Here are some useful structs

// Bot actually...
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot,omitempty"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username,omitempty"`
}

// Food Types...
const (
	PizzaPepperony int = iota
	PizzaFourCheese
	Cola
	KinderSurprise
)

type UpdateChannel <-chan Update

type Response struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type ChatMember interface {
}

type ChatMemberUpdated struct {
	Chat          Chat       `json:"chat"`
	From          User       `json:"from"`
	Date          int        `json:"date"`
	OldChatMember ChatMember `json:"old_chat_member"`
	NewChatMember ChatMember `json:"new_chat_member"`
}

type Update struct {
	UpdateId      int                `json:"update_id"`
	Message       *Message           `json:"message"`
	CallbackQuery *CallbackQuery     `json:"callback_query"`
	MyChatMember  *ChatMemberUpdated `json:"my_chat_member"`
}

type UpdateParams struct {
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
	Timeout int `json:"timeout"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

type Message struct {
	MessageID   int                   `json:"message_id"`
	Chat        *Chat                 `json:"chat"`
	Text        string                `json:"text,omitempty"`
	Entities    []MessageEntity       `json:"entities,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	From        User                  `json:"from"`
}

type Chat struct {
	Id int64 `json:"id"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string  `json:"text"`
	CallbackData *string `json:"callback_data,omitempty"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data,omitempty"`
	From    User     `json:"from"`
}
