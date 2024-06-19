package lib

import (
	"context"
	"os"

	"github.com/paulmach/osm/osmpbf"
)

func OpenFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func CreateFileScanner(ctx context.Context, file *os.File) *osmpbf.Scanner {
	scanner := osmpbf.New(ctx, file, 3)
	return scanner
}
