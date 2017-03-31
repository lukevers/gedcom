package gedcom

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type File struct {
	Nodes []*Node
}

func FromFile(file string) (*File, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	str := string(bytes)
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.TrimSpace(str)

	f := &File{}
	err = f.parse(strings.Split(str, "\n"))

	return f, err
}

func (f *File) parse(lines []string) error {
	var nodes []*Node

	// Convert every line to a node
	for _, line := range lines {
		parts := strings.Split(line, " ")
		n := &Node{}
		var err error
		n.Depth, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		n.Attribute = parts[1]
		n.Data = strings.Join(parts[2:], " ")

		nodes = append(nodes, n)
	}

	// Temporary root node that is changed throughout loop
	var root *Node

	// Loop through every node and assign parent and children nodes
	for index, node := range nodes {
		// If index is 0 we have a new root element
		if node.Depth == 0 {
			f.Nodes = append(f.Nodes, node)
			root = node
			continue
		}

		// If depth is 1, the root element is the parent of this node
		if node.Depth == 1 {
			node.Parent = root
			root.Children = append(root.Children, node)
			continue
		}

		// If depth is > 1, the parent element of this node is a node that we
		// have already processed.
		if node.Depth > 1 {
			for i := index - 1; i > 0; i-- {
				if nodes[i].Depth == node.Depth {
					continue
				}

				if nodes[i].Depth == node.Depth-1 {
					node.Parent = nodes[i]
					nodes[i].Children = append(node.Parent.Children, node)
					break
				}
			}

			continue
		}
	}

	return nil
}