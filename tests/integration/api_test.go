package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"github.com/yasv98/movies-api/internal/delivery/http/handler"
	"github.com/yasv98/movies-api/internal/delivery/http/routes"
	"github.com/yasv98/movies-api/internal/repository/mongodb"
	"github.com/yasv98/movies-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// The following integration tests are setup to test
// enpoints against a pre-seeded test Mongo DB.
//
// Start the test DB docker container before running
// or run via make command `make integration-tests`.
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	s.server.Close()
	_ = s.client.Disconnect(ctx)
}

type IntegrationTestSuite struct {
	suite.Suite
	client *mongo.Client
	db     *mongo.Database
	app    *application
	server *httptest.Server
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.client, s.db = connectDatabase(context.Background())
	s.app = newApp(s.db)
	s.server = httptest.NewServer(s.app.Router)
}

// TODO: Add more throughough testing to below integration test suite:
// - Assert JSON body response on valid requests
// - Assert expected fields are present and have correct types
// - Compare response data against test database records
// - Add test cases for all errors and invalid requests

const validMovieID = "573a1390f29313caabcd4eaf"
const invalidMovieID = "12345"
const missingMovieID = "573a1390f29313caabcd4133"

func (s *IntegrationTestSuite) TestGetMovie_Valid() {
	apitest.New("Get movie with valid ID").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + validMovieID).
		Expect(s.T()).
		Status(http.StatusOK).
		End()
}

func (s *IntegrationTestSuite) TestGetMovie_Invalid() {
	apitest.New("Get movie with invalid ID").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + invalidMovieID).
		Expect(s.T()).
		Body(`{"error": "Invalid movie ID format"}`).
		Status(http.StatusBadRequest).
		End()

	apitest.New("Get movie ID not found").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + missingMovieID).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Body(`{"error": "Movie not found"}`).
		End()
}

func (s *IntegrationTestSuite) TestGetMovies_Valid() {
	apitest.New("Get movies without any parameters").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Expect(s.T()).
		Status(http.StatusOK).
		End()

	apitest.New("Get movies filtered by title").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Query("title", "Blacksmith").
		Expect(s.T()).
		Status(http.StatusOK).
		End()

	apitest.New("Get movies with pagination").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Query("page", "2").
		Query("limit", "10").
		Expect(s.T()).
		Status(http.StatusOK).
		End()

	apitest.New("Get movies with both title filter and pagination").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Query("title", "Matrix").
		Query("page", "1").
		Query("limit", "5").
		Expect(s.T()).
		Status(http.StatusOK).
		End()
}

func (s *IntegrationTestSuite) TestGetMovies_Invalid() {
	apitest.New("Get movies with invalid page parameter").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Query("page", "qwefqwefqwf").
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Body(`{"error": "invalid page parameter"}`).
		End()

	apitest.New("Get movies with invalid limit parameter").
		Handler(s.app.Router).
		Get("/api/v1/movies").
		Query("limit", "invalid").
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Body(`{"error": "invalid limit parameter"}`).
		End()
}

const validCommentID = "5a9427648b0beebeb6957a22"
const invalidCommentID = "12345"
const missingCommentID = "5a9427648b0beebeb69579cd"

func (s *IntegrationTestSuite) TestGetMovieComment_Valid() {
	apitest.New("Get comment with valid movie and comment ID").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + validMovieID + "/comments/" + validCommentID).
		Expect(s.T()).
		Status(http.StatusOK).
		End()
}

func (s *IntegrationTestSuite) TestGetMovieComment_Invalid() {
	apitest.New("Get comment with invalid comment ID format").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + validMovieID + "/comments/" + invalidCommentID).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Body(`{"error": "invalid comment ID format"}`).
		End()

	apitest.New("Get comment with non-existent comment ID").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + validMovieID + "/comments/" + missingCommentID).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Body(`{"error": "Comment not found"}`).
		End()
}

func (s *IntegrationTestSuite) TestGetMovieComments_Valid() {
	apitest.New("Get movie comments without pagination").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + validMovieID + "/comments").
		Expect(s.T()).
		Status(http.StatusOK).
		End()

	apitest.New("Get movie comments with pagination").
		Handler(s.app.Router).
		Get("/api/v1/movies/"+validMovieID+"/comments").
		Query("page", "1").
		Query("limit", "10").
		Expect(s.T()).
		Status(http.StatusOK).
		End()
}

func (s *IntegrationTestSuite) TestGetMovieComments_Invalid() {
	invalidMovieID := "12345"
	apitest.New("Get comments with invalid movie ID format").
		Handler(s.app.Router).
		Get("/api/v1/movies/" + invalidMovieID + "/comments").
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Body(`{"error": "invalid movie ID format"}`).
		End()
}

func (s *IntegrationTestSuite) TestCreateComment_Valid() {
	movieID := "573a1390f29313caabcd4135"
	comment := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
		"text":  "Great movie!",
	}

	apitest.New().
		Handler(s.app.Router).
		Post("/api/v1/movies/" + movieID + "/comments").
		JSON(comment).
		Expect(s.T()).
		Status(http.StatusCreated).
		End()
}

func (s *IntegrationTestSuite) TestCreateComment_Invalid() {
	invalidMovieID := "12345"
	validComment := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
		"text":  "Great movie!",
	}

	apitest.New("Create comment with invalid movie ID format").
		Handler(s.app.Router).
		Post("/api/v1/movies/" + invalidMovieID + "/comments").
		JSON(validComment).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Body(`{"error": "invalid movie ID format"}`).
		End()
}

func (s *IntegrationTestSuite) TestUpdateComment() {
	movieID := "573a1390f29313caabcd4b1b"
	commentID := "5a9427648b0beebeb6957a23"
	update := map[string]string{
		"name":  "John Doe Updated",
		"email": "john.updated@example.com",
		"text":  "Updated comment",
	}

	apitest.New().
		Handler(s.app.Router).
		Put("/api/v1/movies/" + movieID + "/comments/" + commentID).
		JSON(update).
		Expect(s.T()).
		Status(http.StatusOK).
		End()
}

// TODO: Use test setup and teardown to handle populating and cleaning up database.
func (s *IntegrationTestSuite) TestDeleteComment() {
	// First create comment to make sure test re-runs pass.
	movieID := "573a1390f29313caabcd6399"
	comment := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
		"text":  "Great movie!",
	}

	var resp struct {
		CommentID string `json:"id"`
	}

	apitest.New().
		Handler(s.app.Router).
		Post("/api/v1/movies/" + movieID + "/comments").
		JSON(comment).
		Expect(s.T()).
		Status(http.StatusCreated).
		End().
		JSON(&resp)

	// Now test deleting it.
	apitest.New().
		Handler(s.app.Router).
		Delete("/api/v1/movies/" + movieID + "/comments/" + resp.CommentID).
		Expect(s.T()).
		Status(http.StatusNoContent).
		End()
}

// TODO: Use mongo DB test container and seed with deterministic data.
func connectDatabase(ctx context.Context) (*mongo.Client, *mongo.Database) {
	// Connect to existing Docker container.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://host.docker.internal:27017")) // TODO: Pass test parameters in a config file.
	if err != nil {
		panic(err)
	}

	// Connect to pre-seeded database in Docker container.
	return client, client.Database("sample_mflix") // TODO: Pass test parameters in a config file.
}

type application struct {
	Router *gin.Engine
}

func newApp(db *mongo.Database) *application {
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

	return &application{Router: router}
}
