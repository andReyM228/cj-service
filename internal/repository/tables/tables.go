package tables

import (
	"cj_service/internal/config"
	"cj_service/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"net/http"
	"net/url"
)

type Repository struct {
	log    log.Logger
	client *http.Client
	config config.Extra
}

func NewRepository(log log.Logger, client *http.Client, config config.Extra) Repository {
	return Repository{
		log:    log,
		client: client,
		config: config,
	}
}

func (r Repository) Get(name string) (domain.Students, error) {
	ur, err := url.Parse(r.config.GoogleSheets.UrlNames)
	if err != nil {
		return domain.Students{}, errs.InternalError{err.Error()}
	}

	values := ur.Query()

	values.Add("query", name)

	ur.RawQuery = values.Encode()

	resp, err := r.client.Get(ur.String())
	if err != nil {
		return domain.Students{}, errs.InternalError{err.Error()}
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	var students []domain.Student

	if resp.StatusCode != http.StatusOK {
		return domain.Students{}, errs.InternalError{Cause: fmt.Sprintf("status code: %d", resp.StatusCode)}
	}

	err = json.NewDecoder(resp.Body).Decode(&students)
	if err != nil {
		return domain.Students{}, errs.InternalError{err.Error()}
	}

	return domain.Students{Students: students}, nil
}

func (r Repository) GetQuestions(question string) (domain.Questions, error) {
	ur, err := url.Parse(r.config.GoogleSheets.UrlQuestions)
	if err != nil {
		return domain.Questions{}, errs.InternalError{err.Error()}
	}

	values := ur.Query()

	values.Add("query", question)

	ur.RawQuery = values.Encode()

	resp, err := r.client.Get(ur.String())
	if err != nil {
		return domain.Questions{}, errs.InternalError{err.Error()}
	}
	defer resp.Body.Close()

	var questions []domain.Question

	if resp.StatusCode != http.StatusOK {
		return domain.Questions{}, errs.InternalError{Cause: fmt.Sprintf("status code: %d", resp.StatusCode)}
	}

	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		return domain.Questions{}, errs.InternalError{err.Error()}
	}

	return domain.Questions{Questions: questions}, nil
}

// GetAllQuestions returns all questions from google sheets
func (r Repository) GetAllQuestions() (domain.Questions, error) {
	ur, err := url.Parse(r.config.GoogleSheets.UrlAllQuestions)
	if err != nil {
		return domain.Questions{}, errs.InternalError{Cause: err.Error()}
	}

	resp, err := r.client.Get(ur.String())
	if err != nil {
		return domain.Questions{}, errs.InternalError{Cause: err.Error()}
	}
	defer resp.Body.Close()

	var questions []domain.Question

	if resp.StatusCode != http.StatusOK {
		return domain.Questions{}, errs.InternalError{Cause: fmt.Sprintf("status code: %d", resp.StatusCode)}
	}

	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		return domain.Questions{}, errs.InternalError{Cause: err.Error()}
	}

	return domain.Questions{Questions: questions}, nil
}

func (r Repository) GetInfo() (domain.Info, error) {
	ur, err := url.Parse(r.config.GoogleSheets.UrlInfo)
	if err != nil {
		return domain.Info{}, errs.InternalError{Cause: err.Error()}
	}

	resp, err := r.client.Get(ur.String())
	if err != nil {
		return domain.Info{}, errs.InternalError{Cause: err.Error()}
	}
	defer resp.Body.Close()

	var info []domain.Info

	if resp.StatusCode != http.StatusOK {
		return domain.Info{}, errs.InternalError{Cause: fmt.Sprintf("status code: %d", resp.StatusCode)}
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return domain.Info{}, errs.InternalError{Cause: err.Error()}
	}

	return info[0], nil
}
