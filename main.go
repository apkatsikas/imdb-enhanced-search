package main

import (
	"log"
	"os"

	"github.com/apkatsikas/imdb-enhanced-search/client"
	"github.com/apkatsikas/imdb-enhanced-search/search"
)

var (
	basicsFile      = "title.basics.tsv.gz"
	ratingsFile     = "title.ratings.tsv.gz"
	imdbDataBaseUrl = "https://datasets.imdbws.com"
	imdbTitleUrl    = "https://www.imdb.com/title"
)

const (
	basicsFileEnv      = "IMDB_BASICS_FILE"
	ratingsFileEnv     = "IMDB_RATINGS_FILE"
	imdbDataBaseUrlEnv = "IMDB_DATA_BASE_URL"
	imdbTitleUrlEnv    = "IMDB_TITLE_URL"
)

func main() {
	log.Println("IMDB Enhanced Search")
	log.Println("====================")

	config := search.GetConfigFromUser()
	if config.DownloadData {
		if basicsEnv := os.Getenv(basicsFileEnv); basicsEnv != "" {
			basicsFile = basicsEnv
		}
		if ratingsEnv := os.Getenv(ratingsFileEnv); ratingsEnv != "" {
			ratingsFile = ratingsEnv
		}
		if dataEnv := os.Getenv(imdbDataBaseUrlEnv); dataEnv != "" {
			imdbDataBaseUrl = dataEnv
		}

		imdbClient, err := client.NewImdbClient(&client.ImdbConfig{
			BaseURL:     imdbDataBaseUrl,
			BasicsFile:  basicsFile,
			RatingsFile: ratingsFile,
		})
		if err != nil {
			log.Fatalf("Error getting IMDB client: %v", err)
		}
		if err := imdbClient.DownloadAndExtract(); err != nil {
			log.Fatalf("Error downloading and extracting IMDB data: %v", err)
		}
	}

	log.Println("Loading IMDB data...")
	basicsWithoutGzPath := basicsFile[:len(basicsFile)-3]
	movies, err := search.LoadMovies(basicsWithoutGzPath)
	if err != nil {
		log.Fatalf("Error loading movies: %v", err)
	}

	ratingsWithoutGzPath := ratingsFile[:len(ratingsFile)-3]
	ratings, err := search.LoadRatings(ratingsWithoutGzPath)
	if err != nil {
		log.Fatalf("Error loading ratings: %v", err)
	}

	log.Printf("Loaded %d movies and %d ratings", len(movies), len(ratings))

	results := search.FilterMovies(movies, ratings, config)

	log.Printf("Found %d movies matching your criteria\n", len(results))

	if titleEnv := os.Getenv(imdbTitleUrlEnv); titleEnv != "" {
		imdbTitleUrl = titleEnv
	}

	search.OpenMoviesInBrowser(imdbTitleUrl, results)
}
