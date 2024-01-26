package tables

import (
	"cj_service/internal/domain"
	"cj_service/internal/repository/tables"
	"fmt"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"strings"
	"time"
)

type Service struct {
	log        log.Logger
	repository tables.Repository
	storedInfo domain.Info
}

func NewService(log log.Logger, repository tables.Repository) Service {
	return Service{
		log:        log,
		repository: repository,
	}
}

func (s Service) Get(name string) (domain.Students, error) {
	resp, err := s.repository.Get(name)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Students{}, err
	}

	return resp, nil
}

func (s Service) GetQuestions(question string) (domain.Questions, error) {
	resp, err := s.repository.GetQuestions(question)
	if err != nil {
		s.log.Error(err.Error())
		return domain.Questions{}, err
	}

	return resp, nil
}

func (s Service) GetAllQuestions() (domain.Questions, error) {
	resp, err := s.repository.GetAllQuestions()
	if err != nil {
		s.log.Error(err.Error())
		return domain.Questions{}, err
	}

	return resp, nil
}

func (s *Service) SetInfo() error {
	resp, err := s.repository.GetInfo()
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	fmt.Println("setting info")

	s.storedInfo = resp

	return nil
}

func (s *Service) GetInfo() (domain.Info, error) {
	if s.storedInfo.Title == "" {
		err := s.SetInfo()
		if err != nil {
			s.log.Error(err.Error())
			return domain.Info{}, err
		}
	}

	return s.storedInfo, nil
}

func (s *Service) GetDate(day string) (string, error) {
	var result string

	weekday, ok := domain.Weekdays[day]
	if !ok {
		return "", errs.InternalError{Cause: "invalid weekday"}
	}

	if weekday == time.Now().Weekday() {
		return "сегодня", nil
	}

	if weekday > time.Now().Weekday() {
		difference := weekday - time.Now().Weekday()
		day := time.Now().Add((time.Hour * 24) * time.Duration(difference))

		result = day.Format("02.01")
	}

	if weekday < time.Now().Weekday() {
		difference := time.Now().Weekday() - weekday
		day := time.Now().Add(-(time.Hour * 24) * time.Duration(difference))

		if time.Now().Weekday() != time.Sunday {
			day = day.Add((time.Hour * 24) * 7)
		}

		result = day.Format("02.01")
	}

	result = strings.Replace(result, ".", `\.\`, -1)

	return result, nil
}
