package tables

import (
	"cj_service/internal/domain"
	"cj_service/internal/service/tables"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	tablesService tables.Service
	tgBot         *telebot.Bot
}

func NewHandler(service tables.Service, tgbot *telebot.Bot) Handler {
	return Handler{
		tablesService: service,
		tgBot:         tgbot,
	}
}

func (h Handler) Get(name string) (domain.Students, error) {
	resp, err := h.tablesService.Get(name)
	if err != nil {
		return domain.Students{}, err
	}

	return resp, nil
}

func (h Handler) GetQuestions(question string) (domain.Questions, error) {
	resp, err := h.tablesService.GetQuestions(question)
	if err != nil {
		return domain.Questions{}, err
	}

	return resp, nil
}

func (h Handler) GetAllQuestions() (domain.Questions, error) {
	resp, err := h.tablesService.GetAllQuestions()
	if err != nil {
		return domain.Questions{}, err
	}

	return resp, nil
}

func (h *Handler) GetInfo() (domain.Info, error) {
	info, err := h.tablesService.GetInfo()
	if err != nil {
		return domain.Info{}, err
	}
	return info, nil
}

func (h *Handler) GetDate(day string) (string, error) {
	date, err := h.tablesService.GetDate(day)
	if err != nil {
		return "", err
	}

	return date, nil
}
