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

// GetChildNodeByAttribute traverses the children Nodes for an Attribute and
// returns the Node that the Attribute is found on.
func (n *Node) GetChildNodeByAttribute(attribute string) (*Node, error) {
	for _, child := range n.Children {
		if child.Attribute == attribute {
			return child, nil
		}
	}

	return nil, errors.New("Could not find attribute on node")
}

// GetDataByAttributes traverses the children Nodes over a list of multiple (or
// one) attribute(s). If there is no structure of child Nodes that has the
// structure requested, an empty string and an error will be returned.
//
// The following examples assume that the individual record with the depth of 0
// is the current Node.
//
// A simple example is when we need to get the name of an individual:
//
//   0 @I1@ INDI
//   1 NAME John
//
// Example:
//
//   name, err := GetDataByAttributes("NAME")
//   if err != nil {
//     // The Node "NAME" was not found
//   }
//
// Or a more complex example would be when you need to get the birth date of an
// individual:
//
//   0 @I1@ INDI
//   1 BIRT
//   2 DATE 15 Jun 1990
//
// Example:
//
//   birthday, err := GetDataByAttributes("BIRT", "DATE")
//   if err != nil {
//     // Either the Node "BIRT" or "DATE" was not found
//   }
func (n *Node) GetDataByAttributes(attributes ...string) (string, error) {
	var err error
	for _, attribute := range attributes {
		n, err = n.GetChildNodeByAttribute(attribute)
		if err != nil {
			return "", err
		}
	}

	return n.Data, nil
}
