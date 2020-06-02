# Simple REST API

## Author
Andrei Milas

## Description
A basic REST API, connected to a MariaDB database,  providing access to a ```user``` resource 

#### Endpoints 
```
GET /v1/users
```
* Returns a list of users in a JSON array
* Supports pagination (e.g. GET /v1/users?limit=1&offset=2)
```
POST /v1/users
```
* Creates a user instance
```
GET /v1/users/{uuid}
```
* Returns a user instance in JSON format
```
DELETE /v1/users/{uuid}
```
* Deletes a user instance

## Running the project
1. Check the configuration in **config.yml** and adapt it to your environment.
2. You need to make sure the MySQL server is accepting connections and has loaded initial data. You can do this by running ```docker-compose up``` in the root folder of the project. This will start the MySQL server in a docker container (with port 3306 forwarded) and create the database and the required table from the SQL script in **./dumps**.
3. Build the project by running ```go build -o sample-rest-api```.
4. Run the project: ```./sample-rest-api```

You can now perform API calls to the endpoints listed above.

## Testing

In order to run the unit tests, you can call ```go test -v ./...``` in the root folder of the project.