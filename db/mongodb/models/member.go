package models

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/osm"
)

type Member struct {
	Type osm.Type `json:"type" bson:"type"`
	Ref  int64    `json:"ref" bson:"ref"`
	Role string   `json:"role" bson:"role"`

	Location *Coords `json:"location" bson:"location"`

	// Orientation is the direction of the way around a ring of a multipolygon.
	// Only valid for multipolygon or boundary relations.
	Orientation orb.Orientation `json:"orienation" bson:"orienation"`
}
