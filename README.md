Simple read me for go academy todo app

Intial take on ToDo app done

Supported flags are :-  
    action      - Create / Show / Update / Delete (Case insensitive)
    status      - default to empty string - relevant for create action
    description - default to empty string - relevant for create action
    Id          - only relevant for update/delete actions

Show / Create / Delete / Update implemented

All data located in data/todos.json

Example usage :- 

go run todo.go -action "DELETE" -Id nn - Delete specific id nn
if id <= 0 no action will be taken, if id cannot be located no action will be taken
Note :- Case insensitive flag

go run todo.go -action "show" 
List out all elements in the data/todos.json file

Updates:-

Cleaned up Repo

Cleaned up folder structure

Moved type to own 'datatypes' package

Split code into packages datatypes, utils, constants, dataaccess, cli, tests (Not implemented as yet)

Increased error handling

Amended GetNextId to use sort as per recommendations

Added slog & context trace id
