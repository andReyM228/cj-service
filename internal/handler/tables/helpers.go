package tables

import (
	"gopkg.in/telebot.v3"
)

func (h Handler) PrepareButtons() ([]telebot.Btn, error) {
	r := h.tgBot.NewMarkup()

	infoButton := r.Text("информацияℹ️")
	raspButton := r.Text("расписание📅")
	questionsButton := r.Text("вопросы❔")

	return []telebot.Btn{infoButton, raspButton, questionsButton}, nil
}
