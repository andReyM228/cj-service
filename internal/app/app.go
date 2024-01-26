package app

import (
	"cj_service/internal/config"
	"cj_service/internal/handler/shedulers"
	tables_h "cj_service/internal/handler/tables"
	tables_r "cj_service/internal/repository/tables"
	tables_s "cj_service/internal/service/tables"
	"cj_service/internal/tg_handlers"
	"context"
	"github.com/andReyM228/lib/gpt3"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"gopkg.in/telebot.v3"
	stdLog "log"
	"net/http"
	"os"

	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"github.com/go-playground/validator/v10"
)

type App struct {
	config           config.Config
	serviceName      string
	tgbot            *telebot.Bot
	logger           log.Logger
	validator        *validator.Validate
	tablesRepo       tables_r.Repository
	tablesService    tables_s.Service
	tablesHandler    tables_h.Handler
	tgHandler        tg_handlers.Handler
	schedulerHandler shedulers.Handler
	clientHTTP       *http.Client
	errChan          chan errs.TgError
	chatGPT          gpt3.ChatGPT
	db               *sqlx.DB
}

func New(name string) App {
	return App{
		serviceName: name,
	}
}

func (a *App) Run(ctx context.Context) {
	a.initValidator()
	a.populateConfig()
	a.initLogger()
	a.listenErrs(ctx)
	a.initTgBot()
	a.initHTTPClient()
	a.initRepos()
	a.initServices()
	a.initHandlers()
	go a.runShedulers(ctx)
	a.listenTgHandlers()
}

func (a *App) listenErrs(ctx context.Context) {
	a.errChan = make(chan errs.TgError)

	go func() {
		for {
			select {
			case err := <-a.errChan:
				go func(err errs.TgError) {
					errs.HandleError(err.Err, a.logger, nil, err.ChatID)
				}(err)
			case <-ctx.Done():
				a.logger.Debug("ctx is done")
				os.Exit(1)

			}
		}
	}()
}

func (a *App) listenTgHandlers() {
	a.tgbot.Handle("/start", a.tgHandler.StartHandler)

	a.tgbot.Handle("информацияℹ️", a.tgHandler.InfoHandler)

	a.tgbot.Handle("вопросы❔", a.tgHandler.GetAllQuestionsHandler)

	a.tgbot.Handle("расписание📅", a.tgHandler.Schedule)

	a.tgbot.Handle(telebot.OnText, a.tgHandler.GetStudentHandler)

	a.logger.Debug("started tg handlers")

	a.tgbot.Start()
}

func (a *App) runShedulers(ctx context.Context) error {
	err := gocron.Every(5).Minutes().Do(a.schedulerHandler.GetInfo)
	if err != nil {
		return err
	}

	<-gocron.Start()

	return nil
}

func (a *App) initTgBot() {
	var err error
	bot, err := telebot.NewBot(telebot.Settings{
		Token:     a.config.TgBot.Token,
		Poller:    &telebot.LongPoller{Timeout: 10 * 60},
		ParseMode: telebot.ModeMarkdownV2,
	})
	if err != nil {
		return
	}

	a.tgbot = bot
}

func (a *App) initLogger() {
	a.logger = log.Init()
}

func (a *App) initValidator() {
	a.validator = validator.New()
}

func (a *App) initRepos() {
	a.tablesRepo = tables_r.NewRepository(a.logger, a.clientHTTP, a.config.Extra)

	a.logger.Debug("repos created")
}

func (a *App) initServices() {
	a.tablesService = tables_s.NewService(a.logger, a.tablesRepo)

	a.logger.Debug("services created")
}

func (a *App) initHandlers() {
	a.tablesHandler = tables_h.NewHandler(a.tablesService, a.tgbot)
	a.tgHandler = tg_handlers.NewHandler(a.tgbot, a.tablesHandler, a.errChan)
	a.schedulerHandler = shedulers.NewHandler(a.errChan, a.tablesService)

	a.logger.Debug("handlers created")
}

func (a *App) populateConfig() {
	cfg, err := config.ParseConfig()
	if err != nil {
		stdLog.Fatal(err)
	}

	err = cfg.ValidateConfig(a.validator)
	if err != nil {
		stdLog.Fatal(err)
	}

	a.config = cfg
}

func (a *App) initHTTPClient() {
	a.clientHTTP = http.DefaultClient
}
