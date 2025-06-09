package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"todo/constants"
	"todo/datatypes"
)

/*
	Helper functions below this comment line
	GetNextId & CheckfileExists
	Note - Capitalization of function names will export the functions correctly
*/
// Scans the array and returns highest id int + 1
func GetNextId(toDos []datatypes.ToDo) int {
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

func ShowPermittedStatuses() {
	fmt.Printf("Permitted statuses are :- %s, %s, %s\n", constants.StatusCompleted, constants.StatusNotStarted, constants.StatusStarted)

}
