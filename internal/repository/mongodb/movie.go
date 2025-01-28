package mongodb

import (
	"context"

	"github.com/yasv98/movies-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type movieRepository struct {
	db *mongo.Database
}

func NewMovieRepository(db *mongo.Database) domain.MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetMovie(ctx context.Context, id primitive.ObjectID) (*domain.Movie, error) {
	var movie domain.Movie
	if err := r.db.Collection("movies").FindOne(ctx, bson.M{"_id": id}).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *movieRepository) GetMovies(ctx context.Context, titleFilter string, page, limit int) ([]domain.Movie, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	filter := bson.M{}
	if titleFilter != "" {
		filter = bson.M{
			"title": primitive.Regex{
				Pattern: titleFilter,
				Options: "i",
			},
		}
	}

	cursor, err := r.db.Collection("movies").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []domain.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
