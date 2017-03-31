package gedcom

import (
	"errors"
)

// Node is a GEDCOM line, but also contains other information to help traversal
// and ease of use. Some examples of GEDCOM lines look like this:
//
//   0 @I1@ INDI
//   1 NAME John
//   1 SEX  M
//          ^^^^^^^ Data
//     ^^^^-------- Attribute
//   ^------------- Depth
//
// Depth will always be an integer, while Attribute and Data are basically just
// strings. The delimiter between the three data points is a space (" "), but
// Data contains the rest of the line after the second delimiter.
//
// Depth puts the parent/child nodes into perspective. In our example above,
// an individual (@I1@) node is defined with two children nodes, both one node
// deep. While this individual node will have no parent Node, it has two
// children Nodes, which will both have the individual (@I1@) Node as their
// parent Node. Children nodes can have children nodes too. For example:
//
//   0 HEAD
//   1 SOUR PAF
//   2 NAME Personal Ancestral File
//   1 DATE 31 MAR 2017
//
// In the example above, the Node structure would look like this:
//
//   &Node{ // 0xc420014230
//     Depth:     0,
//     Attribute: "HEAD",
//     Data:      "",
//     Parent:    nil,
//     Children:  []*Node{
//       &Node{ // 0xc4200143c0
//         Depth:     1,
//         Attribute: "SOUR",
//         Data:      "PAF",
//         Parent:    0xc420014230,
//         Children:  []*Node{
//           &Node{ // 0xc420014410
//             Depth:     2,
//             Attribute: "NAME",
//             Data:      "Personal Ancestral File",
//             Parent:    0xc4200143c0,
//             Children:  []*Node{},
//           },
//         },
//       },
//       &Node{ // 0xc420014500
//         Depth:     1,
//         Attribute: "DATE",
//         Data:      "31 MAR 2017",
//         Parent:    0xc420014230,
//         Children:  []*Node{},
//       },
//     },
//   }
//
// As you can see, a Node with a Depth of n will only ever have children with a
// Depth of n+1, and is not limited to any number of children.
type Node struct {
	Depth     int
	Attribute string
	Data      string

	Parent   *Node
	Children []*Node
}

// GetAttribute traverses the children Nodes for an Attribute and returns the
// Data attached to the Attribute. If there is no children Node that contains
// the Attribute searched for, an error will be returned with an empty string
// instead of any Data.
func (n *Node) GetAttribute(attribute string) (string, error) {
	for _, child := range n.Children {
		if child.Attribute == attribute {
			return child.Data, nil
		}
	}

	return "", errors.New("Could not find attribute on node")
}
