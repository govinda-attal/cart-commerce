## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Notes
Following not addressed due to lack of time.

No Business Validation Logic is implemented
No Unit Tests are written - but dependency injection is followed well so writing unit tests should be easy. Prefer BDD frameworks like ginkgo and gomega libraries

### Prerequisites

1. Linux OS - Ubuntu or Fedora is preferred.
2. git
3. Golang Setup locally on OS
4. Postman for REST API testing
5. Docker & Docker Compose

### Installing

This application can be setup locally in following ways:

#### Option A
```
go get github.com/govinda-attal/cart-commerce
```

#### Option B (Preferred) :heavy_check_mark:
```
cd $GOPATH/src/github.com/
mkdir govinda-attal
cd govinda-attal
git clone https://github.com/govinda-attal/cart-commerce.git
```

### Application Development Setup

'Makefile' will be used to setup cartcom-api quickly on the development workstation.

```
cd $GOPATH/src/github.com/govinda-attal/cart-commerce
make install # This will go get 'dep' and use it to install vendor dependencies.
```

## Running Tests

### Unit tests

Have followed principles of dependency injection well. So in theory it should be easy to write unit tests with stubs implementing the interfaces. Would prefer ginkgo and gomega.

No time for writing BDD tests though.


### Integration tests

This microservice achives given requirements with Golang and Postgres Database as backend. To keep this foot-print of this application minimum postgres db will execute within a docker container. Where as following backend microservice can be hosted within a docker container or local OS.

#### Option A: Docker Compose - orchestrate DB and Microservice as docker containers (Preferred) :heavy_check_mark:

```
cd $GOPATH/src/github.com/govinda-attal/cart-commerce
make docker-compose-start # This will start Postgres DB, cartcom-api Microservice and Swagger-UI which will point to microservice swagger definition.
```

Docker compose will orchestrate containers and they can be accessed from Local OS as below:
1. Postgres DB on localhost:5432
2. Microservice on :earth_asia: http://localhost:9080
3. Swagger-UI on :earth_asia: http://localhost:9090
4. Adminer-UI on :earth_asia: http://localhost:8080

##### Setup Database
1. Login to :earth_asia: http://localhost:8080 with details as below
```
server:db
username: postgres
password: postgres
database: postgres
```
2. Ensure you select DB as *Postgres* & Schema as *Public*
3. Click on File upload or SQL Command
4. select file ./migration/setup.sql to import or copy contents of this file and past same into SQL command text block and click execute.

##### Simlified Testing with Swagger UI
1. Browse to :earth_asia: http://localhost:9090
2. Feel free to test API - with expected happy scenarios :-) [as haven't implemented business/system integrity validation rules for lack of time]