package main

import (
	"fmt"
	"sync"
	"time"
)

type Person struct {
	name string
	age  uint64
}

func PrintPerson(person Person, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(person.age) * 10 * time.Millisecond)
	fmt.Printf("Person name: %s, Person age: %d\n", person.name, person.age)
}

// Simulation of Concurrent Transactions:
// Use threads to simulate transactions that attempt to access two resources shared: data items X and Y.
func Ex1(persons []Person) {
	var wg sync.WaitGroup

	for _, person := range persons {
		wg.Add(1)
		go PrintPerson(person, &wg)
	}

	wg.Wait()
}

// Access Control to Shared Resources:
// Use a mutex-like structure (binary lock) to control exclusive access to data items X and Y. Threads must attempt to obtain the lock to access these resources concurrently.
func Ex2(persons []Person) {
	var wg sync.WaitGroup

	for _, person := range persons {
		wg.Add(1)
		go PrintPerson(person, &wg)
	}

	wg.Wait()
}

func main() {
	persons := []Person{
		{name: "Eduardo Savian", age: 100},
		{name: "Henrique Zimmermann", age: 1},
		{name: "Steff Tousant", age: 50},
	}

	Ex1(persons)
	Ex2(persons)
}
