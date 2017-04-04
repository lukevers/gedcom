package gedcom

// Individual represents an individual record in a node structure with links to
// other individuals.
type Individual struct {
	Node   *Node
	Father *Individual
	Mother *Individual
}

// GetName returns the name of the individual. If there is no name attached to
// the individual, an empty string is returned instead.
func (i *Individual) GetName() string {
	name, err := i.Node.GetAttribute("NAME")
	if err != nil {
		return ""
	}

	return name
}
