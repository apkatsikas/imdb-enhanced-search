package search

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	defaultMinYear      = 0
	defaultMaxYear      = math.MaxInt
	defaultMinRating    = 0.0
	defaultMinVotes     = 0
	defaultMaxVotes     = math.MaxInt
	defaultMinRuntime   = 0
	defaultMaxRuntime   = math.MaxInt
	defaultExcludeAdult = false
)

var defaultGenres = []string{}

func GetConfigFromUser() config {
	reader := bufio.NewReader(os.Stdin)

	config := config{
		DownloadData: false,
		minYear:      defaultMinYear,
		maxYear:      defaultMaxYear,
		minRating:    defaultMinRating,
		minVotes:     defaultMinVotes,
		maxVotes:     defaultMaxVotes,
		minRuntime:   defaultMinRuntime,
		maxRuntime:   defaultMaxRuntime,
		excludeAdult: defaultExcludeAdult,
		genres:       defaultGenres,
	}

	fmt.Println("Download fresh dataset from IMDB? y=yes:")
	downloadData := strings.ToLower(readLine(reader))
	if downloadData == "y" || downloadData == "yes" {
		config.DownloadData = true
	}

	fmt.Print("\nEnter minimum year: ")
	setConfigInt(reader, func(i int) {
		config.minYear = i
	})

	fmt.Print("Enter maximum year: ")
	setConfigInt(reader, func(i int) {
		config.maxYear = i
	})

	fmt.Print("Enter minimum run time: ")
	setConfigInt(reader, func(i int) {
		config.minRuntime = i
	})

	fmt.Print("Enter maximum run time: ")
	setConfigInt(reader, func(i int) {
		config.maxRuntime = i
	})

	fmt.Print("Enter minimum rating: ")
	if rating := readLine(reader); rating != "" {
		if r, err := strconv.ParseFloat(rating, 64); err == nil {
			config.minRating = r
		} else {
			log.Println("Invalid value provided, ignoring input")
		}
	}

	fmt.Print("Enter minimum votes: ")
	setConfigInt(reader, func(i int) {
		config.minVotes = i
	})

	fmt.Print("Enter maximum votes: ")
	setConfigInt(reader, func(i int) {
		config.maxVotes = i
	})

	fmt.Print("Enter genres (comma-separated, e.g., Action,Drama): ")
	if genres := readLine(reader); genres != "" {
		config.genres = strings.Split(genres, ",")
		for i := range config.genres {
			config.genres[i] = strings.TrimSpace(config.genres[i])
		}
	}

	return config
}

func OpenMoviesInBrowser(imdbTitleUrl string, results []Movie) {
	scanner := bufio.NewScanner(os.Stdin)
	for len(results) > 0 {
		movie := results[0]

		fmt.Print("Press Enter to open in browser (or 'q' to quit): ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if strings.ToLower(input) == "q" {
			log.Println("Quitting")
			break
		}

		url := fmt.Sprintf("%s/%s/", imdbTitleUrl, movie.Id)
		if output, err := openBrowser(url); err != nil {
			log.Fatalln("Error opening browser", err, string(output))
		}

		results = results[1:]
	}

	if len(results) == 0 {
		log.Println("No more results")
		os.Exit(0)
	}
}

func openBrowser(url string) ([]byte, error) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		if isWSL() {
			cmd = "wslview"
		} else {
			cmd = "xdg-open"
		}
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).CombinedOutput()
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func setConfigInt(reader *bufio.Reader, configFunc func(int)) {
	if input := readLine(reader); input != "" {
		if v, err := strconv.Atoi(input); err == nil {
			configFunc(v)
			return
		}
		log.Println("Invalid value provided, ignoring input")
	}
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
