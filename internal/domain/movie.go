package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Plot             string             `bson:"plot" json:"plot"`
	Genres           []string           `bson:"genres" json:"genres"`
	Runtime          int                `bson:"runtime" json:"runtime"`
	Cast             []string           `bson:"cast" json:"cast"`
	NumMflixComments int                `bson:"num_mflix_comments" json:"num_mflix_comments"`
	Title            string             `bson:"title" json:"title"`
	Fullplot         string             `bson:"fullplot" json:"fullplot"`
	Countries        []string           `bson:"countries" json:"countries"`
	Languages        []string           `bson:"languages" json:"languages"`
	Released         primitive.DateTime `bson:"released" json:"released"`
	Directors        []string           `bson:"directors" json:"directors"`
	Rated            string             `bson:"rated" json:"rated"`
	Awards           Awards             `bson:"awards" json:"awards"`
	LastUpdated      string             `bson:"lastupdated" json:"lastupdated"`
	Year             int                `bson:"year" json:"year"`
	IMDB             IMDB               `bson:"imdb" json:"imdb"`
	Type             string             `bson:"type" json:"type"`
	Tomatoes         Tomatoes           `bson:"tomatoes" json:"tomatoes"`
	Poster           string             `bson:"poster" json:"poster"`
}

type Awards struct {
	Wins        int    `bson:"wins" json:"wins"`
	Nominations int    `bson:"nominations" json:"nominations"`
	Text        string `bson:"text" json:"text"`
}

type IMDB struct {
	Rating float64 `bson:"rating" json:"rating"`
	Votes  int     `bson:"votes" json:"votes"`
	ID     int     `bson:"id" json:"id"`
}

type Tomatoes struct {
	Viewer      Viewer             `bson:"viewer" json:"viewer"`
	Fresh       int                `bson:"fresh" json:"fresh"`
	Critic      Viewer             `bson:"critic" json:"critic"`
	Rotten      int                `bson:"rotten" json:"rotten"`
	LastUpdated primitive.DateTime `bson:"lastUpdated" json:"lastUpdated"`
	DVD         primitive.DateTime `bson:"dvd" json:"dvd"`
}

type Viewer struct {
	Rating     float64 `bson:"rating" json:"rating"`
	NumReviews int     `bson:"numReviews" json:"numReviews"`
	Meter      int     `bson:"meter" json:"meter"`
}

type MovieRepository interface {
	GetMovie(ctx context.Context, id primitive.ObjectID) (*Movie, error)
	GetMovies(ctx context.Context, titleFilter string, page, limit int) ([]Movie, error)
}
