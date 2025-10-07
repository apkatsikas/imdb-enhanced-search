package search

import (
	"math"
	"os"
	"slices"
	"testing"
)

const (
	searchWorkersEnv = "IMDB_SEARCH_WORKERS"
)

// Test helpers
func intPtr(i int) *int {
	return &i
}

func createTestMovie(id, title string, isAdult bool, year, runtime int, genres []string) Movie {
	return Movie{
		Id:             id,
		PrimaryTitle:   title,
		isAdult:        isAdult,
		StartYear:      intPtr(year),
		runtimeMinutes: intPtr(runtime),
		Genres:         genres,
	}
}

func createTestRating(avgRating float64, numVotes int) rating {
	return rating{
		AverageRating: avgRating,
		NumVotes:      numVotes,
	}
}

func setupTestData() (map[string]Movie, map[string]rating) {
	os.Setenv(searchWorkersEnv, "")
	movies := map[string]Movie{
		"1": createTestMovie("1", "Action Hero", false, 2020, 120, []string{"Action"}),
		"2": createTestMovie("2", "Adult Drama", true, 2019, 90, []string{"Drama"}),
		"3": createTestMovie("3", "Thriller Night", false, 2021, 150, []string{"Action", "Thriller"}),
		"4": createTestMovie("4", "Comedy Show", false, 2018, 60, []string{"Comedy"}),
		"5": createTestMovie("5", "Action Pack", false, 2022, 130, []string{"Action"}),
		"6": createTestMovie("6", "Old Classic", false, 2005, 100, []string{"Drama"}),
		"7": createTestMovie("7", "Long Epic", false, 2020, 200, []string{"Drama"}),
		"8": createTestMovie("8", "Short Film", false, 2020, 45, []string{"Documentary"}),
	}

	ratings := map[string]rating{
		"1": createTestRating(8.5, 10000),
		"2": createTestRating(7.8, 8000),
		"3": createTestRating(9.0, 15000),
		"4": createTestRating(6.5, 3000),
		"5": createTestRating(8.0, 12000),
		"6": createTestRating(7.5, 20000),
		"7": createTestRating(8.2, 5000),
		"8": createTestRating(7.0, 2000),
	}

	return movies, ratings
}

func TestFilterMovies_AdultFilter(t *testing.T) {
	const resultCount = 7
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: true,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
		genres:       nil,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
			for _, movie := range results {
				if movie.isAdult {
					t.Errorf("%v is adult and should have been filtered", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_AdultFilterOff(t *testing.T) {
	const resultCount = 8
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: false,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
		genres:       nil,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
		})
	}
}

func TestFilterMovies_YearFilter(t *testing.T) {
	const resultCount = 6
	const minYear = 2018
	const maxYear = 2021
	movies, ratings := setupTestData()

	cfg := config{
		minYear:      minYear,
		maxYear:      maxYear,
		excludeAdult: false,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
		genres:       nil,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}

			for _, movie := range results {
				if *movie.StartYear < minYear {
					t.Errorf("%v start year is less than minimum year", movie.Id)
				}
				if *movie.StartYear > maxYear {
					t.Errorf("%v start year is greater than maximum year", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_RuntimeFilter(t *testing.T) {
	const resultCount = 5
	const minRuntime = 90
	const maxRuntime = 150
	movies, ratings := setupTestData()

	cfg := config{
		minRuntime:   minRuntime,
		maxRuntime:   maxRuntime,
		excludeAdult: false,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		genres:       nil,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}

			for _, movie := range results {
				runtime := *movie.runtimeMinutes
				if runtime < minRuntime || runtime > cfg.maxRuntime {
					t.Errorf("Movie %s has runtime %d, outside range", movie.Id, runtime)
				}
			}
		})
	}
}

func TestFilterMovies_RatingFilter(t *testing.T) {
	const resultCount = 3
	const maxVotes = 15_000
	const minVotes = 10_000
	const minRating = 8.0

	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: false,
		minYear:      0,
		maxYear:      99999999,
		minRuntime:   0,
		maxRuntime:   999999999,
		minRating:    minRating,
		minVotes:     minVotes,
		maxVotes:     maxVotes,
		genres:       nil,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)
			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
			for _, movie := range results {
				r := ratings[movie.Id]
				if r.AverageRating < minRating {
					t.Errorf("Movie %s has rating %.1f, below minimum", movie.Id, r.AverageRating)
				}
				if r.NumVotes < minVotes {
					t.Errorf("Movie %s has %d votes, below minimum", movie.Id, r.NumVotes)
				}
				if r.NumVotes > maxVotes {
					t.Errorf("Movie %s has %d votes, above maximum", movie.Id, r.NumVotes)
				}
			}
		})
	}
}

