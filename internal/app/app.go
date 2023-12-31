package app

import (
	"context"
	stdLog "log"
	"net/http"
	"os"
	"strings"

	"cj_service/internal/config"
	"cj_service/internal/domain"
	car_handler "cj_service/internal/handler/car"
	user_handler "cj_service/internal/handler/user"
	"cj_service/internal/repository/cars"
	"cj_service/internal/repository/user"
	"cj_service/internal/service/car"
	user_service "cj_service/internal/service/user"
	"cj_service/internal/tg_handlers"

	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/gpt3"
	"github.com/andReyM228/lib/log"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	config                      config.Config
	serviceName                 string
	tgbot                       *tgbotapi.BotAPI
	logger                      log.Logger
	validator                   *validator.Validate
	usersRepo                   user.Repository
	carsRepo                    cars.Repository
	userService                 user_service.Service
	carService                  car.Service
	userHandler                 user_handler.Handler
	carHandler                  car_handler.Handler
	tgHandler                   tg_handlers.Handler
	clientHTTP                  *http.Client
	errChan                     chan errs.TgError
	loginUsers                  map[int64]string
	chatGPT                     gpt3.ChatGPT
	processingRegistrationUsers domain.ProcessingRegistrationUsers
	processingLoginUsers        domain.ProcessingLoginUsers
}

func New(name string) App {
	return App{
		serviceName: name,
	}
}

func (a *App) initGPT() {
	a.chatGPT = gpt3.Init(a.config.ChatGPT.Key, a.config.ChatGPT.Model)
}

func (a *App) Run(ctx context.Context) {
	a.initValidator()
	a.populateConfig()
	a.initLogger()
	a.initGPT()
	a.listenErrs(ctx)
	a.initTgBot()
	a.initHTTPClient()
	a.initRepos()
	a.initServices()
	a.initHandlers()
	a.listenTgBot()
}

func (a *App) listenErrs(ctx context.Context) {
	a.errChan = make(chan errs.TgError)

	go func() {
		for {
			select {
			case err := <-a.errChan:
				go func(err errs.TgError) {
					errs.HandleError(err.Err, a.logger, a.tgbot, err.ChatID)
				}(err)
			case <-ctx.Done():
				a.logger.Debug("ctx is done")
				os.Exit(1)

			}
		}
	}()
}

func (a *App) initTgBot() {
	var err error
	a.tgbot, err = tgbotapi.NewBotAPI(a.config.TgBot.Token)
	if err != nil {
		a.errChan <- errs.TgError{
			Err: err,
		}
		return
	}

}

func (a *App) listenTgBot() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 1
	updates := a.tgbot.GetUpdatesChan(updateConfig)

	a.logger.Debug("tg_bot api started")

	for update := range updates {
		if update.Message != nil {
			if a.processingRegistrationUsers.IfExists(update.Message.Chat.ID) {
				go a.tgHandler.RegistrationHandler(update)

				continue
			}

			if a.processingLoginUsers.IfExists(update.Message.Chat.ID) {
				go a.tgHandler.LoginHandler(update)

				continue
			}
		}
		if update.Message == nil {
			if update.CallbackQuery != nil {
				switch {
				case strings.Contains(update.CallbackQuery.Data, "buy_data"):
					go a.tgHandler.BuyDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "sell_data"):
					go a.tgHandler.SellDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "view_data"):
					go a.tgHandler.ViewDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "characteristics_data"):
					go a.tgHandler.CharacteristicsDataButton(update)
					continue

				case strings.Contains(update.CallbackQuery.Data, "all_car_data"):
					go a.tgHandler.AllCarDataButton(update, a.loginUsers)
					continue

				}
			}

			continue
		}

		switch {
		case strings.Contains(update.Message.Text, "/start"):
			go a.tgHandler.StartHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/get-car"):
			go a.tgHandler.GetCarHandler(update, a.loginUsers)
			continue

		case strings.Contains(update.Message.Text, "/all-cars"):
			go a.tgHandler.AllCarsHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/get-user"):
			go a.tgHandler.GetUserHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/get-my-cars"):
			go a.tgHandler.GetMyCarsHandler(update, a.loginUsers)
			continue

		case strings.Contains(update.Message.Text, "/registration"):
			go a.tgHandler.RegistrationHandler(update)
			continue

		case strings.Contains(update.Message.Text, "/login"):
			go a.tgHandler.LoginHandler(update)
			continue

		}
	}
}

func (a *App) initLogger() {
	a.logger = log.Init()
}

func (a *App) initValidator() {
	a.validator = validator.New()
}

func (a *App) initRepos() {
	a.carsRepo = cars.NewRepository(a.logger, a.clientHTTP)
	a.usersRepo = user.NewRepository(a.logger, a.clientHTTP)

	a.logger.Debug("repos created")
}

func (a *App) initServices() {
	a.carService = car.NewService(a.carsRepo, a.logger)
	a.userService = user_service.NewService(a.usersRepo, a.logger)

	a.logger.Debug("services created")
}

func (a *App) initHandlers() {
	a.loginUsers = map[int64]string{}
	a.carHandler = car_handler.NewHandler(a.carService, a.tgbot)
	a.userHandler = user_handler.NewHandler(a.userService, a.tgbot, a.loginUsers, &a.processingRegistrationUsers, &a.processingLoginUsers)
	a.tgHandler = tg_handlers.NewHandler(a.tgbot, a.userHandler, a.carHandler, a.errChan, a.chatGPT)

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
