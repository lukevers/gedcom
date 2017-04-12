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
			log.Printf("Father: %s", family.Father.GetName())

			for j, child := range family.Father.Children {
				log.Printf("Child %d of Father: %s", j, child.GetName())
			}
		}
		if family.Mother != nil {
			log.Printf("Mother: %s", family.Mother.GetName())

			for j, child := range family.Mother.Children {
				log.Printf("Child %d of Mother: %s", j, child.GetName())
			}
		}

		for j, child := range family.Children {
			log.Printf("Child %d: %s", j, child.GetName())

			if child.Father != nil {
				log.Printf("Child %d Father: %s", j, child.Father.GetName())
			}

			if child.Mother != nil {
				log.Printf("Child %d Mother: %s", j, child.Mother.GetName())
			}
		}
	}
}
