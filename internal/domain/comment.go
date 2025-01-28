package domain

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrCommentNotFoundForMovie = errors.New("comment not found for movie")

type Comment struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	MovieID primitive.ObjectID `bson:"movie_id" json:"movie_id"`
	Name    string             `bson:"name" json:"name"`
	Email   string             `bson:"email" json:"email"`
	Text    string             `bson:"text" json:"text"`
	Date    primitive.DateTime `bson:"date" json:"date"`
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, movieID, commentID primitive.ObjectID) error
	GetMovieComment(ctx context.Context, movieID, commentID primitive.ObjectID) (*Comment, error)
	GetMovieComments(ctx context.Context, movieID primitive.ObjectID) ([]Comment, error)
}
