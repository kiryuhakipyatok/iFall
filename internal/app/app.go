package app

import (
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

	"github.com/gofiber/fiber/v2/log"
)

func Run() {
	cfg := config.MustLoad(os.Getenv("CONFIG_PATH"))
	logger := logger.NewLogger(cfg.App)
	logger.Info("config loaded successfully")

	validator := validator.NewValidator()

	storage := storage.MustConnect(cfg.Storage)
	logger.Info("connected to postgres successfully")
	defer func() {
		storage.Close()
		logger.Info("postgres closed successfully")
	}()

	server := server.NewServer(cfg.Server, cfg.App)
	logger.Info("server created successfully")
	defer func() {
		server.MustClose()
		logger.Info("server closed successfully")
	}()

	client := client.NewClient(cfg.ApiClient)

	userRepository := repositories.NewUserRepository(storage)
	iphoneRepository := repositories.NewIPhoneRepository(storage)

	userService := services.NewUserService(userRepository, logger)
	iphoneService := services.NewIPhoneService(iphoneRepository, client, logger)

	userHandler := handlers.NewUsersHandler(userService, validator)

	routesSetup := routes.NewRoutesSetup(server.App, userHandler)
	routesSetup.SetupRoutes()

	smtpAuth := smtp.PlainAuth("", cfg.Email.Address, cfg.Email.Password, cfg.Email.SmtpAddress)
	emailSendler := email.NewEmailSender(smtpAuth, cfg.Email)

	scheduler := scheduler.NewScheduler(iphoneService, userRepository, emailSendler, logger, cfg.Scheduler, cfg.IPhones)
	scheduler.Start()
	defer func() {
		scheduler.Stop()
	}()

	go func() {
		server.Start()
		logger.Info("server started successfully")
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info("app shutting down...")
}
