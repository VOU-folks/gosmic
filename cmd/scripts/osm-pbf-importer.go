package main

import (
	"flag"
	"fmt"
	"path/filepath"

	. "osm-api/internal/config"
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
}
