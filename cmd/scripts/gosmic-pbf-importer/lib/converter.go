package lib

import (
	"github.com/paulmach/osm"

	. "gosmic/db/mongodb/models"
)

func ConvertWay(osmWay *osm.Way) Object {
	nodes := make([]int64, 0, len(osmWay.Nodes))
	for _, node := range osmWay.Nodes {
		nodes = append(nodes, int64(node.ID))
	}

	return Object{
		ID: ID{
			ID:      int64(osmWay.ID),
			Type:    WayType,
			Version: osmWay.Version,
		},
		Tags:      ConvertTags(osmWay.Tags),
		Timestamp: osmWay.Timestamp,
		Nodes:     nodes,
	}
}

func ConvertNode(osmNode *osm.Node) Object {
	return Object{
		ID: ID{
			ID:      int64(osmNode.ID),
			Type:    NodeType,
			Version: osmNode.Version,
		},
		Tags:      ConvertTags(osmNode.Tags),
		Timestamp: osmNode.Timestamp,
		Location: Coords{
			Type: "Point",
			Coordinates: []float64{
				osmNode.Lon,
				osmNode.Lat,
			},
		},
	}
}

func ConvertRelation(osmRelation *osm.Relation) Object {
	members := make([]Member, 0, len(osmRelation.Members))
	for _, member := range osmRelation.Members {
		var location *Coords
		if member.Lat != 0.0 && member.Lon != 0.0 {
			location = &Coords{
				Type: "Point",
				Coordinates: []float64{
					member.Lon,
					member.Lat,
				},
			}
		}
		members = append(
			members,
			Member{
				Type:        member.Type,
				Orientation: member.Orientation,
				Ref:         member.Ref,
				Role:        member.Role,
				Location:    location,
			},
		)
	}

	return Object{
		ID: ID{
			ID:      int64(osmRelation.ID),
			Type:    RelationType,
			Version: osmRelation.Version,
		},
		Tags:      ConvertTags(osmRelation.Tags),
		Timestamp: osmRelation.Timestamp,
		Members:   members,
	}
}

func ConvertTags(tags osm.Tags) []Tag {
	result := make([]Tag, 0, len(tags))

	for _, tag := range tags {
		result = append(
			result,
			Tag{
				Key:   tag.Key,
				Value: tag.Value,
			},
		)
	}

	return result
}
