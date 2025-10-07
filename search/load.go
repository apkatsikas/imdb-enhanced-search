package search

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const tabComma = '\t'

func LoadMovies(filename string) (map[string]Movie, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = tabComma
	reader.LazyQuotes = true

	movies := make(map[string]Movie)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	requiredCols := []string{"tconst", "titleType", "primaryTitle", "originalTitle",
		"isAdult", "startYear", "endYear", "runtimeMinutes", "genres"}
	for _, col := range requiredCols {
		if _, exists := colIndex[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in file", col)
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Omitting record due to error on loadMovies:", err)
			continue
		}

		if len(record) < len(header) {
			log.Println("loadMovies: len(record) < len(header):", err)
			continue
		}

		movie := Movie{
			Id:            record[colIndex["tconst"]],
			titleType:     record[colIndex["titleType"]],
			PrimaryTitle:  record[colIndex["primaryTitle"]],
			originalTitle: record[colIndex["originalTitle"]],
			isAdult:       record[colIndex["isAdult"]] == "1",
		}

		if record[colIndex["startYear"]] != "\\N" {
			year, _ := strconv.Atoi(record[colIndex["startYear"]])
			movie.StartYear = &year
		}

		if record[colIndex["endYear"]] != "\\N" {
			year, _ := strconv.Atoi(record[colIndex["endYear"]])
			movie.endYear = &year
		}

		if record[colIndex["runtimeMinutes"]] != "\\N" {
			runtime, _ := strconv.Atoi(record[colIndex["runtimeMinutes"]])
			movie.runtimeMinutes = &runtime
		}

		if record[colIndex["genres"]] != "\\N" {
			movie.Genres = strings.Split(record[colIndex["genres"]], ",")
		}

		if validMovie(movie.titleType) {
			movies[movie.Id] = movie
		}
	}

	return movies, nil
}

func LoadRatings(filename string) (map[string]rating, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = tabComma

	ratings := make(map[string]rating)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	requiredCols := []string{"tconst", "averageRating", "numVotes"}
	for _, col := range requiredCols {
		if _, exists := colIndex[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in file", col)
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Omitting record due to error on loadRatings:", err)
			continue
		}

		if len(record) < len(header) {
			log.Println("loadRatings: len(record) < len(header):", err)
			continue
		}

		avgRating, _ := strconv.ParseFloat(record[colIndex["averageRating"]], 64)
		numVotes, _ := strconv.Atoi(record[colIndex["numVotes"]])

		ratings[record[colIndex["tconst"]]] = rating{
			id:            record[colIndex["tconst"]],
			AverageRating: avgRating,
			NumVotes:      numVotes,
		}
	}

	return ratings, nil
}

func validMovie(titleType string) bool {
	return titleType == "movie" || titleType == "short" || titleType == "tvMovie" || titleType == "tvShort"
}
