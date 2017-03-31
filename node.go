package gedcom

import (
	"errors"
)

type Node struct {
	Depth     int
	Attribute string
	Data      string

	Parent   *Node
	Children []*Node
}

func (n *Node) GetAttribute(attribute string) (string, error) {
	for _, child := range n.Children {
		if child.Attribute == attribute {
			return child.Data, nil
		}
	}

	return "", errors.New("Could not find attribute on node")
}
