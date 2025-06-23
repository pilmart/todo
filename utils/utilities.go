package utils

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"todo/constants"
	"todo/model"
)

/*
	Helper functions below this comment line
	GetNextId, CheckfileExists, ValidateStatus
	Note - Capitalization of function names will export the functions correctly
*/
// Scans the array and returns highest id int + 1, use sort to put them in desc order
// then pick the id of the first element
func GetNextId(toDos []model.ToDo) int {
	if len(toDos) == 0 {
		// no elements so set initial value, not an error state as we are
		// will create a new json file starting from element id 1
		return 1
	}
	sort.Slice(toDos, func(i, j int) bool {
		return toDos[i].Id > toDos[j].Id
	})
	return toDos[0].Id + 1
}

// Make sure the file exists, probably better ways to do this
func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

// validate that the incoming status is one of the statuses in constants package
// Cases MUST match ie Completed != COMPLETED
func ValidateStatus(status string) bool {

	return status == constants.StatusCompleted ||
		status == constants.StatusNotStarted ||
		status == constants.StatusStarted

}

// Show all the permitted statuses
func ShowPermittedStatuses() string {
	return fmt.Sprintf("Permitted statuses are :- %s, %s, %s\n", constants.StatusCompleted, constants.StatusNotStarted, constants.StatusStarted)
}

// scan the slice of toDos for id and return index of matching
// element or return -1
func FindTodoIndex(id int, toDos []model.ToDo) int {

	for i, toDo := range toDos {
		if toDo.Id == id {
			return i
		}
	}
	return -1
}
