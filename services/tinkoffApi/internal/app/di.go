package app

import (
	"context"
	"log/slog"
	"tinkoffApi/internal/configs"
	"tinkoffApi/internal/handlers"
	redisClient "tinkoffApi/internal/repository/redis"
	"tinkoffApi/internal/service"
	loggeradapter "tinkoffApi/lib/logger/loggerAdapter"

	"github.com/gladinov/cryptotoken"
	"github.com/redis/go-redis/v9"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type diContainer struct {
	logger            *slog.Logger
	configs           *configs.Configs
	loggerAdapter     investgo.Logger
	analyticsService  service.AnalyticsService
	portfolioService  service.PortfolioService
	instrumentService service.InstrumentService
	service           *service.Service
	tokenCrypter      *cryptotoken.TokenCrypter // TODO: Создать интерфейс
	redis             *redis.Client
	handlers          *handlers.Handlers // TODO: Хэндлер слой должен быть интерфейсом
}

func newDIContainer(logger *slog.Logger, configs *configs.Configs) *diContainer {
	return &diContainer{
		logger:  logger,
		configs: configs,
	}
}

func (d *diContainer) LoggerAdapter() investgo.Logger {
	if d.loggerAdapter == nil {
		d.logger.Info("initialize logger adapter")
		d.loggerAdapter = loggeradapter.NewLoggerAdapter(d.logger)
	}

	return d.loggerAdapter
}

func (d *diContainer) AnalyticsService() service.AnalyticsService {
	if d.analyticsService == nil {
		d.logger.Info("initialize analytics service")
		d.analyticsService = service.NewAnalyticsServiceClient(d.configs.TinkoffApiConfig, d.LoggerAdapter())
	}

	return d.analyticsService
}

func (d *diContainer) PortfolioService() service.PortfolioService {
	if d.portfolioService == nil {
		d.logger.Info("initialize portfolio service")
		d.portfolioService = service.NewPortfolioServiceClient(d.configs.TinkoffApiConfig, d.LoggerAdapter())
	}

	return d.portfolioService
}

func (d *diContainer) InstrumentService() service.InstrumentService {
	if d.instrumentService == nil {
		d.logger.Info("initialize instrument service")
		d.instrumentService = service.NewInstrumentServiceClient(d.configs.TinkoffApiConfig, d.LoggerAdapter())
	}

	return d.instrumentService
}

func (d *diContainer) Service() *service.Service {
	if d.service == nil {
		d.logger.Info("initialize service client")
		d.service = service.NewService(
			d.AnalyticsService(),
			d.PortfolioService(),
			d.InstrumentService())
	}

	return d.service
}

func (d *diContainer) TokenCrypter() *cryptotoken.TokenCrypter {
	if d.tokenCrypter == nil {
		d.logger.Info("initialize token crypter")
		d.tokenCrypter = cryptotoken.NewTokenCrypter(d.configs.Config.Key)
	}

	return d.tokenCrypter
}

func (d *diContainer) Redis() *redis.Client {
	if d.redis == nil {
		d.logger.Info("initialize redis", slog.String("address", d.configs.Config.RedisHTTPServer.GetAddress()))

		redis, err := redisClient.NewClient(context.Background(), d.configs.Config)
		if err != nil {
			d.logger.Error("haven't connect with redis", slog.String("error", err.Error()))
		}

		d.redis = redis
	}

	return d.redis
}

func (d *diContainer) Handlers() *handlers.Handlers {
	if d.handlers == nil {
		d.logger.Info("initialize handlers")
		d.handlers = handlers.NewHandlers(d.logger, d.Service(), d.TokenCrypter(), d.Redis())
	}

	return d.handlers
}
