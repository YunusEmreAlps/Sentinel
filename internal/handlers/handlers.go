package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sentinel/config"
	"sentinel/pkg/metric"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	API_PREFIX = "/sentinel"
	RN_PREFIX  = "cld:::sentinel:::"
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

type RespondJson struct {
	Status  bool        `json:"status"`
	Intent  string      `json:"intent"`
	Message interface{} `json:"message"`
}

func respondJson(ctx *gin.Context, code int, intent string, message interface{}, err error) {
	if err == nil {
		ctx.JSON(code, RespondJson{
			Status:  true,
			Intent:  intent,
			Message: message,
		})
	} else {
		ctx.JSON(code, RespondJson{
			Status:  false,
			Intent:  intent,
			Message: err.Error(),
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

	// Set up the Gin router
	/*r.LoadHTMLGlob("pkg/templates/*")

	// Serve the welcome page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	})*/

	// List all certificates
	v1.GET("/certificates", func(ctx *gin.Context) {
		code, data, err := bs.ListCertificates(ctx)
		respondJson(ctx, code, RN_PREFIX+"/certificates", data, err)
	})

	// Certificate information for a specific domain
	v1.GET("/certificates/:domain", func(ctx *gin.Context) {
		code, data, err := bs.GetCertificateInfo(ctx)
		respondJson(ctx, code, RN_PREFIX+"/certificates/:domain", data, err)
	})

	// Get all certificates from the utility/data.go file
	v1.GET("/certificates/all", func(ctx *gin.Context) {
		code, data, err := bs.GetAllExpirations(ctx)
		respondJson(ctx, code, RN_PREFIX+"/check/all", data, err)
	})

	// Health check
	health.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Sentinel Service " + config.C.App.Version + " is running on port " + config.C.App.Port + ".",
		})
	})
}
