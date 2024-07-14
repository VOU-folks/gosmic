package lib

import (
	"github.com/paulmach/osm"

	. "gosmic/db/models"
)

func ConvertWay(osmWay *osm.Way) Way {
	nodes := make([]int64, 0, len(osmWay.Nodes))
	for _, node := range osmWay.Nodes {
		nodes = append(nodes, int64(node.ID))
	}

	return Way{
		Type: "way",

		ID:        uint64(osmWay.ID),
		Version:   osmWay.Version,
		Tags:      ConvertTags(osmWay.Tags),
		Timestamp: osmWay.Timestamp,
		Nodes:     nodes,
	}
}

func ConvertNode(osmNode *osm.Node) Node {
	return Node{
		Type: "node",

		ID:        uint64(osmNode.ID),
		Version:   osmNode.Version,
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