func TestFilterMovies_GenreFilter(t *testing.T) {
	const resultCount = 3
	const genre = "Action"
	movies, ratings := setupTestData()

	cfg := config{
		genres:       []string{genre},
		excludeAdult: false,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)
			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
			for _, movie := range results {
				if !slices.Contains(movie.Genres, genre) {
					t.Errorf("%v did not contain expected genre", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_GenreFilterMultiple(t *testing.T) {
	const resultCount = 6
	const genre1 = "Action"
	const genre2 = "Drama"
	movies, ratings := setupTestData()

	cfg := config{
		genres:       []string{genre1, genre2},
		excludeAdult: false,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)
			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
			for _, movie := range results {
				if !slices.Contains(movie.Genres, genre1) && !slices.Contains(movie.Genres, genre2) {
					t.Errorf("%v did not contain either of the expected genres", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_GenreFilterCaseSensitivity(t *testing.T) {
	const resultCount = 6
	const genre1 = "ACTION"
	const genre2 = "drama"
	movies, ratings := setupTestData()

	cfg := config{
		genres:       []string{genre1, genre2},
		excludeAdult: false,
		minYear:      0,
		maxYear:      9999,
		minRating:    0,
		minVotes:     0,
		maxVotes:     99999999999999,
		maxRuntime:   99999999999,
		minRuntime:   0,
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)
			if len(results) != resultCount {
				t.Errorf("Expected %v results, got %v", resultCount, len(results))
			}
			for _, movie := range results {
				if !slices.Contains(movie.Genres, "Action") && !slices.Contains(movie.Genres, "Drama") {
					t.Errorf("%v did not contain either of the expected genres", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_CombinedFilters(t *testing.T) {
	var expectedIDs = map[string]bool{"1": true, "5": true, "3": true}
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: true,
		minYear:      2019,
		maxYear:      2022,
		minRuntime:   100,
		maxRuntime:   180,
		minRating:    8.0,
		minVotes:     10000,
		maxVotes:     9999999,
		genres:       []string{"Action"},
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != len(expectedIDs) {
				t.Errorf("Expected %d movies passing all filters, got %d", len(expectedIDs), len(results))
			}

			for _, movie := range results {
				if !expectedIDs[movie.Id] {
					t.Errorf("Movie %s should not be in results", movie.Id)
				}
			}
		})
	}
}

func TestFilterMovies_NoResultsCase(t *testing.T) {
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: false,
		minYear:      2030,
		maxYear:      2040,
		minRuntime:   0,
		maxRuntime:   300,
		minRating:    0,
		minVotes:     0,
		genres:       []string{},
	}

	tests := []struct {
		name       string
		filterFunc func(map[string]Movie, map[string]rating, config) []Movie
	}{
		{"Sync", FilterMoviesSync},
		{"Async", FilterMovies},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := tt.filterFunc(movies, ratings, cfg)

			if len(results) != 0 {
				t.Errorf("Expected 0 results with impossible filter, got %d", len(results))
			}
		})
	}
}

func TestFilterMovies_MoreWorkersThanFilms(t *testing.T) {
	const resultCount = 8
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: false,
		minYear:      0,
		maxYear:      math.MaxInt,
		minRuntime:   0,
		maxRuntime:   math.MaxInt,
		minRating:    0,
		minVotes:     0,
		maxVotes:     math.MaxInt,
		genres:       nil,
	}

	os.Setenv(searchWorkersEnv, "9")

	results := FilterMovies(movies, ratings, cfg)

	if len(results) != resultCount {
		t.Errorf("Got %v results, expected %v", len(results), resultCount)
	}
}

func TestFilterMovies_OneWorker(t *testing.T) {
	const resultCount = 8
	movies, ratings := setupTestData()

	cfg := config{
		excludeAdult: false,
		minYear:      0,
		maxYear:      math.MaxInt,
		minRuntime:   0,
		maxRuntime:   math.MaxInt,
		minRating:    0,
		minVotes:     0,
		maxVotes:     math.MaxInt,
		genres:       nil,
	}

	os.Setenv(searchWorkersEnv, "1")

	results := FilterMovies(movies, ratings, cfg)

	if len(results) != resultCount {
		t.Errorf("Got %v results, expected %v", len(results), resultCount)
	}
}

func TestFilterMovies_ConsistencyWithSync(t *testing.T) {
	movies, ratings := setupTestData()

	testCases := []config{
		{excludeAdult: true, maxVotes: math.MaxInt, minYear: 2000, maxYear: 2025, minRuntime: 0, maxRuntime: 300, minRating: 0, minVotes: 0, genres: []string{}},
		{excludeAdult: false, maxVotes: math.MaxInt, minYear: 2018, maxYear: 2021, minRuntime: 90, maxRuntime: 150, minRating: 7.5, minVotes: 5000, genres: []string{"Action"}},
		{excludeAdult: true, maxVotes: math.MaxInt, minYear: 2019, maxYear: 2022, minRuntime: 100, maxRuntime: 180, minRating: 8.0, minVotes: 10000, genres: []string{"Action", "Drama"}},
	}

	for i, cfg := range testCases {
		syncResults := FilterMoviesSync(movies, ratings, cfg)
		asyncResults := FilterMovies(movies, ratings, cfg)

		if len(syncResults) == 0 || len(asyncResults) == 0 {
			t.Error("Expected at least 1 result for both sets, got 0")
		}

		if len(syncResults) != len(asyncResults) {
			t.Errorf("Test case %d: Sync returned %d movies, async returned %d", i, len(syncResults), len(asyncResults))
		}

		syncIDs := make(map[string]bool)
		for _, movie := range syncResults {
			syncIDs[movie.Id] = true
		}

		for _, movie := range asyncResults {
			if !syncIDs[movie.Id] {
				t.Errorf("Test case %d: Async result contains movie %s not in sync results", i, movie.Id)
			}
		}
	}
}

// Benchmark helpers
func generateTestMovies(n int) map[string]Movie {
	movies := make(map[string]Movie)
	genres := []string{"Action", "Drama", "Comedy", "Thriller", "Horror"}

	for i := 0; i < n; i++ {
		id := string(rune('a'+(i%26))) + string(rune('0'+(i/26)))
		title := "Movie " + id
		movies[id] = createTestMovie(
			id,
			title,
			i%5 == 0,
			2000+(i%25),
			60+(i%150),
			[]string{genres[i%len(genres)], genres[(i+1)%len(genres)]},
		)
	}

	return movies
}

func generateTestRatings(n int) map[string]rating {
	ratings := make(map[string]rating)

	for i := 0; i < n; i++ {
		id := string(rune('a'+(i%26))) + string(rune('0'+(i/26)))
		ratings[id] = createTestRating(
			5.0+float64(i%50)/10.0,
			1000+(i*100),
		)
	}

	return ratings
}

// Benchmarks
func BenchmarkFilterMoviesSync_Small(b *testing.B) {
	movies := generateTestMovies(100)
	ratings := generateTestRatings(100)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMoviesSync(movies, ratings, cfg)
	}
}

func BenchmarkFilterMovies_Small(b *testing.B) {
	movies := generateTestMovies(100)
	ratings := generateTestRatings(100)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMovies(movies, ratings, cfg)
	}
}

func BenchmarkFilterMoviesSync_Medium(b *testing.B) {
	movies := generateTestMovies(1000)
	ratings := generateTestRatings(1000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMoviesSync(movies, ratings, cfg)
	}
}

func BenchmarkFilterMovies_Medium(b *testing.B) {
	movies := generateTestMovies(1000)
	ratings := generateTestRatings(1000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMovies(movies, ratings, cfg)
	}
}

func BenchmarkFilterMoviesSync_Large(b *testing.B) {
	movies := generateTestMovies(10000)
	ratings := generateTestRatings(10000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMoviesSync(movies, ratings, cfg)
	}
}

func BenchmarkFilterMovies_Large(b *testing.B) {
	movies := generateTestMovies(10000)
	ratings := generateTestRatings(10000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMovies(movies, ratings, cfg)
	}
}

func BenchmarkFilterMoviesSync_VeryLarge(b *testing.B) {
	movies := generateTestMovies(100000)
	ratings := generateTestRatings(100000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMoviesSync(movies, ratings, cfg)
	}
}

func BenchmarkFilterMovies_VeryLarge(b *testing.B) {
	movies := generateTestMovies(100000)
	ratings := generateTestRatings(100000)
	cfg := config{
		excludeAdult: true,
		minYear:      2010,
		maxYear:      2024,
		minRuntime:   90,
		maxRuntime:   180,
		minRating:    7.0,
		minVotes:     5000,
		genres:       []string{"Action"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterMovies(movies, ratings, cfg)
	}
}
