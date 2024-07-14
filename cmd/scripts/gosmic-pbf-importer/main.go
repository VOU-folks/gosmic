package main

import (
	"context"
	"flag"
	"fmt"
	"gosmic/db/models"
	"gosmic/internal/lo"
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

var mxBatch sync.Mutex
var batch []interface{}
var batchSize uint = 10000

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

	mxBatch = sync.Mutex{}
	batch = make([]interface{}, batchSize)

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
		Nodes: mongodb.GetCollection(ctx, dbInstance, "nodes"),
		Ways:  mongodb.GetCollection(ctx, dbInstance, "ways"),
	}
	// --------------------------------------------

	// --------------------------------------------
	// Step 4: Starting object processor
	// --------------------------------------------
	go objectProcessor(ctx, objectsChannel)
	go batchProcessor(ctx, batchChannel, dbCollections)

	// --------------------------------------------
	// Step 5: Launching the import process
	// --------------------------------------------
	for _, pbfSource := range config.Osm.Sources.PBFs {
		pbfFileName := pbfSource.FileName
		pbfStorage := config.Storage.PBFs
		pbfFilePath := filepath.Join(pbfStorage, pbfFileName)

		fmt.Println("Importing PBF file: ", pbfFilePath)
		scanFile(ctx, pbfFilePath)
	}
}

func scanFile(ctx context.Context, pathToPBF string) {
	file := OpenFile(pathToPBF)
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

func objectProcessor(ctx context.Context, objectsChannel chan osm.Object) {
	for {
		select {
		case object := <-objectsChannel:
			err := processObject(ctx, object)
			if err != nil {
				fmt.Errorf("Error processing object: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func processObject(ctx context.Context, osmObject osm.Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var dbObject interface{}

	switch object := osmObject.(type) {
	case *osm.Way:
		dbObject = ConvertWay(object)
	case *osm.Node:
		dbObject = ConvertNode(object)
	default:
		return fmt.Errorf("unknown object type: %T", osmObject)
	}

	putDbObjectToBatch(dbObject)
	processObjectCheckPoint(dbObject)

	return nil
}

func putDbObjectToBatch(dbObject interface{}) {
	counter := processedObjectsCount.Load()
	increment := counter % uint64(batchSize)

	mxBatch.Lock()
	batch[increment] = dbObject
	if increment == uint64(batchSize)-1 {
		batchChannel <- batch
		batch = make([]interface{}, batchSize)
	}
	mxBatch.Unlock()
}

func processObjectCheckPoint(dbObject interface{}) {
	processedObjectsCount.Add(1)
	if processedObjectsCount.Load()%1000000 == 0 {
		fmt.Println("\n----\\Checkpoint\\----")
		fmt.Println("Processed objects count: ", processedObjectsCount.Load())
		fmt.Printf("Last object: %v\n", dbObject)
		fmt.Println("Time elapsed: ", time.Since(timeCheckpoint))
		fmt.Println("----/Checkpoint/----")

		timeCheckpoint = time.Now()
	}
}

func batchProcessor(ctx context.Context, batchChannel chan []interface{}, dbCollections *db.Collections) {
	var err error

	for {
		select {
		case batchObjects := <-batchChannel:
			// Saving Nodes to DB
			nodes := lo.Filter(batchObjects, func(obj interface{}) bool {
				return obj.(models.Node).Type == "node"
			})

			if len(nodes) > 0 {
				_, err = saveNodesToDb(ctx, nodes, dbCollections)
				if err != nil {
					fmt.Println("Error saving nodes to db: ", err)
				}
			}

			// Saving Ways to DB
			ways := lo.Filter(batchObjects, func(obj interface{}) bool {
				return obj.(models.Way).Type == "way"
			})

			if len(ways) > 0 {
				_, err = saveWaysToDb(ctx, ways, dbCollections)
				if err != nil {
					fmt.Println("Error saving ways to db: ", err)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func saveNodesToDb(ctx context.Context, nodes []interface{}, dbCollections *db.Collections) (*mongodb.InsertManyResult, error) {
	return mongodb.InsertMany(ctx, nodes, dbCollections.Nodes)
}

func saveWaysToDb(ctx context.Context, ways []interface{}, dbCollections *db.Collections) (*mongodb.InsertManyResult, error) {
	return mongodb.InsertMany(ctx, ways, dbCollections.Ways)
}
