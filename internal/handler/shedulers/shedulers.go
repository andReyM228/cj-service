package shedulers

import (
	"cj_service/internal/service/tables"
	"github.com/andReyM228/lib/errs"
)

type Handler struct {
	errChan       chan errs.TgError
	tablesService tables.Service
}

func NewHandler(errChan chan errs.TgError, service tables.Service) Handler {
	return Handler{
		errChan:       errChan,
		tablesService: service,
	}
}

func (h Handler) GetInfo() error {
	err := h.tablesService.SetInfo()
	if err != nil {
		return err
	}

	return nil
}
