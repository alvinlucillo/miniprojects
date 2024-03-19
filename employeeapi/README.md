### Employee CRUD API

#### Running the server

- Run at the project root: `go run ./cmd/`
- Access the API at `http://localhost:9000`
- See sample requests in `test.http`

#### Running the tests

`go test -count=1 ./...`

#### Project structure

```
├── cmd
│   └── main.go                         -> entrypoint of the server
├── go.mod
├── go.sum
├── internal
│   ├── handlers                        -> contains the handlers for the endpoints
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── repos                           -> contains the repository layer objects
│   │   ├── repos.go
│   │   └── repos_test.go
│   └── services                        -> contains the service layer objects
│       └── services.go
└── test.http
```

#### Go packages used

- `envconfig` - environment variable parsing
- `gin-gonic/gin` - web framework
- `go-playground/validator/v10` - used for data validation
- `testify` - testing framework
- `zerolog` - used for logging
