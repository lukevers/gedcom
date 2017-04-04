package gedcom

// Individual represents an individual record in a node structure with links to
// other individuals.
type Individual struct {
	Node   *Node
	Father *Individual
	Mother *Individual
}
