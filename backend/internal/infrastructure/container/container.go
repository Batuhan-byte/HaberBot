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
	UserRepo    port.UserRepository
	CommentRepo port.CommentRepository
	Fetchers    []port.ContentFetcher
	AIProcessor port.AIProcessor

	// Use Cases
	FetchArticlesUC     *usecase.FetchArticlesUseCase
	ProcessArticlesUC   *usecase.ProcessArticlesUseCase
	ListArticlesUC      *usecase.ListArticlesUseCase
	GetArticleUC        *usecase.GetArticleUseCase
	ManageTopicsUC      *usecase.ManageTopicsUseCase
	ManageArticlesUC    *usecase.ManageArticlesUseCase
	DailyPipelineUC     *usecase.DailyPipelineUseCase
	SearchArticlesUC    *usecase.SearchArticlesUseCase
	SummarizeArticleUC  *usecase.SummarizeArticleUseCase
	AuthRegisterUC      *usecase.AuthRegisterUseCase
	AuthLoginUC         *usecase.AuthLoginUseCase
	AuthRefreshUC       *usecase.AuthRefreshUseCase
	AuthLogoutUC        *usecase.AuthLogoutUseCase
	CreateCommentUC     *usecase.CreateCommentUseCase
	ListCommentsUC      *usecase.ListCommentsUseCase

	// Handlers
	ArticleHandler *handler.ArticleHandler
	TopicHandler   *handler.TopicHandler
	AdminHandler   *handler.AdminHandler
	HealthHandler  *handler.HealthHandler
	AuthHandler    *handler.AuthHandler
	CommentHandler *handler.CommentHandler

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
	userRepo := repository.NewPostgresUserRepo(pool)
	commentRepo := repository.NewPostgresCommentRepo(pool)

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
	fetchArticlesUC := usecase.NewFetchArticlesUseCase(fetchers, articleRepo, cfg.PendingLimitPerTopic)
	processArticlesUC := usecase.NewProcessArticlesUseCase(aiProcessor, articleRepo)
	listArticlesUC := usecase.NewListArticlesUseCase(articleRepo)
	getArticleUC := usecase.NewGetArticleUseCase(articleRepo)
	manageTopicsUC := usecase.NewManageTopicsUseCase(topicRepo)
	manageArticlesUC := usecase.NewManageArticlesUseCase(articleRepo)
	dailyPipelineUC := usecase.NewDailyPipelineUseCase(fetchArticlesUC, processArticlesUC, manageTopicsUC)
	searchArticlesUC := usecase.NewSearchArticlesUseCase(articleRepo)
	summarizeArticleUC := usecase.NewSummarizeArticleUseCase(aiProcessor, articleRepo)

	authRegisterUC := usecase.NewAuthRegisterUseCase(userRepo)
	authLoginUC := usecase.NewAuthLoginUseCase(userRepo, cfg.JWTSecret, cfg.JWTAccessTTLMinutes, cfg.JWTRefreshTTLDays)
	authRefreshUC := usecase.NewAuthRefreshUseCase(userRepo, cfg.JWTSecret, cfg.JWTAccessTTLMinutes, cfg.JWTRefreshTTLDays)
	authLogoutUC := usecase.NewAuthLogoutUseCase(userRepo)

	createCommentUC := usecase.NewCreateCommentUseCase(commentRepo, articleRepo)
	listCommentsUC := usecase.NewListCommentsUseCase(commentRepo)

	// 5. Handlers
	articleHandler := handler.NewArticleHandler(listArticlesUC, getArticleUC, searchArticlesUC, summarizeArticleUC, cfg.AdminAPIKey)
	topicHandler := handler.NewTopicHandler(manageTopicsUC, listArticlesUC)
	adminHandler := handler.NewAdminHandler(manageTopicsUC, manageArticlesUC, dailyPipelineUC, cfg.AdminAPIKey)
	healthHandler := handler.NewHealthHandler()

	authHandler := handler.NewAuthHandler(authRegisterUC, authLoginUC, authRefreshUC, authLogoutUC, cfg.JWTRefreshTTLDays)
	commentHandler := handler.NewCommentHandler(createCommentUC, listCommentsUC)

	// 6. Scheduler
	cronScheduler := scheduler.NewScheduler(dailyPipelineUC)

	return &Container{
		Config:             cfg,
		DBPool:             pool,
		ArticleRepo:        articleRepo,
		TopicRepo:          topicRepo,
		UserRepo:           userRepo,
		CommentRepo:        commentRepo,
		Fetchers:           fetchers,
		AIProcessor:        aiProcessor,
		FetchArticlesUC:    fetchArticlesUC,
		ProcessArticlesUC:  processArticlesUC,
		ListArticlesUC:     listArticlesUC,
		GetArticleUC:       getArticleUC,
		ManageTopicsUC:     manageTopicsUC,
		ManageArticlesUC:   manageArticlesUC,
		DailyPipelineUC:    dailyPipelineUC,
		SearchArticlesUC:   searchArticlesUC,
		SummarizeArticleUC: summarizeArticleUC,
		AuthRegisterUC:     authRegisterUC,
		AuthLoginUC:        authLoginUC,
		AuthRefreshUC:      authRefreshUC,
		AuthLogoutUC:       authLogoutUC,
		CreateCommentUC:    createCommentUC,
		ListCommentsUC:     listCommentsUC,
		ArticleHandler:     articleHandler,
		TopicHandler:       topicHandler,
		AdminHandler:       adminHandler,
		HealthHandler:      healthHandler,
		AuthHandler:        authHandler,
		CommentHandler:     commentHandler,
		Scheduler:          cronScheduler,
	}, nil
}
