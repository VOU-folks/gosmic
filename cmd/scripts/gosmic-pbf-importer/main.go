package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

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

var objectsChannel chan osm.Object

var batchChannel chan []interface{}
var batchSize uint = 10000
var mxBatch sync.Mutex
var batch []interface{}

var processedObjectsCount atomic.Uint64

var timeCheckpoint time.Time
var startedAt time.Time

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

	objectsChannel = make(chan osm.Object, runtime.NumCPU()*10)
	batchChannel = make(chan []interface{})
	batch = make([]interface{}, batchSize)
	mxBatch = sync.Mutex{}

	timeCheckpoint = time.Now()
	startedAt = time.Now()
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
	// Step 5: Starting object processor
	// --------------------------------------------
	go objectProcessor(ctx, objectsChannel, dbInstance, dbCollections)
	go batchProcessor(ctx, batchChannel, dbInstance, dbCollections)

	// --------------------------------------------
	// Step 6: Launching the import process
	// --------------------------------------------
	run(ctx, scanner, dbInstance, dbCollections)
}

func run(ctx context.Context, scanner *osmpbf.Scanner, dbInstance *mongodb.Database, dbCollections *db.Collections) {
	for scanner.Scan() {
		objectsChannel <- scanner.Object()
	}

	if len(batch) > 0 {
		batchChannel <- batch
	}

	fmt.Println("Total processed objects count: ", processedObjectsCount.Load())
	fmt.Println("Total time elapsed: ", time.Since(startedAt))

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}
}

func objectProcessor(ctx context.Context, objectsChannel chan osm.Object, dbInstance *mongodb.Database, dbCollections *db.Collections) {
	for {
		select {
		case object := <-objectsChannel:
			err := processObject(ctx, object, dbInstance, dbCollections)
			if err != nil {
				fmt.Errorf("Error processing object: %v", err)
			}
		case <-ctx.Done():
			return
		}
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

	mxBatch.Lock()
	counter := processedObjectsCount.Load()
	increment := counter % uint64(batchSize)
	batch[increment] = dbObject
	if increment == uint64(batchSize)-1 {
		batchChannel <- batch
		batch = make([]interface{}, batchSize)
	}
	mxBatch.Unlock()

	processedObjectsCount.Add(1)
	if processedObjectsCount.Load()%1000000 == 0 {
		fmt.Println("\n----\\Checkpoint\\----")
		fmt.Println("Processed objects count: ", processedObjectsCount.Load())
		fmt.Printf("Last object: %v\n", dbObject)
		fmt.Println("Time elapsed: ", time.Since(timeCheckpoint))
		fmt.Println("----/Checkpoint/----")

		timeCheckpoint = time.Now()
	}

	return nil
}

func batchProcessor(ctx context.Context, batchChannel chan []interface{}, dbInstance *mongodb.Database, dbCollections *db.Collections) {
	for {
		select {
		case batchObjects := <-batchChannel:
			_, err := mongodb.InsertMany(ctx, batchObjects, dbCollections.Objects)
			if err != nil {
				fmt.Errorf("Error inserting batch: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
