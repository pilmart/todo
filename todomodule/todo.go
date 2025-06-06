package main

// imports
import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

// struct used to maintain an individual todo item
type ToDo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {

	// Declare variable for our new todo item
	var toDo ToDo

	// ToDos array of ToDo items
	var toDos []ToDo

	var filePath = "./data/todos.json"

	/*
	 Initial flag setup to capture new todo information
	 id - generate random int simulate db id generation
	 description - default to empty string
	 status - not started
	*/

	// Load todos from json file
	toDos = loadAll(filePath)

	toDo.Id = rand.IntN(100000) + 1
	log.Printf("Generated id : %d", toDo.Id)

	flag.StringVar(&toDo.Description, "description", " ", "Description of to do item")
	flag.StringVar(&toDo.Status, "status", "not started", "Status of to do item")
	flag.Parse()

	// actions - Create / Read / Update / Delete

	// print the single item to console
	log.Printf("New to do item %+v\n", toDo)

	// Add the new todo item to the array
	toDos = append(toDos, toDo)

	// Dump the array to the console
	var sb strings.Builder

	for i := 0; i < len(toDos); i++ {

		sb.WriteString("id : " + strconv.Itoa(toDos[i].Id) + " ")
		sb.WriteString("description : " + toDos[i].Description + " ")
		sb.WriteString("status : " + toDos[i].Status + "\n")
	}

	log.Printf("%s : "+sb.String(), filePath)

	// save the revised data here
	saveAll(toDos)

}

func saveAll(todos []ToDo) {
	log.Println("Starting saveAll")
	log.Printf("saveAll : New to do items %+v\n", todos)

	filePath := "./data/todos.json"

	// Marshal the struct into JSON
	jsonData, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
		return
	}
	log.Println("Marshalling completed")

	// open up a new file....
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return
	}
	defer file.Close()
	log.Printf("File Created %s", filePath)

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
		return
	}
	log.Printf("File Written %s", filePath)
	log.Println("saveAll completes")
}

func loadAll(filePath string) []ToDo {

	log.Printf("Starting loadAll from file %s", filePath)

	// is of type ToDos
	var todos []ToDo

	// Check file existence
	jsonFileExists := checkFileExists(filePath)

	if !jsonFileExists {
		log.Printf("Unable to open file : %s an empty array will be returned ", filePath)
		return todos
	}
	log.Printf("File : %s located...", filePath)

	data, err := os.Open(filePath)
	// Handle any error
	if err != nil {
		log.Fatalf("Unable to open file : %v", err)
	}
	log.Println("Successfully Opened todos.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer data.Close()

	// Read the file in....
	byteValue, err := io.ReadAll(data)
	if err != nil {
		log.Fatalf("loadAll - Unable to ReadAll :  %v", err)
	}
	log.Println("Successfully Read bytevalue")

	// unmarshal the json...
	err = json.Unmarshal(byteValue, &todos)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %v", err)
	}
	log.Println("To Do's contains " + strconv.Itoa(len(todos)) + " items ")
	// hand back our result..
	log.Println("loadAll completes...")
	return todos
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
