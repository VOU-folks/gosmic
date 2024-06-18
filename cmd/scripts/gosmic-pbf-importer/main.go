package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/paulmach/osm/osmpbf"
	"os"
	"path/filepath"

	. "gosmic/internal/config"
)

func main() {
	var err error

	configFile := flag.String("config", "config.yaml", "path to the config file. Example: api -config /full/path/to/config.yaml")
	flag.Parse()

	config, err := GetConfigService(*configFile)
	if err != nil {
		panic(err)
	}

	storageConfig := config.GetStorage()
	osmConfig := config.GetOsmConfig()

	pbfFileName := osmConfig.Sources.PBF.FileName
	pbfStorage := storageConfig.PBFs

	pbfFilePath := filepath.Join(pbfStorage, pbfFileName)

	fmt.Println("Importing PBF file: ", pbfFilePath)

	file := openFile(pbfFilePath)
	defer file.Close()

	scanner := createFileScanner(file)
	defer scanner.Close()

	scan(scanner)
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func createFileScanner(file *os.File) *osmpbf.Scanner {
	scanner := osmpbf.New(context.Background(), file, 3)
	return scanner
}

func scan(scanner *osmpbf.Scanner) {
	for scanner.Scan() {
		object := scanner.Object()
		fmt.Println(object)
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}
}
