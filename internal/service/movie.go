package service

import (
	"context"

	"github.com/yasv98/movies-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService struct {
	movieRepo domain.MovieRepository
}

func NewMovieService(movieRepo domain.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (u *MovieService) GetMovie(ctx context.Context, id primitive.ObjectID) (*domain.Movie, error) {
	return u.movieRepo.GetMovie(ctx, id)
}

func (u *MovieService) GetMovies(ctx context.Context, titleFilter string, page, limit int) ([]domain.Movie, error) {
	return u.movieRepo.GetMovies(ctx, titleFilter, page, limit)
}
