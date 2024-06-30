package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"

	mongodb "gosmic/db/drivers/mongodb"

	"gosmic/db"
	"gosmic/db/indexes"
	"gosmic/db/models"

	. "gosmic/internal/config"
	"gosmic/internal/structs"

	. "gosmic/cmd/scripts/gosmic-pbf-importer/lib"
)

var config struct {
	Database structs.DatabaseConfig
	Storage  structs.StorageConfig
	Osm      structs.OSMConfig
}

var dbClient *mongodb.Client
var dbInstance *mongodb.Database
var dbCollections *db.Collections

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

	ctx := context.Background()

	// --------------------------------------------
	// Step 1: Connection to Database
	// --------------------------------------------
	dbClient, err = mongodb.Connect(ctx, config.Database.ConnectionString)
	if err != nil {
		panic(err)
	}

	err = mongodb.Ping(ctx, dbClient)
	if err != nil {
		panic(err)
	}

	defer func(ctx context.Context, client *mongodb.Client) {
		err := mongodb.Disconnect(ctx, client)
		if err != nil {
			fmt.Println("Error disconnecting from Database: ", err)
		}
	}(ctx, dbClient)

	dbInstance = mongodb.SwitchToDB(ctx, dbClient, config.Database.DatabaseName)

	fmt.Println("Successfully connected to Database: ", config.Database.DatabaseName)
	// --------------------------------------------

	// --------------------------------------------
	// Step 2: Ensuring db indexes
	// --------------------------------------------
	err = indexes.CreateIndexes(ctx, dbInstance)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created (ensured) db indexes")
	// --------------------------------------------

	// --------------------------------------------
	// Step 3: Initiating collections
	// --------------------------------------------
	dbCollections = &db.Collections{
		Objects: mongodb.GetCollection(ctx, dbInstance, "objects"),
	}
	// --------------------------------------------

	// --------------------------------------------
	// Step 4: Initiating PBF file scanner
	// --------------------------------------------
	pbfFileName := config.Osm.Sources.PBF.FileName
	pbfStorage := config.Storage.PBFs

	pbfFilePath := filepath.Join(pbfStorage, pbfFileName)

	fmt.Println("Importing PBF file: ", pbfFilePath)

	file := OpenFile(pbfFilePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file: ", err)
		}
	}(file)

	scanner := CreateFileScanner(ctx, file)
	defer func(scanner *osmpbf.Scanner) {
		err := scanner.Close()
		if err != nil {
			fmt.Println("Error closing scanner: ", err)
		}
	}(scanner)
	// --------------------------------------------

	// --------------------------------------------
	// Step 3: Launching the import process
	// --------------------------------------------
	run(ctx, scanner, dbInstance, dbCollections)
}

func run(ctx context.Context, scanner *osmpbf.Scanner, dbInstance *mongodb.Database, dbCollections *db.Collections) {
	for scanner.Scan() {
		object := scanner.Object()

		processErr := processObject(ctx, object, dbInstance, dbCollections)
		if processErr != nil {
			panic(processErr)
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}
}

func processObject(ctx context.Context, osmObject osm.Object, dbInstance *mongodb.Database, dbCollections *db.Collections) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var dbObject models.Object

	switch osmObject := osmObject.(type) {
	case *osm.Way:
		dbObject = ConvertWay(osmObject)
	case *osm.Node:
		dbObject = ConvertNode(osmObject)
	case *osm.Relation:
		dbObject = ConvertRelation(osmObject)
	default:
		return fmt.Errorf("unknown object type: %T", osmObject)
	}

	fmt.Println("ID:", dbObject.ID.ID, "Type:", dbObject.ID.Type, "Version:", dbObject.ID.Version)
	fmt.Println("Tags:", dbObject.Tags)
	fmt.Println("Timestamp:", dbObject.Timestamp)
	fmt.Println("Nodes:", dbObject.Nodes)
	fmt.Println("Location:", dbObject.Location.Type, dbObject.Location.Coordinates)
	fmt.Println("Members:", dbObject.Members)

	_, err := saveObject(ctx, dbObject, dbCollections.Objects)
	if err != nil {
		return err
	}

	return nil
}

func saveObject(ctx context.Context, object models.Object, objectsCollection *mongodb.Collection) (*models.Object, error) {
	fmt.Println("// ---- OBJECT ---- //")
	fmt.Printf("%+v\n", object)
	fmt.Println("// ---- ------ ---- //")

	_, err := mongodb.InsertOne(ctx, object, objectsCollection)
	if err != nil {
		panic(err)
	}

	//fmt.Println("saved")
	//fmt.Printf("%+v\n", result)
	//fmt.Println("// ---- ------ ---- //")

	return nil, nil
}
