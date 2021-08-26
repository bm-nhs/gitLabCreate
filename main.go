package main

import (
	"encoding/json"
	"fmt"
	"goGitBack/models"
)

// Student declares `Student` structure
type Student struct {
	FirstName, lastName string
	Email string
	Age int
	HeightInMeters float64
	IsMale bool
}

func main() {

	// define `john` struct
	john := Student{
		FirstName: "John",
		lastName: "Doe",
		Age: 21,
		HeightInMeters: 1.75,
		IsMale: true,
	}

	// encode `john` as JSON
	johnJSON, _ := json.Marshal( john )

	// print JSON string
	fmt.Println( string(johnJSON) )
}