package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultTimeoutMinutes = 2

type ImdbClient struct {
	baseURL     string
	httpClient  *http.Client
	basicsFile  string
	ratingsFile string
}

type ImdbConfig struct {
	BaseURL     string
	BasicsFile  string
	RatingsFile string
	Timeout     time.Duration
}

func NewImdbClient(cfg *ImdbConfig) (*ImdbClient, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	cfg.withDefaults()

	return &ImdbClient{
		baseURL:     cfg.BaseURL,
		basicsFile:  cfg.BasicsFile,
		ratingsFile: cfg.RatingsFile,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

func (c *ImdbClient) DownloadAndExtract() error {
	log.Println("Downloading datasets from", c.baseURL)

	for _, downloadPath := range []string{c.basicsFile, c.ratingsFile} {
		url := c.baseURL + "/" + downloadPath
		log.Println("Downloading", downloadPath)

		if err := c.downloadFile(url, downloadPath); err != nil {
			return fmt.Errorf("failed to download %s: %w", downloadPath, err)
		}

		log.Println("Extracting", downloadPath)
		extractWithoutGzPath := downloadPath[:len(downloadPath)-3]

		if err := extractGzip(downloadPath, extractWithoutGzPath); err != nil {
			return fmt.Errorf("failed to extract %s: %w", downloadPath, err)
		}

		if err := os.Remove(downloadPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", downloadPath, err)
		}
	}

	log.Println("Done downloading IMDB data")
	return nil
}

func (c *ImdbConfig) validate() error {
	return requireArgs(map[string]string{
		"BaseURL":     c.BaseURL,
		"BasicsFile":  c.BasicsFile,
		"RatingsFile": c.RatingsFile,
	})
}

func (c *ImdbConfig) withDefaults() *ImdbConfig {
	if c.Timeout == 0 {
		c.Timeout = defaultTimeoutMinutes * time.Minute
	}
	return c
}

func requireArgs(args map[string]string) error {
	var missing []string
	for name, value := range args {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required arguments missing: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (c *ImdbClient) downloadFile(url, filepath string) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status downloading file %v: %s", filepath, resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
