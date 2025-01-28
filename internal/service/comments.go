package service

import (
	"context"

	"github.com/yasv98/movies-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	commentRepo domain.CommentRepository
}

func NewCommentService(commentRepo domain.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (c *CommentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	return c.commentRepo.Create(ctx, comment)
}

func (c *CommentService) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	return c.commentRepo.Update(ctx, comment)
}

func (c *CommentService) DeleteComment(ctx context.Context, movieID, commentID primitive.ObjectID) error {
	return c.commentRepo.Delete(ctx, movieID, commentID)
}

func (c *CommentService) GetMovieComment(ctx context.Context, movieID, commentID primitive.ObjectID) (*domain.Comment, error) {
	return c.commentRepo.GetMovieComment(ctx, movieID, commentID)
}

func (c *CommentService) GetMovieComments(ctx context.Context, movieID primitive.ObjectID) ([]domain.Comment, error) {
	return c.commentRepo.GetMovieComments(ctx, movieID)
}
