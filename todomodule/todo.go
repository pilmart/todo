package main

// imports
import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

// struct used to maintain an individual todo item

func main() {

	// Declare variable for our new todo item
	var toDo ToDo

	/*
	 Initial flag setup to capture new todo information
	 id - generate random int simulate db id generation
	 description - default to empty string
	 status - not started
	*/

	toDo.Id = rand.IntN(100000) + 1
	log.Printf("Generated id : %d", toDo.Id)

	flag.StringVar(&toDo.Description, "description", " ", "Description of to do item")
	flag.StringVar(&toDo.Status, "status", "not started", "Status of to do item")
	flag.Parse()

	log.Printf("New to do item %+v\n", toDo)

	// todos from json file
	var todos = loadAll()

	// bolt in our new todo
	todos.ToDos = append(todos.ToDos, toDo)

	var sb strings.Builder

	for i := 0; i < len(todos.ToDos); i++ {
		sb.WriteString("id : " + strconv.Itoa(todos.ToDos[i].Id) + " ")
		sb.WriteString("description : " + todos.ToDos[i].Description + " ")
		sb.WriteString("status : " + todos.ToDos[i].Status + "\n")
	}

	todos.ToDos = append(todos.ToDos, toDo)

	fmt.Println("Data/todos.json " + sb.String())

}

func loadAll() ToDos {

	log.Println("Starting loadAll")

	// is of type ToDos
	var todos ToDos

	todoFile, err := os.Open("./data/todos.json")
	// if we returns an error then handle it
	if err != nil {
		log.Fatalf("Unable to open file : %s", err)
	}
	fmt.Println("Successfully Opened todos.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer todoFile.Close()

	// Read the file in....
	byteValue, err := io.ReadAll(todoFile)
	if err != nil {
		log.Fatalf("loadAll - Unable to ReadAll :  %s", err)
	}
	log.Println("Successfully Read bytevalue")

	// marshal the json...
	err = json.Unmarshal(byteValue, &todos)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}
	log.Println("To Do's contains " + strconv.Itoa(len(todos.ToDos)) + " items ")
	// hand back our result..
	return todos
}

func saveAll() {

}

// ToDo Struct outlining a simple to do item
type ToDo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// ToDos struct which contains array of ToDo items
type ToDos struct {
	ToDos []ToDo `json:"todos"`
}
