package cli

// imports
import (
	"context"
	"log"
	"log/slog"
	"strings"
	"todo/dataaccess"
	"todo/model"
	"todo/utils"
)

func StartCLI(ctx context.Context, action string, status string, description string, Id int) {

	// actions - Create / Show / Update / Delete
	slog.Info("Selected action", "action", strings.ToLower(action))
	switch strings.ToLower(action) {
	case "show":
		// show all records no params needed
		dataaccess.ShowAllRecords(ctx)
	case "create":
		if !utils.ValidateStatus(status) {
			utils.ShowPermittedStatuses()
			log.Fatalf("Status of %s is not permitted", status)
			return
		}
		// All ok create new todo item
		dataaccess.Create(ctx, description, status)
	case "update":
		var updatedToDo model.ToDo
		updatedToDo.Id = Id
		updatedToDo.Description = description
		updatedToDo.Status = status
		if !utils.ValidateStatus(status) {
			utils.ShowPermittedStatuses()
			log.Fatalf("Status of %s is not permitted", status)
			return
		}
		// all ok update may continue
		dataaccess.Update(ctx, updatedToDo)
	case "delete":
		// delete function
		dataaccess.Delete(ctx, Id)
	default:
		log.Printf("Unsupported action..%s passed to procedure\n", strings.ToLower(action))
	}
}
