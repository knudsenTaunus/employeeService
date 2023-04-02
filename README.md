# User Service

This service provides simple CRUD Operations for Users which have the following format:   



```json
{
    "id": "8f29039b-7603-42fc-954a-24f7b8c7f548",
    "first_name": "first",
    "last_name": "last",
    "nickname": "nick",
    "password": "pass",
    "email": "example@example.org",
    "country": "UK"
}
```
The Endpoints for the managing users are the following:   

``` 
/user       - Method POST
/users      - Method GET
/users/{id} - Methods GET PATCH DELETE
```

## Running the Service

To run the service please make sure you installed all dependencies with:   

```go get ./... ```

The service uses a MySQL Database which can be configured through the   
```config.yml``` file. Further you can configure host and port for both, the HTTP   
and the gRPC Server.

For development purposes there is a ```docker-compose.yml ``` which provides a MySQL   
database. Please note that the current image is for Computers with ARM Chips like the   
Apple M1. If you use another Architecture, please adjust the image.

After you started the MySQL database, you need to execute the migrations which can be   
found in the ```migrations``` folder. To do that, please use the Makefile Command   
```make migrate```.

After that you can execute the service with this command:

``` go run cmd/userService/main.go```

The service offers an HTTP Interface for general usage and also a gRPC Interface   
on that a client can register to receive instant information when a user got updated.

An example client can be found in the ```grpcclient``` folder.   
The protobuf definitions can be found under ```proto``` folder, for generating server   
and client code a tool called ```buf``` which can be used by executing the Makefile Command   
```make protobuf ```
