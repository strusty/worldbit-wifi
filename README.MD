# Prerequisites
* [freeRADIUS](https://freeradius.org/releases/) **IMPORTANT: The backend requires Postgres database to be used for RADIUS**
* Postgres v≥9.0
* Golang v≥1.9.2
* [dep](https://github.com/golang/dep) latest


# How to start
* Clone the project to your $GOPATH
* 'dep ensure' all the dependencies
* Fill out config.yaml
* `go run main.go`

# How to use api_docs.yaml
Navigate to [Swagger Editor](https://editor.swagger.io) and import the file.