package engine

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// ReadImdbTitleFiles parses the results of imdb_downloader into a format that
// Suggest can use.
func ReadImdbTitleFiles(titleFilePaths ImdbTitleFilePaths) map[string]Title {
	titles := make(map[string]Title)

	ratingsColumns := []string{"tconst", "averageRating", "numVotes"}
	processTsvGzFile(titleFilePaths.titleRatingsTsvGzPath, ratingsColumns, titles, addPartialTitleFromRating)

	basicsColumns := []string{"tconst", "titleType", "primaryTitle", "originalTitle", "isAdult", "startYear", "endYear", "runtimeMinutes", "genres"}
	processTsvGzFile(titleFilePaths.titleBasicsTsvGzPath, basicsColumns, titles, updatePrimaryTitleFromTitleBasicsRecord)

	return titles
}

// parses records from https://datasets.imdbws.com/title.basics.tsv.gz. Example of format:
//
// tconst	titleType	primaryTitle	originalTitle	isAdult	startYear	endYear	runtimeMinutes	genres
// tt0000001	short	Carmencita	Carmencita	0	1894	\N	1	Documentary,Short
func updatePrimaryTitleFromTitleBasicsRecord(record []string, titles map[string]Title) {
	tconst := record[0]
	if title, ok := titles[tconst]; ok {
		title.TitleType = record[1]
		title.PrimaryTitle = record[2]
		title.RuntimeMinutes = record[7]
		title.Genres = record[8]
		titles[tconst] = title
	}
}

// parses records from https://datasets.imdbws.com/title.ratings.tsv.gz. Example of format:
//
// tconst	averageRating	numVotes
// tt0000001	5.6	1668
func addPartialTitleFromRating(record []string, titles map[string]Title) {
	tconst := record[0]

	parsedRating, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	parsedRatingCount, err := strconv.ParseInt(record[2], 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	titles[tconst] = Title{
		tconst:         record[0],
		TitleType:      "",
		PrimaryTitle:   "",
		RuntimeMinutes: "",
		Genres:         "",
		Rating:         parsedRating,
		RatingCount:    int32(parsedRatingCount),
	}
}

type sanitizingReader struct {
	src io.Reader
}

func newSanitizingReader(src io.Reader) io.Reader {
	return sanitizingReader{src: src}
}

func (sr sanitizingReader) Read(p []byte) (int, error) {
	count, err := sr.src.Read(p)
	if err != nil {
		return count, err
	}
	for i := 0; i < count; i++ {
		if p[i] == '"' {
			p[i] = '\''
		}
	}
	return count, err
}

type titleProcessor = func([]string, map[string]Title)

func processTsvGzFile(filepath string, expectedColumns []string, titles map[string]Title, processor titleProcessor) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	sanitizedReader := newSanitizingReader(gr)

	cr := csv.NewReader(sanitizedReader)
	cr.Comma = '\t'
	cr.LazyQuotes = true
	header, err := cr.Read()
	if err == io.EOF {
		log.Fatal("Expected to find a tsv header, but instead got empty file")
	}
	if err != nil {
		log.Fatal(err)
	}

	if len(header) != len(expectedColumns) {
		log.Fatalf("Expected tsv header with %d columns but found %d\n", len(expectedColumns), len(header))
	}
	for index, expectedColumn := range expectedColumns {
		actualColumn := header[index]
		if actualColumn != expectedColumn {
			log.Fatalf("Expected tsv column %d to be %s but found %s\n", index+1, expectedColumn, actualColumn)
		}
	}

	for i := 0; i < 100000; i++ {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		processor(record, titles)
	}
}
