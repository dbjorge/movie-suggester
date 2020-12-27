package engine

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

// DownloadImdbTitleData downloads current copies of IMDB's title data to a
// temp directory and outputs the resulting file paths.
func DownloadImdbTitleData() ImdbTitleFilePaths {
	return ImdbTitleFilePaths{
		titleBasicsTsvGzPath:  refreshTempFile("https://datasets.imdbws.com/title.basics.tsv.gz"),
		titleRatingsTsvGzPath: refreshTempFile("https://datasets.imdbws.com/title.ratings.tsv.gz"),
	}
}

func needsRefresh(filePath string) bool {
	cutoffTime := time.Now().Add(-1 * time.Hour * 24)
	stats, oops := os.Stat(filePath)
	if oops != nil {
		if e, ok := oops.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return true
		}
		log.Fatal(oops)
	}
	return stats.ModTime().Before(cutoffTime)
}

func refreshTempFile(rawURLString string) string {
	// Build fileName from fullPath
	parsedURL, err := url.Parse(rawURLString)
	if err != nil {
		log.Fatal(err)
	}
	urlPath := parsedURL.Path
	segments := strings.Split(urlPath, "/")
	fileName := "movie_suggester_" + segments[len(segments)-1]
	filePath := path.Join(os.TempDir(), fileName)

	if needsRefresh(filePath) {
		// log.Printf("Refreshing data file at %s", filePath)
		downloadFile(rawURLString, filePath)
	} else {
		// log.Printf("Reusing existing data file at %s", filePath)
	}

	return filePath
}

func downloadFile(rawURLString string, filePath string) {
	// Create blank file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(rawURLString)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
}
