package engine

import (
	"math"
	"math/rand"
)

// Unrated is a constant to be used with Title.rating to indicate that a title
// has no rating
var Unrated float64 = math.NaN()

// Title represents one line of IMDB title.basics.tsv data
type Title struct {
	tconst         string
	TitleType      string
	PrimaryTitle   string
	RuntimeMinutes string
	Genres         string
	Rating         float64
	RatingCount    int32
}

// ImdbTitleFilePaths represents local file paths corresponding to downloaded IMDB data
type ImdbTitleFilePaths struct {
	titleBasicsTsvGzPath  string
	titleRatingsTsvGzPath string
}

// RandTitle selects a title from an array of title uniformly at random
func RandTitle(list []Title) Title {
	return list[rand.Intn(len(list))]
}
