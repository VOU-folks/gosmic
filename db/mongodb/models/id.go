package models

type ItemType string

const (
	NodeType     ItemType = "node"
	WayType      ItemType = "way"
	RelationType ItemType = "relation"
)

type ID struct {
	ID      int64    `json:"id" bson:"id"`
	Type    ItemType `json:"type" bson:"type"`
	Version int      `json:"version" bson:"version"`
}
