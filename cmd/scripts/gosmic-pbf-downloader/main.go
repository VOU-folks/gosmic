package main

import (
	"flag"
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

	pbfSource := config.Osm.Sources.PBF.Url
	pbfFileName := config.Osm.Sources.PBF.FileName
	pbfStorage := config.Storage.PBFs

	err = EnsureFolderExists(pbfStorage)
	if err != nil {
		panic(err)
	}

	err = DownloadPBF(pbfSource, pbfStorage, pbfFileName)
	if err != nil {
		panic(err)
	}
}
