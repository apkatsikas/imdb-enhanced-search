package search

type config struct {
	DownloadData bool
	minYear      int
	maxYear      int
	minRating    float64
	minVotes     int
	maxVotes     int
	maxRuntime   int
	minRuntime   int
	genres       []string
	excludeAdult bool
}
