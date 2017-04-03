package main

import (
	"github.com/lukevers/gedcom"
	"log"
)

func main() {
	tree, err := gedcom.ParseFromFile("simple.ged")
	if err != nil {
		log.Fatalln("Error loading file:", err)
	}

	err = tree.TraverseFamilies()
	if err != nil {
		log.Fatalln("Error traversing families:", err)
	}

	for i, family := range tree.Families {
		log.Printf("===== FAMILY %d =====", i)
		if family.Father != nil {
			name, _ := family.Father.GetAttribute("NAME")
			log.Printf("Father: %s", name)
		}
		if family.Mother != nil {
			name, _ := family.Mother.GetAttribute("NAME")
			log.Printf("Mother: %s", name)
		}

		for j, child := range family.Children {
			name, _ := child.GetAttribute("NAME")
			log.Printf("Child %d: %s", j, name)
		}
	}
}
