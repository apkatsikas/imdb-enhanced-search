package search

type Movie struct {
	Id             string
	titleType      string
	PrimaryTitle   string
	originalTitle  string
	isAdult        bool
	StartYear      *int // Pointer to allow nil for missing data
	endYear        *int
	runtimeMinutes *int // Pointer to allow nil for missing data
	Genres         []string
}
