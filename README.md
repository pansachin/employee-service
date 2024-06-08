# Employee Service

Employee service deals with the employee data. It provides the following functionalities:
- Create an employee
- Get an employee
- Get all employees
- Update an employee
- Delete an employee(Soft delete)

## Tools/Software used:
- Taskfile:
    + https://taskfile.dev/
    + Task is a task runner / build tool that aims to be simpler and easier to use than, for example, GNU Make.
- Mysql:
    + For storing results/data inside database.
- Flyway:
    + https://documentation.red-gate.com/fd/flyway-documentation-138346877.html
    + Flyway is an industry leading database versioning framework that aims to unlock DevOps for the database. It strongly favors simplicity and convention over configuration.
    + We use community edition.
- Docker Desktop:
    + https://www.docker.com/products/docker-desktop/


## Local Setup

### Prerequisite
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) must be running on the system.


### Local Setup
```bash
task local-setup
```
This will cerate a MySql database as docker containers and generate all db schemas using flyway tool.

### Run Local Package Index Service
```bash
task run:service
```
This will run employee service  locally.


### Local Clean
```bash
task local-clean
```
Cleans up db docker containers.

### Local Unit Test
```bash
task test-unit
```
This will run unit tests for all the modules.

### Local Lint Check
```bash
task lint-all
```
This will check the linting for the go file [https://github.com/golangci/golangci-lint]


### Local swagger document
```bash
task swagger-generate
```
This will generate the API document with example for employee-service.

### Local swagger serve
```bash
task swagger-serve
```
This will start serving Swagger documentation locally.