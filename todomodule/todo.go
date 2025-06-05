package main

// imports
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	todoFile, err := os.Open("./data/todos.json")
	// if we returns an error then handle it
	if err != nil {
		log.Fatalf("Unable to open file : %s", err)
	}
	fmt.Println("Successfully Opened users.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer todoFile.Close()

	byteValue, err := io.ReadAll(todoFile)
	if err != nil {
		log.Fatalf("Unable to ReadAll :  %s", err)
	}
	fmt.Println("Successfully Read bytevalue")
	fmt.Printf("%s", byteValue)

	var todos ToDos

	err = json.Unmarshal(byteValue, &todos)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}
	fmt.Println("length " + strconv.Itoa(len(todos.ToDos)))

	for i := 0; i < len(todos.ToDos); i++ {
		fmt.Println("id : " + strconv.Itoa(todos.ToDos[i].Id))
		fmt.Println("description : " + todos.ToDos[i].Description)
		fmt.Println("status : " + todos.ToDos[i].Status)
	}
}

// struct used to maintain an individual todo item
type ToDo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// ToDo struct which contains
// an array of users
type ToDos struct {
	ToDos []ToDo `json:"todos"`
}
