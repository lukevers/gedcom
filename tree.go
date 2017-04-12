package gedcom

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// Tree contains a node structure of a GEDCOM file.
type Tree struct {
	Nodes       []*Node
	Families    []*Family
	Individuals []*Individual
}

// ParseFromFile loads a file into memory and parses it to a Tree.
func ParseFromFile(file string) (*Tree, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	str := string(bytes)
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.TrimSpace(str)

	return Parse(strings.Split(str, "\n"))
}

// Parse takes a slice of GEDCOM lines and converts it to a Node structure in a
// Tree.
func Parse(lines []string) (*Tree, error) {
	t := &Tree{}
	var nodes []*Node

	// Convert every line to a node
	for _, line := range lines {
		parts := strings.Split(line, " ")
		n := &Node{}
		var err error
		n.Depth, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		n.Attribute = strings.TrimSpace(parts[1])
		n.Data = strings.TrimSpace(strings.Join(parts[2:], " "))

		nodes = append(nodes, n)
	}

	// Temporary root node that is changed throughout loop
	var root *Node

	// Loop through every node and assign parent and children nodes
	for index, node := range nodes {
		// If index is 0 we have a new root element
		if node.Depth == 0 {
			t.Nodes = append(t.Nodes, node)
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

	return t, nil
}

// TraverseFamilies loops over all nodes and creates a slice of Families and a
// slice of Individuals.
func (t *Tree) TraverseFamilies() error {
	individuals := make(map[string]*Individual)
	fams := make(map[string]*Node)
	families := make(map[string]*Family)

	for _, node := range t.Nodes {
		switch node.Data {
		case "INDI":
			individual := &Individual{Node: node}
			individuals[node.Attribute] = individual
			t.Individuals = append(t.Individuals, individual)
		case "FAM":
			fams[node.Attribute] = node
		}
	}

	for id, family := range fams {
		f := &Family{}

		for _, node := range family.Children {
			individual := individuals[node.Data]

			switch node.Attribute {
			case "HUSB":
				f.Father = individual
			case "WIFE":
				f.Mother = individual
			case "CHIL":
				f.Children = append(f.Children, individual)
			}
		}

		families[id] = f
		t.Families = append(t.Families, f)
	}

	for _, individual := range individuals {
		fam, err := individual.Node.GetDataByAttributes("FAMC")

		// If there is an error that means that the individual is not a child
		// of a family and has no parents set and we should move on to the next
		// individual.
		if err != nil {
			continue
		}

		if f, exists := families[fam]; exists {
			if f.Father != nil {
				individual.Father = f.Father
				f.Father.Children = append(f.Father.Children, individual)
			}

			if f.Mother != nil {
				individual.Mother = f.Mother
				f.Mother.Children = append(f.Mother.Children, individual)
			}
		}
	}

	return nil
}

// FindIndividualByAttribute searches through the slice of Individuals on the
// Tree and looks for an Individual that's internal Node structure has the
// attribute provided with the same data provided.
func (t *Tree) FindIndividualByAttribute(attribute, data string) *Individual {
	for _, individual := range t.Individuals {
		if a, err := individual.Node.GetDataByAttributes(attribute); err == nil {
			if a == data {
				return individual
			}
		}
	}

	return nil
}
