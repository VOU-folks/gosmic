package models

import "time"

type Node struct {
	Type string

	ID        uint64    `json:"_id" bson:"_id"`
	Version   int       `json:"version" bson:"version"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Tags      []Tag     `json:"tags" bson:"tags"`
	Location  Coords    `json:"location,omitempty" bson:"location,omitempty"`
	Nodes     []int64   `json:"nodes,omitempty" bson:"nodes,omitempty"`
}
