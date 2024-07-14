package main

import (
	"flag"
	"fmt"
	. "gosmic/cmd/scripts/gosmic-pbf-downloader/lib"

	. "gosmic/internal/config"
	"gosmic/internal/structs"
)

var config struct {
	Database structs.DatabaseConfig
	Storage  structs.StorageConfig
	Osm      structs.OSMConfig
}

func init() {
	var err error

	configFile := flag.String("config", "config.yaml", "path to the config file. Example: api -config /full/path/to/config.yaml")
	flag.Parse()

	configService, err := GetConfigService(*configFile)
	if err != nil {
		panic(err)
	}

	config.Osm = configService.GetOsmConfig()
	config.Storage = configService.GetStorageConfig()
	config.Database = configService.GetDatabaseConfig()
}

func main() {
	var err error

	pbfSources := config.Osm.Sources.PBFs
	pbfStorage := config.Storage.PBFs

	err = EnsureFolderExists(pbfStorage)
	if err != nil {
		panic(err)
	}

	for _, pbfSource := range pbfSources {
		fmt.Println("Downloading PBF file for region:", pbfSource.Region)
		err = DownloadPBF(pbfSource.Url, pbfStorage, pbfSource.FileName)
		if err != nil {
			fmt.Printf("Error downloading PBF file for region %s: %s\n", pbfSource.Region, err.Error())
		}
	}
}
