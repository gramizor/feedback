package app

import (
	"context"

	"rest-apishka/internal/config"
	"rest-apishka/internal/dsn"
	"rest-apishka/internal/http/feedback"
	"rest-apishka/internal/http/repository"
	"rest-apishka/internal/http/usecase"
)

// Application представляет основное приложение.
type Application struct {
    Config    *config.Config
    Repository *repository.Repository
	UseCase    *usecase.UseCase
	Handler    *feedback.Handler
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
    // Инициализируйте конфигурацию
    cfg, err := config.NewConfig(ctx)
    if err != nil {
        return nil, err
    }

    // Инициализируйте подключение к базе данных (DB)
    repo, err := repository.New(dsn.FromEnv())
    if err != nil {
        return nil, err
    }
    uc := usecase.NewUseCase(repo)
    h := feedback.NewHandler(uc)
    // Инициализируйте и настройте объект Application
    app := &Application{
        Config: cfg,
        Repository: repo,
        UseCase: uc,
        Handler: h,
    }

    return app, nil
}

