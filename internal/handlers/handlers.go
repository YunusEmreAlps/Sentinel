package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sentinel/config"
	"sentinel/pkg/metric"

	"sentinel/internal/models"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	API_PREFIX = "/sentinel"
)

// Sentinel is a struct for auth core
type Sentinel struct {
	inAppCache   *persistence.InMemoryStore
	cache        *redis.Client
	cacheContext context.Context
	db           *gorm.DB
	// s3sess       *session.Session
}

func NewSentinelService(
	inAppCache *persistence.InMemoryStore,
	cache *redis.Client,
	cacheContext context.Context,
	db *gorm.DB,
	// s3sess *session.Session,
) *Sentinel {
	return &Sentinel{
		inAppCache:   inAppCache,
		cache:        cache,
		cacheContext: cacheContext,
		db:           db,
		// s3sess:       s3sess,
	}
}

// respondJson standardizes API responses with consistent format
func respondJson(ctx *gin.Context, code int, path string, data interface{}, err error) {
	if err == nil {
		ctx.JSON(code, models.APIResponse{
			Success: true,
			Path:    path,
			Data:    data,
		})
	} else {
		ctx.JSON(code, models.APIResponse{
			Success: false,
			Path:    path,
			Data:    err.Error(),
		})
	}
}

func (bs *Sentinel) InitRouter(r *gin.Engine) {
	// Prometheus metrics
	metrics, err := metric.CreateMetrics(config.C.Metric.Url, config.C.Metric.Service)
	if err != nil {
		fmt.Println("INIT: Cannot create metrics")
	}
	fmt.Println("INIT: Application configuration success.", metrics)

	// -- my service routes (group)
	v1 := r.Group(API_PREFIX)
	health := v1.Group("/health")

	// List all certificates
	v1.GET("/domains", func(ctx *gin.Context) {
		code, data, err := bs.ListCertificates(ctx)
		respondJson(ctx, code, API_PREFIX+"/domains", data, err)
	})

	// Certificate information for a specific domain
	v1.GET("/certificates/:domain", func(ctx *gin.Context) {
		code, data, err := bs.GetCertificateInfo(ctx)
		respondJson(ctx, code, API_PREFIX+"/certificates/:domain", data, err)
	})

	// Get all certificates from the utility/data.go file
	v1.GET("/certificates/scan", func(ctx *gin.Context) {
		code, data, err := bs.GetAllExpirations(ctx)
		respondJson(ctx, code, API_PREFIX+"/certificates/scan", data, err)
	})

	// Health check
	health.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Sentinel Service " + config.C.App.Version + " is running on port " + config.C.App.Port + ".",
		})
	})
}
