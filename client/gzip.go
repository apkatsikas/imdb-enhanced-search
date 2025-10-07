package client

import (
	"compress/gzip"
	"io"
	"os"
)

func extractGzip(gzipPath, outputPath string) error {
	gzipFile, err := os.Open(gzipPath)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzipReader)
	return err
}
