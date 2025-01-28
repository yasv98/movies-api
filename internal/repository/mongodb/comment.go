package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/yasv98/movies-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type commentRepository struct {
	db *mongo.Database
}

func NewCommentRepository(db *mongo.Database) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	comment.ID = primitive.NewObjectID()
	comment.Date = primitive.NewDateTimeFromTime(time.Now())

	_, err := r.db.Collection("comments").InsertOne(ctx, comment)
	return err
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	update := bson.M{
		"$set": bson.M{
			"name":  comment.Name,
			"email": comment.Email,
			"text":  comment.Text,
			"date":  time.Now(),
		},
	}

	result, err := r.db.Collection("comments").UpdateOne(
		ctx,
		bson.M{
			"_id":      comment.ID,
			"movie_id": comment.MovieID,
		},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrCommentNotFoundForMovie
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, movieID, commentID primitive.ObjectID) error {
	result, err := r.db.Collection("comments").DeleteOne(ctx, bson.M{
		"_id":      commentID,
		"movie_id": movieID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrCommentNotFoundForMovie
	}

	// Decrement the movie's comment count as it has been deleted.
	update := bson.M{
		"$inc": bson.M{
			"num_mflix_comments": -1,
		},
	}
	_, err = r.db.Collection("movies").UpdateOne(
		ctx,
		bson.M{"_id": movieID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update movie comment count: %w", err)
	}

	return nil
}

func (r *commentRepository) GetMovieComment(ctx context.Context, movieID, commentID primitive.ObjectID) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.Collection("comments").FindOne(ctx, bson.M{
		"_id":      commentID,
		"movie_id": movieID,
	}).Decode(&comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) GetMovieComments(ctx context.Context, movieID primitive.ObjectID) ([]domain.Comment, error) {
	cursor, err := r.db.Collection("comments").Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, fmt.Errorf("failed to find comments: %w", err)
	}
	defer cursor.Close(ctx)

	var comments []domain.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}
