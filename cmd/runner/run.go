package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yasv98/movies-api/internal/config"
	"github.com/yasv98/movies-api/internal/delivery/http/handler"
	"github.com/yasv98/movies-api/internal/delivery/http/routes"
	"github.com/yasv98/movies-api/internal/repository/mongodb"
	"github.com/yasv98/movies-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context, configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := initializeMongoDB(ctx, cfg.MonogoDB.URI)
	if err != nil {
		return fmt.Errorf("initialize mongo db: %w", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("error disconnecting mongo client: %v", err)
		}
	}()

	db := client.Database(cfg.MonogoDB.Database)

	// Repository.
	movieRepo := mongodb.NewMovieRepository(db)
	commentRepo := mongodb.NewCommentRepository(db)

	// Service.
	movieUsecase := service.NewMovieService(movieRepo)
	commentUsecase := service.NewCommentService(commentRepo)

	// Handler.
	movieHandler := handler.NewMovieHandler(movieUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	// Router.
	router := gin.Default()
	routes.SetupRoutes(router, movieHandler, commentHandler)

	return router.Run(":" + cfg.Port)
}

func initializeMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
