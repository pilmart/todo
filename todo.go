package main

// imports
import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"todo/datatypes"
)

func main() {

	// Flag values
	var Id int
	var action string
	var description string
	var status string

	/*
	 Initial flag setup to capture new todo information
	 id - only relevant for update/delete options
	 action - one of show/update/create/delete
	 description - default to empty string - relevant for create
	 status - default to not started - relevant for create
	*/

	flag.StringVar(&action, "action", "show", "Selected action")
	flag.StringVar(&description, "description", " ", "Description of to do item")
	flag.StringVar(&status, "status", "not started", "Status of to do item")
	flag.IntVar(&Id, "Id", 0, "Mandatory for both update/delete actions")

	flag.Parse()

	// actions - Create / Show / Update / Delete

	fmt.Printf("Selected action..%s\n", strings.ToLower(action))
	switch strings.ToLower(action) {
	case "show":
		// show all records no params needed
		showAllRecords()
	case "create":
		// create new todo item
		create(description, status)
	case "update":
		var updatedToDo datatypes.ToDo
		updatedToDo.Id = Id
		updatedToDo.Description = description
		updatedToDo.Status = status
		update(updatedToDo)
	case "delete":
		// delete function
		delete(Id)
	default:
		log.Printf("Unsupported action..%s passed to procedure\n", strings.ToLower(action))
	}
}

// Show all the current items on the console.
func showAllRecords() {
	// ToDos array of ToDo items
	var toDos []datatypes.ToDo
	filePath := "./data/todos.json"
	log.Println("Starting showAllRecords..")

	// Load todos from json file
	toDos = loadAll(filePath)

	// Dump the array to the console
	var sb strings.Builder

	for i := 0; i < len(toDos); i++ {

		sb.WriteString("id : " + strconv.Itoa(toDos[i].Id) + " ")
		sb.WriteString("description : " + toDos[i].Description + " ")
		sb.WriteString("status : " + toDos[i].Status + "\n")
	}

	log.Printf("%s : "+sb.String(), filePath)
	log.Println("showAllRecords completes")

}

// Create a single ToDo item and persist back to file
// Note :-  we have default values set in the flags so we can just create with those
func create(description string, status string) {
	filePath := "./data/todos.json"
	var toDos []datatypes.ToDo
	log.Printf("Starting create..with desc :%s and status: %s", description, status)

	// Load todos from json file
	toDos = loadAll(filePath)

	// Add the new todo item to the array
	var toDo datatypes.ToDo
	toDo.Id = getNextId(toDos)
	toDo.Status = status
	toDo.Description = description

	// add new item to array
	toDos = append(toDos, toDo)

	// save the revised data here
	saveAll(toDos)

	log.Printf("create completes after saving record %v\n", toDo)

}

// Update a single ToDo item and persist back to file, make sure
// incoming ToDo item has a sensible Id otherwise bail out
func update(toDo datatypes.ToDo) {
	// check for uninitialised Id
	if toDo.Id <= 0 {
		log.Println("ToDo Id uninitialised - no action taken")
		return
	}

	filePath := "./data/todos.json"
	var toDos []datatypes.ToDo
	log.Printf("Starting update..for record : %v", toDo)

	// same as delete here - probably refactor out later
	// Load todos from json file
	toDos = loadAll(filePath)

	// scan the array for the required id and capture its index
	var currIndx int = -1
	for i := 0; i < len(toDos); i++ {
		if toDos[i].Id == toDo.Id {
			currIndx = i
			break
		}
	}

	if currIndx > -1 {
		log.Printf("Id..%d is located at position %d in Array, and can be updated with the new item ", toDo.Id, currIndx)
		// In place array update using currindx
		toDos[currIndx].Description = toDo.Description
		toDos[currIndx].Status = toDo.Status
		log.Printf("Updated Array..%v", toDos)
		// persist back to file
		saveAll(toDos)
		log.Printf("Update for record..%v complete", toDo)
	} else {
		log.Printf("Id..%d cannot be located - no action taken", toDo.Id)
	}
	log.Println("delete completes")

}

// delete a record by id with check to ensure id is sensible
func delete(Id int) {

	// check for uninitialised Id
	if Id <= 0 {
		log.Println("Id uninitialised - no action")
		return
	}

	var toDos []datatypes.ToDo
	filePath := "./data/todos.json"
	log.Printf("Starting delete for id..%d", Id)

	// Load todos from json file
	toDos = loadAll(filePath)

	// scan the array for the required id and capture its index
	var currIndx int = -1
	for i := 0; i < len(toDos); i++ {
		if toDos[i].Id == Id {
			currIndx = i
			break
		}
	}

	if currIndx > -1 {
		log.Printf("Id..%d is located at position %d in Array, and would be deleted !!!", Id, currIndx)
		var newToDos = append(toDos[:currIndx], toDos[currIndx+1:]...)
		log.Printf("New Array..%v", newToDos)
		// persist back to file
		saveAll(newToDos)
		log.Printf("Delete for id..%d complete", Id)
	} else {
		log.Printf("Id..%d cannot be located - no action taken", Id)
	}
	log.Println("delete completes")
}

// saves all items in ToDo array back to the specified json file
func saveAll(todos []datatypes.ToDo) {
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

// Loads all items in from the specified file and returns an array of ToDo Items
func loadAll(filePath string) []datatypes.ToDo {

	log.Printf("Starting loadAll from file %s", filePath)

	// is of type ToDos
	var todos []datatypes.ToDo

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

/*
	Helper functions below this comment line
*/
// Scans the array and returns highest id int + 1
func getNextId(toDos []datatypes.ToDo) int {
	if len(toDos) == 0 {
		// no elements so set initial value...
		return 1
	}

	var highCount int
	for i := 0; i < len(toDos); i++ {
		log.Printf("Comparing %d with %d", toDos[i].Id, highCount)
		if toDos[i].Id > highCount {
			highCount = toDos[i].Id
		}
	}
	log.Printf("highCount  : %d", highCount)
	return highCount + 1
}

// Make sure the file exists, probably better ways to do this
func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
