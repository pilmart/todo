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
Split code into packages datatypes, utils, constants, dataaccess
Increased error handling
Amended GetNextId to use sort as per recommendations
Added Context & traceId across data access code

Server code created - go run todo.go

HTML based Endpoints are :-
/todolist  - simple html list of json file
/about - static html page

Rest API endpoints specified as :-

GET http://localhost:3000/todo/:id - return single todo item by :id

DELETE http://localhost:3000/todo/:id - Delete record with specified id

PUT http://localhost:3000/todo - Update record - requires JSON payload as 
{
    "id" : 15,
    "description":"A new record test",
  "status":"STARTED"
} 

POST http://localhost:3000/todo - Requires JSON payload as 
{
  "description":"A new record test",
  "status":"STARTED"
} 
Note :- A new id will be generated

Actor pattern implemented and code split out of MUX.go file and moved to todoactor.go

Parallel tests set up for get/put/post, delete still outstanding
