package utils

import (
	"testing"
	"todo/model"
)

func TestGetNextID_Sequential(t *testing.T) {

	t.Log("TestGetNextID_Sequential - starts")
	// Create dummy todos array in sequence
	toDos := []model.ToDo{
		{Id: 1, Description: "description here", Status: "COMPLETED"},
		{Id: 2, Description: "description here", Status: "COMPLETED"},
		{Id: 3, Description: "description here", Status: "COMPLETED"},
	}

	result := GetNextId(toDos)
	expected := 4

	if result != expected {
		t.Errorf("GetNextId = %d want %d", result, expected)
	}
	t.Log("TestGetNextID_Sequential - ends")
}
func TestGetNextID_Unordered(t *testing.T) {

	t.Log("TestGetNextID_Unordered - starts")
	// Create dummy todos array in sequence, sequences out of order
	toDos := []model.ToDo{
		{Id: 3, Description: "description here", Status: "COMPLETED"},
		{Id: 17, Description: "description here", Status: "COMPLETED"},
		{Id: 6, Description: "description here", Status: "COMPLETED"},
	}

	result := GetNextId(toDos)
	expected := 18

	if result != expected {
		t.Errorf("GetNextId = %d want %d", result, expected)
	}
	t.Log("TestGetNextID_Unordered - ends")

}

func TestGetNextID_EmptyToDoList(t *testing.T) {

	t.Log("TestGetNextID_EmptyToDoList - starts")

	// Create dummy todos array in sequence, sequences out of order
	var toDos []model.ToDo

	result := GetNextId(toDos)
	expected := 1

	if result != expected {
		t.Errorf("GetNextId = %d want %d", result, expected)
	}
	t.Log("TestGetNextID_EmptyToDoList - ends")
}

func TestGetNextIdTableDriven(t *testing.T) {

	t.Log("TestGetNextIdTableDriven - starts")
	// Set up our arrays
	// Unordered
	toDos_Unordered := []model.ToDo{
		{Id: 3, Description: "description here", Status: "COMPLETED"},
		{Id: 17, Description: "description here", Status: "COMPLETED"},
		{Id: 6, Description: "description here", Status: "COMPLETED"},
	}

	// Empty
	var toDos_Empty []model.ToDo

	// Ordered
	toDos_Ordered := []model.ToDo{
		{Id: 1, Description: "description here", Status: "COMPLETED"},
		{Id: 2, Description: "description here", Status: "COMPLETED"},
		{Id: 3, Description: "description here", Status: "COMPLETED"},
	}

	ptrUnordered := &toDos_Unordered
	ptrOrdered := &toDos_Ordered
	ptrEmpty := &toDos_Empty

	tests := []struct {
		name     string
		todo     *[]model.ToDo
		expected int
	}{
		{"Empty Array", ptrEmpty, 1},
		{"Ordered", ptrOrdered, 4},
		{"Unordered", ptrUnordered, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Actioning test %s", tt.name)
			result := GetNextId(*tt.todo)
			if result != tt.expected {
				t.Errorf("GetNextId = %d want %d", result, tt.expected)
			}
		})
	}
	t.Log("TestGetNextIdTableDriven - ends")
}

func TestCheckFileExists_ReturnType(t *testing.T) {

	// Should only return a boolean
	t.Log("TestCheckFileExists_ReturnType - starts")
	filePath := "./data/todos.json"
	result := CheckFileExists(filePath)
	var expected bool

	// SHould only get a boolean irrespective of file name
	if result != expected {
		t.Errorf("TestCheckFileExists_ReturnType = %T want %T", result, expected)
	}

	t.Log("TestCheckFileExists_ReturnType - ends")
}

func TestValidateStatusTableDriven(t *testing.T) {

	// Should only return a boolean
	t.Log("TestValidateStatus - starts")

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Valid_Completed", "COMPLETED", true},
		{"Valid_NotStarted", "NOT STARTED", true},
		{"Valid_Started", "STARTED", true},
		{"Invalid", "random junk", false},
		{"zero length", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Actioning test %s", tt.name)
			result := ValidateStatus(tt.value)
			if result != tt.expected {
				t.Errorf("ValidateStatus = %t want %t", result, tt.expected)
			}
		})
	}
	t.Log("TestValidateStatus - ends")
}
