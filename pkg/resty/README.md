# resty

`resty` is a generic JSON REST API client package.

## Examples

```go
client := resty.NewClient("http://todolist.com")

// call Get & Post and handle response yourself
bytes, err := client.Post("/api/todo", &Todo{Name: "Take the bins out", Due: time.Now()})
bytes, err := client.Get("/api/todos")

// use generic functions to marshal/unmarshal the response
saved, err := resty.Post[Todo, Todo](client, "/api/todo", &Todo{Name: "Take the bins out", Due: time.Now()})
todos, err := resty.Get[[]Todo](client, "/api/todos")
```
