package engine

import (
	"log"
)

// SuggestOptions is the options type for the Suggest function
// ^^^ This is why required comments are idiotic
type SuggestOptions struct {
	MinRating      float64
	MinRatingCount int32
	SeenTitles     []string
}

// Suggest uses IMDB data to suggest the title of a movie to watch
func Suggest(options SuggestOptions) Title {
	var titleFiles = DownloadImdbTitleData()
	var allTitles = ReadImdbTitleFiles(titleFiles)
	var movieTitles = filterNonMovies(allTitles)
	// log.Printf("Parsed %d rated titles from IMDB", len(movieTitles))
	var candidateTitles = filterByRating(movieTitles, options.MinRating, options.MinRatingCount)
	// log.Printf("%d titles after filtering by rating", len(wellRatedTitles))
	var unseenCandidateTitles = filterAlreadySeen(candidateTitles, options.SeenTitles)
	// log.Printf("%d titles after filtering already-seen titles", len(candidateTitles))
	log.Printf("you still have %d (of %d) movies to watch!", len(unseenCandidateTitles), len(candidateTitles))
	if len(unseenCandidateTitles) == 0 {
		log.Fatal("you've seen it all!")
	}
	var suggestion = RandTitle(unseenCandidateTitles)
	return suggestion
}

func filterNonMovies(titles map[string]Title) []Title {
	var filteredTitles = make([]Title, 0, len(titles))
	for _, title := range titles {
		if title.TitleType == "movie" {
			filteredTitles = append(filteredTitles, title)
		}
	}
	return filteredTitles
}

func filterByRating(titles []Title, minRating float64, minRatingCount int32) []Title {
	var filteredTitles = make([]Title, 0, len(titles))
	for _, title := range titles {
		if title.Rating >= minRating && title.RatingCount >= minRatingCount {
			filteredTitles = append(filteredTitles, title)
		}
	}
	return filteredTitles
}

func filterAlreadySeen(titles []Title, alreadySeenTitleNames []string) []Title {
	var filteredTitles = make([]Title, 0, len(titles))

title_loop:
	for _, title := range titles {
		for _, alreadySeenName := range alreadySeenTitleNames {
			if title.PrimaryTitle == alreadySeenName {
				continue title_loop
			}
		}
		filteredTitles = append(filteredTitles, title)
	}

	return filteredTitles
}
