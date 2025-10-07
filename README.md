# imdb-enhanced-search

A partially vibe-coded port of my old Python version of [imdb-enhanced-search](https://github.com/Vampire-Computer-People/imdb-enhanced-search), re-built in Go.

Ever been on the hunt for a cool movie to watch but sick of browsing whatever is on Netflix? Want to learn how to hack data and write personalized queries custom built to your unique taste? 

Thanks to the [Internet Movie Database](https://www.imdb.com/) (IMDB) for the data.

For WSL users, make sure [wslu/wslview](https://github.com/wslutilities/wslu) is installed.

## Running

Grab a release from the releases page and follow the prompts after running the executable.

## Building from source

Run `make build`.

## Environment variables

The following environment variables can be used to overwrite some hard-coded values:

* `IMDB_BASICS_FILE` - defaults to `title.basics.tsv.gz`
* `IMDB_RATINGS_FILE` - defaults to `title.ratings.tsv.gz`
* `IMDB_DATA_BASE_URL` - defaults to `https://datasets.imdbws.com`
* `IMDB_TITLE_URL` - defaults to `https://www.imdb.com/title`

`IMDB_SEARCH_WORKERS` defaults to `runtime.NumCPU()`

## Future improvements

* Make the randomization of results configurable - currently results are always randomized
* Allow genres to be configurable to filter as "all" or "any" - currently, the filter works as "any"
* Make valid movie types configurable
* Add a configurable option sanitize quotes, to fix some problematic IMDB data, could be accomplished like:

```
	func sanitizeQuotes(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			// Remove all quote characters — like Excel does
			line = strings.ReplaceAll(line, `"`, "")
			fmt.Fprintln(pw, line)
		}
		pw.CloseWithError(scanner.Err())
	}()
	return pr
}
```
