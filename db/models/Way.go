package models

import "time"

type Way struct {
	ID        uint64    `json:"_id" bson:"_id"`
	Version   int       `json:"version" bson:"version"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Tags      []Tag     `json:"tags" bson:"tags"`
	Location  Coords    `json:"location,omitempty" bson:"location,omitempty"`
	Nodes     []int64   `json:"nodes,omitempty" bson:"nodes,omitempty"`
}

func (w Way) GetType() string {
	return "way"
}
