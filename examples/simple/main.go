package main

import (
	"github.com/4thel00z/fakenews/pkg/v1/fakenews"
	"log"
)

type Person struct {
	Age  int
	Name string
}

func main() {
	// Create a rule to set a random integer between 1 and 10 for the Age field
	ageRule := fakenews.RandomInt[Person]("Age", 1, 10)
	nameRule := fakenews.RandomPick[Person, string]("Name", "John", "Jane", "Jack")
	fake, err := fakenews.Fake[Person](Person{}, ageRule, nameRule)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("My age is", fake.Age)
	log.Println("My name is", fake.Name)
}
