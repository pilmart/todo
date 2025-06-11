package dataaccess

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"todo/model"
	"todo/utils"
)

// Show all the current items on the console.
// make functions public by capitalising function names
func ShowAllRecords(ctx context.Context) {
	// ToDos array of ToDo items
	var toDos []model.ToDo
	filePath := "./data/todos.json"
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "not found"
	}
	slog.Info("Starting ShowAllRecords for traceID", "traceID", traceID)

	// Load todos from json file
	toDos = loadAll(filePath)

	// Dump the array to the console
	var sb strings.Builder

	for i := 0; i < len(toDos); i++ {

		sb.WriteString("id : " + strconv.Itoa(toDos[i].Id) + " ")
		sb.WriteString("description : " + toDos[i].Description + " ")
		sb.WriteString("status : " + toDos[i].Status + "\n")
	}

	slog.Info(sb.String())
	slog.Info("ShowAllRecords completes")

}

// Create a single ToDo item and persist back to file
// Note :-  we have default values set in the flags so we can just create with those
func Create(ctx context.Context, description string, status string) {
	filePath := "./data/todos.json"
	var toDos []model.ToDo

	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "not found"
	}
	slog.Info("Starting Create..with ", "description", description, "status", status, "traceID", traceID)

	// Load todos from json file
	toDos = loadAll(filePath)

	// Add the new todo item to the array
	var toDo model.ToDo
	toDo.Id = utils.GetNextId(toDos)
	toDo.Status = status
	toDo.Description = description

	// add new item to array
	toDos = append(toDos, toDo)

	// save the revised data here
	saveAll(toDos)

	slog.Info("Create completes after saving ", "record", toDo)

}

// Update a single ToDo item and persist back to file, make sure
// incoming ToDo item has a sensible Id otherwise bail out
func Update(ctx context.Context, toDo model.ToDo) {
	// check for uninitialised Id
	if toDo.Id <= 0 {
		slog.Info("ToDo Id uninitialised - no action taken")
		return
	}

	filePath := "./data/todos.json"
	var toDos []model.ToDo

	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "not found"
	}

	slog.Info("Starting Update..for ", "record", toDo, "traceID", traceID)

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
		// In place array update using currindx
		toDos[currIndx].Description = toDo.Description
		toDos[currIndx].Status = toDo.Status
		// persist back to file
		saveAll(toDos)
		slog.Info("Update for ", "record", toDo, "status", "completed")
	} else {
		slog.Warn("Update not run as record id, cannot be located - no action taken ", "ID", toDo.Id)
	}
	slog.Info("Update completes")

}

// delete a record by id with check to ensure id is sensible
func Delete(ctx context.Context, Id int) {

	// check for uninitialised Id
	if Id <= 0 {
		slog.Info("Id uninitialised - no action")
		return
	}

	var toDos []model.ToDo
	filePath := "./data/todos.json"

	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "not found"
	}

	slog.Info("Starting Delete for ID", "ID", Id, "traceID", traceID)

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
		var newToDos = append(toDos[:currIndx], toDos[currIndx+1:]...)
		// persist back to file
		saveAll(newToDos)
		slog.Info("Delete for id complete", "ID", Id)
	} else {
		slog.Warn("Delete not run as record id, cannot be located - no action taken ", "ID", Id)
	}
	slog.Info("Delete completes")
}

// saves all items in ToDo array back to the specified json file
// leave as private
func saveAll(todos []model.ToDo) {
	slog.Info("Starting saveAll")
	filePath := "./data/todos.json"

	// Marshal the struct into JSON
	jsonData, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		slog.Error("Error marshalling JSON: ", "error", err)
		return
	}
	slog.Info("Marshalling completed")

	// open up a new file....
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("Error creating file ", "error", err)
		return
	}
	defer file.Close()
	slog.Info("File Created ", "file", filePath)

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		slog.Error("Error writing to file ", "error", err)
		return
	}

	slog.Info("File Written ok, SaveAll completes")
}

// Loads all items in from the specified file and returns an array of ToDo Items
// leave as private
func loadAll(filePath string) []model.ToDo {

	slog.Info("Starting loadAll from file", "file", filePath)

	// is of type ToDos
	var todos []model.ToDo

	// Check file existence
	jsonFileExists := utils.CheckFileExists(filePath)

	if !jsonFileExists {
		slog.Info("Unable to open file, an empty array will be returned ", "file", filePath)
		return todos
	}
	slog.Info("File located", "file", filePath)

	data, err := os.Open(filePath)
	// Handle any error
	if err != nil {
		slog.Error("Unable to open file", "error", err)
	}
	slog.Info("Successfully Opened todos.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer data.Close()

	// Read the file in....
	byteValue, err := io.ReadAll(data)
	if err != nil {
		slog.Error("loadAll - Unable to ReadAll", "error", err)
	}
	slog.Info("Successfully Read bytevalue")

	// unmarshal the json...
	err = json.Unmarshal(byteValue, &todos)
	if err != nil {
		slog.Error("Unable to marshal JSON", "error", err)
	}

	slog.Info("To Do's contains " + strconv.Itoa(len(todos)) + " items ")
	// hand back our result..
	slog.Info("loadAll completes...")
	return todos
}
