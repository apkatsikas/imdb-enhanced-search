package search

import (
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// FilterMoviesSync filters movies synchronously
func FilterMoviesSync(movies map[string]Movie, ratings map[string]rating, config config) []Movie {
	movieSlice := mapToSlice(movies)
	filtered := filterMovieSlice(movieSlice, ratings, config)
	randomizeResults(filtered)
	return filtered
}

// FilterMovies filters movies concurrently using worker pool
func FilterMovies(movies map[string]Movie, ratings map[string]rating, config config) []Movie {
	movieSlice := mapToSlice(movies)

	resultsChan := make(chan Movie, len(movieSlice))

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	altWorkers := os.Getenv("IMDB_SEARCH_WORKERS")
	if altWorkers != "" {
		if num, err := strconv.Atoi(altWorkers); err == nil {
			numWorkers = num
		}
	}
	chunks := partitionSlice(movieSlice, numWorkers)

	for _, chunk := range chunks {
		wg.Add(1)
		go func(movies []Movie) {
			defer wg.Done()
			filtered := filterMovieSlice(movies, ratings, config)
			for _, movie := range filtered {
				resultsChan <- movie
			}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	results := collectFromChannel(resultsChan)
	randomizeResults(results)
	return results
}

func filterMovieSlice(movies []Movie, ratings map[string]rating, cfg config) []Movie {
	results := make([]Movie, 0, len(movies))

	for _, movie := range movies {
		rating, hasRating := ratings[movie.Id]

		if shouldIncludeMovie(movie, rating, hasRating, cfg) {
			results = append(results, movie)
		}
	}

	return results
}

func shouldIncludeMovie(movie Movie, rating rating, hasRating bool, cfg config) bool {
	return passesAdultFilter(movie, cfg) &&
		passesYearFilter(movie, cfg) &&
		passesRuntimeFilter(movie, cfg) &&
		passesRatingFilter(rating, hasRating, cfg) &&
		passesGenreFilter(movie, cfg)
}

func passesAdultFilter(movie Movie, cfg config) bool {
	return !cfg.excludeAdult || !movie.isAdult
}

func passesYearFilter(movie Movie, cfg config) bool {
	if movie.StartYear == nil {
		return false
	}

	year := *movie.StartYear
	return year >= cfg.minYear && year <= cfg.maxYear
}

func passesRuntimeFilter(movie Movie, cfg config) bool {
	if movie.runtimeMinutes == nil {
		return false
	}

	runtime := *movie.runtimeMinutes
	return runtime >= cfg.minRuntime && runtime <= cfg.maxRuntime
}

func passesRatingFilter(rating rating, hasRating bool, cfg config) bool {
	if !hasRating {
		return false
	}

	return rating.AverageRating >= cfg.minRating &&
		rating.NumVotes >= cfg.minVotes && rating.NumVotes <= cfg.maxVotes
}

func passesGenreFilter(movie Movie, cfg config) bool {
	if len(cfg.genres) == 0 {
		return true
	}

	return hasGenre(movie.Genres, cfg.genres)
}

func randomizeResults(movies []Movie) {
	for i := range movies {
		j := rand.IntN(i + 1)
		movies[i], movies[j] = movies[j], movies[i]
	}
}

func mapToSlice(movies map[string]Movie) []Movie {
	slice := make([]Movie, 0, len(movies))
	for _, movie := range movies {
		slice = append(slice, movie)
	}
	return slice
}

func partitionSlice(movies []Movie, numWorkers int) [][]Movie {
	if len(movies) == 0 {
		return nil
	}

	chunkSize := (len(movies) + numWorkers - 1) / numWorkers
	chunks := make([][]Movie, 0, numWorkers)

	for i := 0; i < len(movies); i += chunkSize {
		end := i + chunkSize
		if end > len(movies) {
			end = len(movies)
		}
		chunks = append(chunks, movies[i:end])
	}

	return chunks
}

func collectFromChannel(ch chan Movie) []Movie {
	results := make([]Movie, 0)
	for movie := range ch {
		results = append(results, movie)
	}
	return results
}

func hasGenre(movieGenres, filterGenres []string) bool {
	genreSet := make(map[string]bool)
	for _, g := range movieGenres {
		genreSet[g] = true
	}

	for _, g := range filterGenres {
		if genreSet[normalizeGenre(g)] {
			return true
		}
	}
	return false
}

func normalizeGenre(genre string) string {
	if genre == "" {
		return ""
	}
	lower := strings.ToLower(genre)
	return strings.ToUpper(string(lower[0])) + lower[1:]
}
