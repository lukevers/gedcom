package gedcom

// Family represents a family record in a node structure. It can be used to
// traverse a Tree with a more friendly API.
type Family struct {
	Father   *Individual
	Mother   *Individual
	Children []*Individual
}
