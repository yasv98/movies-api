package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yasv98/movies-api/internal/delivery/http/handler"
)

func SetupRoutes(
	r *gin.Engine,
	movieHandler *handler.MovieHandler,
	commentHandler *handler.CommentHandler,
) {
	api := r.Group("/api/v1")
	{
		// Movie routes.
		api.GET("/movies/:movieId", movieHandler.GetMovie)
		api.GET("/movies", movieHandler.GetMovies)

		// Comment routes.
		api.GET("/movies/:movieId/comments/:commentId", commentHandler.GetMovieComment)
		api.GET("/movies/:movieId/comments", commentHandler.GetMovieComments)
		api.POST("/movies/:movieId/comments", commentHandler.CreateComment)
		api.PUT("/movies/:movieId/comments/:commentId", commentHandler.UpdateComment)
		api.DELETE("/movies/:movieId/comments/:commentId", commentHandler.DeleteComment)
	}
}
