package bootstrap

import (
	"bug_triage/internal/auth"
	"bug_triage/internal/cache"
	"bug_triage/internal/config"
	"bug_triage/internal/database"
	"bug_triage/internal/handler"
	"bug_triage/internal/kafka"
	"bug_triage/internal/pkg"
	"bug_triage/internal/repository"
	"bug_triage/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// AppDependencies holds all initialized application dependencies
type AppDependencies struct {
	DB           *sqlx.DB
	Redis        *redis.Client
	Kafka        *KafkaDependencies
	Repositories *RepositoryDependencies
	Services     *ServiceDependencies
	Auth         *AuthDependencies
	Handlers     *HandlerDependencies
	RateLimiter  *pkg.RateLimiter
	Logger       *zap.Logger
}

// KafkaDependencies holds Kafka-related dependencies
type KafkaDependencies struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

// RepositoryDependencies holds all repositories
type RepositoryDependencies struct {
	UserRepo repository.UserRepository
	BugRepo  repository.BugRepository
}

// ServiceDependencies holds all services
type ServiceDependencies struct {
	UserService *service.UserService
	BugService  *service.BugService
}

// AuthDependencies holds authentication-related dependencies
type AuthDependencies struct {
	PasswordManager *auth.PasswordManager
	JWTManager      *auth.JWTManager
}

// HandlerDependencies holds all handlers
type HandlerDependencies struct {
	AuthHandler *handler.AuthHandler
	BugHandler  *handler.BugHandler
}

// NewAppDependencies initializes all application dependencies
func NewAppDependencies(cfg *config.Config, log *zap.Logger) (*AppDependencies, error) {
	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DBUrl, log)
	if err != nil {
		return nil, err
	}

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, log)
	if err != nil {
		return nil, err
	}

	// Initialize Kafka
	kafkaProducer := kafka.NewProducerWithBrokers([]string{cfg.KafkaBroker}, log)

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepo(db)
	bugRepo := repository.NewPostgresBugRepo(db)

	// Initialize auth utilities
	passwordManager := auth.NewPasswordManager()
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Initialize services
	userService := service.NewUserService(userRepo, passwordManager, jwtManager)
	bugService := service.NewBugService(bugRepo, kafkaProducer, log)

	// Initialize rate limiter
	rateLimiter := pkg.NewRateLimiter(redisClient, log)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, log)
	bugHandler := handler.NewBugHandler(bugService, log)

	return &AppDependencies{
		DB:    db,
		Redis: redisClient,
		Kafka: &KafkaDependencies{
			Producer: kafkaProducer,
		},
		Repositories: &RepositoryDependencies{
			UserRepo: userRepo,
			BugRepo:  bugRepo,
		},
		Services: &ServiceDependencies{
			UserService: userService,
			BugService:  bugService,
		},
		Auth: &AuthDependencies{
			PasswordManager: passwordManager,
			JWTManager:      jwtManager,
		},
		Handlers: &HandlerDependencies{
			AuthHandler: authHandler,
			BugHandler:  bugHandler,
		},
		RateLimiter: rateLimiter,
		Logger:      log,
	}, nil
}

// Close closes all closeable dependencies
func (d *AppDependencies) Close() error {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.Kafka != nil && d.Kafka.Producer != nil {
		d.Kafka.Producer.Close()
	}
	if d.Kafka != nil && d.Kafka.Consumer != nil {
		d.Kafka.Consumer.Close()
	}
	return nil
}
