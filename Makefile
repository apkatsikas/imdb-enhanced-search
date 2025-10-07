build:
	go build -o imdb-enhanced-search main.go

run: build
	./imdb-enhanced-search

test:
	go test -p 1 -race -coverprofile coverage.out ./...

check-formatting:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@if grep -qi microsoft /proc/version; then \
		explorer.exe coverage.html; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		open coverage.html; \
	else \
		xdg-open coverage.html; \
	fi

install-staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

check-build:
	go build -o imdb-enhanced-search -race .
