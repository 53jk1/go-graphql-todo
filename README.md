# GraphQL Todo List
This is a simple project to demonstrate how to create a GraphQL API using Go. It implements a todo list with basic CRUD functionality.

## Requirements
- Go 1.16 or higher
- github.com/99designs/gqlgen
## Installation
1. Clone the repository:
```bash
git clone https://github.com/your-username/your-repository.git
```
2. Install dependencies:
```bash
go mod download
```
## Usage
1. Start the server:
```bash
make run-dev
```
2. Open your web browser and go to http://localhost:8080/playground to access the GraphQL Playground.
3. Use the Playground to create, update, delete, and retrieve todos.

## Example queries
### Create a new todo
```graphql
mutation {
    createTodo(title: "Buy groceries", description: "Milk, eggs, bread") {
        id
        title
        description
        isCompleted
    }
}
```
### Update a todo
```graphql
mutation {
    updateTodo(
        id: "1"
        title: "Buy groceries"
        description: "Milk, eggs, bread, cheese"
        isCompleted: true
    ) {
        id
        title
        description
        isCompleted
    }
}
```
### Delete a todo
```graphql
mutation {
    deleteTodo(id: "a219e689-5f2c-4066-92cf-1059e01fb1aa")
}
```
### Retrieve all todos
```graphql
query {
    todos {
        id
        title
        description
        isCompleted
    }
}
```
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.