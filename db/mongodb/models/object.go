package models

import "time"

type Object struct {
	ID        ID        `json:"_id" bson:"_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Tags      []Tag     `json:"tags" bson:"tags"`
	Location  Coords    `json:"location,omitempty" bson:"location,omitempty"`
	Nodes     []int64   `json:"nodes,omitempty" bson:"nodes,omitempty"`
	Members   []Member  `json:"members,omitempty" bson:"members,omitempty"`
}
