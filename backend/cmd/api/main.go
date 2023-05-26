package main

import (
	redis2 "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"backend/config"
	"backend/internal/server"
	"backend/pkg/db/postgres"
	"backend/pkg/db/redis"
	"backend/pkg/logger"
	"backend/pkg/utils"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// @title Go Example REST API
// @version 1.0
// @description Example Golang REST API
// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru
// @BasePath /api/v1
func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer func(psqlDB *sqlx.DB) {
		err := psqlDB.Close()
		if err != nil {

		}
	}(psqlDB)

	redisClient := redis.NewRedisClient(cfg)
	defer func(redisClient *redis2.Client) {
		err := redisClient.Close()
		if err != nil {

		}
	}(redisClient)
	appLogger.Info("Redis connected")

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {

		}
	}(closer)
	appLogger.Info("Opentracing connected")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
