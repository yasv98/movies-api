package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yasv98/movies-api/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	movieUsecase *service.MovieService
}

func NewMovieHandler(movieUsecase *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieUsecase: movieUsecase,
	}
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.movieUsecase.GetMovie(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	titleFilter := c.Query("title")
	movies, err := h.movieUsecase.GetMovies(c.Request.Context(), titleFilter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}
