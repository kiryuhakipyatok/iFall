package app

import (
	"iFall/internal/bot"
	"iFall/internal/client"
	"iFall/internal/config"
	"iFall/internal/delivery/handlers"
	"iFall/internal/delivery/routes"
	"iFall/internal/domain/repositories"
	"iFall/internal/domain/services"
	"iFall/internal/email"
	"iFall/internal/scheduler"
	"iFall/pkg/logger"
	"iFall/pkg/server"
	"iFall/pkg/storage"
	"iFall/pkg/validator"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))
	logger := logger.NewLogger(cfg.App)
	logger.Info("config loaded successfully")

	validator := validator.NewValidator()

	storage := storage.MustConnect(cfg.Storage)
	logger.Info("connected to postgres successfully")
	defer func() {
		storage.MustClose()
		logger.Info("postgres closed successfully")
	}()

	server := server.NewServer(cfg.Server, cfg.App)
	logger.Info("server created successfully")
	defer func() {
		server.MustClose()
		logger.Info("server closed successfully")
	}()

	client := client.NewClient(cfg.ApiClient)

	smtpAuth := smtp.PlainAuth("", cfg.Email.Address, cfg.Email.Password, cfg.Email.SmtpAddress)
	emailSender := email.NewEmailSender(smtpAuth, cfg.Email)

	userRepository := repositories.NewUserRepository(storage)
	iphoneRepository := repositories.NewIPhoneRepository(storage)

	bot := bot.NewTelegramBot(cfg.TelegramBot, logger, userRepository)
	logger.Info("bot created successfully")
	bot.StoreChatId()
	defer func() {
		bot.Stop()
		logger.Info("bot stopped successfully")
	}()

	userService := services.NewUserService(userRepository, logger)

	iphoneService := services.NewIPhoneService(iphoneRepository, client, logger, emailSender, cfg.IPhones)
	iphoneReportService := services.NewIPhoneReportService(iphoneService, userRepository, logger, bot, emailSender, cfg.IPhones)

	userHandler := handlers.NewUsersHandler(userService, validator)

	routesSetup := routes.NewRoutesSetup(server.App, userHandler)
	routesSetup.SetupRoutes()

	scheduler := scheduler.NewScheduler(iphoneService, iphoneReportService, logger, cfg.Scheduler)
	scheduler.Start()
	defer func() {
		scheduler.Stop()
	}()

	go func() {
		logger.Info("bot started successfully")
		bot.Start()
	}()

	go func() {
		server.Start()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Info("app shutting down...")
}
