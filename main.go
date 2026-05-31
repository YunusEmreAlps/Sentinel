package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	docs "sentinel/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"sentinel/config"
	"sentinel/internal/handlers"
	"sentinel/pkg/constants"
	"sentinel/pkg/db/postgres"
	rediscfg "sentinel/pkg/db/redis"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"
)

// @title Sentinel API
// @description: This is a simple API for checking the expiration date of the certificates of the given domain names.
// @version 1.0.0
// @schemes http https

// @contact.name   Yunus Emre Alpu
// @contact.url    https://yunusemrealpu.netlify.app
// @contact.email  YunusAlpu@icloud.com

// @BasePath /sentinel

var isConfigSuccess = false

// var equals string = strings.Repeat("=", 50)

// APP_NAME = "localhost:8080/sentinel/"
const (
	APP_NAME = "sentinel"
)

func main() {
	var dbConn *gorm.DB
	var cacheConn *redis.Client
	var cacheContext context.Context
	var inAppCache *persistence.InMemoryStore

	mode := config.C.App.Mode
	port := config.C.App.Port

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(gin.DefaultErrorWriter))

	// If cache is active, create cache connection
	if config.C.Cache.Active {
		inAppCache = rediscfg.NewInAppCacheStore(time.Duration(config.C.App.Expire) * time.Second)
		cacheConn, cacheContext = rediscfg.NewRedisCacheConnection(config.C.Cache.Url)
	}

	// If db is active, create db connection
	if config.C.DB.Active {
		dbConn = postgres.NewPostgresDB(config.C.DB.Url)
		utils.GlobalDataService = utils.NewDataService(dbConn)
		logger.CLogger.Info("GlobalDataService initialized with database connection")
	}

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: config.C.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.C.Jaeger.LogSpans,
			LocalAgentHostPort: config.C.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Fatal("cannot create tracer", err)
	}

	// create application service
	sentinelsvc := handlers.NewSentinelService(
		inAppCache,
		cacheConn,
		cacheContext,
		dbConn,
	)

	// check env and set gin mode
	setApplicationMode(mode, router)
	sentinelsvc.InitRouter(router)

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	logger.CLogger.Info("Opentracing connected")

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    constants.ServerReadTimeout,
		WriteTimeout:   constants.ServerWriteTimeout,
		IdleTimeout:    constants.ServerIdleTimeout,
		MaxHeaderBytes: constants.ServerMaxHeaderBytes,
	}

	// Run server in goroutine
	go func() {
		logger.CLogger.Info("INIT: Application " + APP_NAME + " started on port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.CLogger.Fatal("INIT: Server failed to start: ", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.CLogger.Info("SHUTDOWN: Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.CLogger.Fatal("SHUTDOWN: Server forced to shutdown: ", err)
	}

	// Close database connection if active
	if config.C.DB.Active && dbConn != nil {
		sqlDB, err := dbConn.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// Close Redis connection if active
	if config.C.Cache.Active && cacheConn != nil {
		cacheConn.Close()
	}

	logger.CLogger.Info("SHUTDOWN: Server exited gracefully")
}

// Initialize Application
func init() {
	isConfigSuccess = configureApplication()
	if !isConfigSuccess {
		logger.CLogger.Error("INIT: Application configuration failed.")
		os.Exit(1)
	} else {
		logger.CLogger.Info("INIT: Application configuration success.")
	}
}

// Configure Application with config file
func configureApplication() bool {
	dir, err := os.Getwd()
	if err != nil {
		logger.CLogger.Error("INIT: Cannot get current working directory os.Getwd()")
		return false
	} else {
		config.ReadConfig(dir)
		return true
	}
}

// Set Application Mode
func setApplicationMode(md string, router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	if md == "prod" || md == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.Use(gin.Logger())
		// logger time gonna be in UTC+3 (Asia/Istanbul)
		gin.SetMode(gin.DebugMode)
	}

	// check env and set swagger
	if !(md == "prod" || md == "production") {
		docs.SwaggerInfo.BasePath = handlers.API_PREFIX
		// Endpoint for swagger: http://localhost:{targetPort}/sentinel/swagger/index.html
		router.GET(handlers.API_PREFIX+"/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
