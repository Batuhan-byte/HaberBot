package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"haberbot/internal/adapter/fetcher"
	"haberbot/internal/adapter/gateway"
	"haberbot/internal/adapter/handler"
	"haberbot/internal/adapter/repository"
	"haberbot/internal/domain/port"
	"haberbot/internal/infrastructure/config"
	"haberbot/internal/infrastructure/database"
	"haberbot/internal/infrastructure/scheduler"
	"haberbot/internal/usecase"
)

// Container holds and manages all runtime dependencies for the HaberBot application.
type Container struct {
	Config      config.Config
	DBPool      *pgxpool.Pool
	ArticleRepo port.ArticleRepository
	TopicRepo   port.TopicRepository
	Fetchers    []port.ContentFetcher
	AIProcessor port.AIProcessor

	// Use Cases
	FetchArticlesUC     *usecase.FetchArticlesUseCase
	ProcessArticlesUC   *usecase.ProcessArticlesUseCase
	ListArticlesUC      *usecase.ListArticlesUseCase
	GetArticleUC        *usecase.GetArticleUseCase
	ManageTopicsUC      *usecase.ManageTopicsUseCase
	DailyPipelineUC     *usecase.DailyPipelineUseCase
	SearchArticlesUC    *usecase.SearchArticlesUseCase
	SummarizeArticleUC  *usecase.SummarizeArticleUseCase

	// Handlers
	ArticleHandler *handler.ArticleHandler
	TopicHandler   *handler.TopicHandler
	AdminHandler   *handler.AdminHandler
	HealthHandler  *handler.HealthHandler

	// Scheduler
	Scheduler *scheduler.Scheduler
}

// NewContainer initializes and wires all application components together.
func NewContainer(ctx context.Context, cfg config.Config) (*Container, error) {
	// 1. DB connection pool
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("container DB setup: %w", err)
	}

	// 2. Repositories
	articleRepo := repository.NewPostgresArticleRepo(pool)
	topicRepo := repository.NewPostgresTopicRepo(pool)

	// 3. Gateways & Fetchers
	hnFetcher := fetcher.NewHackerNewsFetcher()

	// Default feeds for RSS fetcher.
	// You can also add more standard tech feeds here.
	defaultFeeds := []string{
		"https://news.ycombinator.com/rss",
		"https://dev.to/feed",
		"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
	}
	rssFetcher := fetcher.NewRSSFetcher(defaultFeeds)

	fetchers := []port.ContentFetcher{hnFetcher, rssFetcher}

	var aiProcessor port.AIProcessor
	if strings.HasPrefix(cfg.AIAPIKey, "AIzaSy") || (cfg.AIAPIKey == "" && strings.HasPrefix(cfg.GeminiAPIKey, "AIzaSy")) {
		apiKey := cfg.GeminiAPIKey
		if apiKey == "" {
			apiKey = cfg.AIAPIKey
		}
		aiProcessor = gateway.NewGeminiProcessor(apiKey)
	} else {
		aiProcessor = gateway.NewOpenAIProcessor(cfg.AIAPIKey, cfg.AIBaseURL, cfg.AIModel)
	}

	// 4. Use Cases
	fetchArticlesUC := usecase.NewFetchArticlesUseCase(fetchers, articleRepo)
	processArticlesUC := usecase.NewProcessArticlesUseCase(aiProcessor, articleRepo)
	listArticlesUC := usecase.NewListArticlesUseCase(articleRepo)
	getArticleUC := usecase.NewGetArticleUseCase(articleRepo)
	manageTopicsUC := usecase.NewManageTopicsUseCase(topicRepo)
	dailyPipelineUC := usecase.NewDailyPipelineUseCase(fetchArticlesUC, processArticlesUC, manageTopicsUC)
	searchArticlesUC := usecase.NewSearchArticlesUseCase(articleRepo)
	summarizeArticleUC := usecase.NewSummarizeArticleUseCase(aiProcessor, articleRepo)

	// 5. Handlers
	articleHandler := handler.NewArticleHandler(listArticlesUC, getArticleUC, searchArticlesUC, summarizeArticleUC)
	topicHandler := handler.NewTopicHandler(manageTopicsUC, listArticlesUC)
	adminHandler := handler.NewAdminHandler(manageTopicsUC, dailyPipelineUC, cfg.AdminAPIKey)
	healthHandler := handler.NewHealthHandler()

	// 6. Scheduler
	cronScheduler := scheduler.NewScheduler(dailyPipelineUC)

	return &Container{
		Config:            cfg,
		DBPool:            pool,
		ArticleRepo:       articleRepo,
		TopicRepo:         topicRepo,
		Fetchers:          fetchers,
		AIProcessor:       aiProcessor,
		FetchArticlesUC:   fetchArticlesUC,
		ProcessArticlesUC: processArticlesUC,
		ListArticlesUC:    listArticlesUC,
		GetArticleUC:      getArticleUC,
		ManageTopicsUC:    manageTopicsUC,
		DailyPipelineUC:   dailyPipelineUC,
		SearchArticlesUC:  searchArticlesUC,
		SummarizeArticleUC: summarizeArticleUC,
		ArticleHandler:    articleHandler,
		TopicHandler:      topicHandler,
		AdminHandler:      adminHandler,
		HealthHandler:     healthHandler,
		Scheduler:         cronScheduler,
	}, nil
}
