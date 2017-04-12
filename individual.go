package gedcom

// Individual represents an individual record in a node structure with links to
// other individuals.
type Individual struct {
	Node     *Node
	Father   *Individual
	Mother   *Individual
	Children []*Individual
}

// GetName returns the name of the individual. If there is no name attached to
// the individual, an empty string is returned instead.
func (i *Individual) GetName() string {
	name, err := i.Node.GetDataByAttributes("NAME")
	if err != nil {
		return ""
	}

	return name
}

// GetBirthday returns the birth date of the individual. If there is no birth
// date attached to the individual, an empty string is returned instead. A
// string is used here, as the GEDCOM format does not specify one specific
// format for date objects.
func (i *Individual) GetBirthday() string {
	birthday, err := i.Node.GetDataByAttributes("BIRT", "DATE")
	if err != nil {
		return ""
	}

	return birthday
}
