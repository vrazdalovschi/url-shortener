[![Go Report Card](https://goreportcard.com/badge/vrazdalovschi/url-shortener)](https://goreportcard.com/report/vrazdalovschi/url-shortener)
![Go](https://github.com/vrazdalovschi/url-shortener/workflows/Go/badge.svg?branch=main)

# URL Shortener Service

URL Shortener service in Go. Functionality:
* Create/Describe/Delete shortened urls 
* Redirect shortened urls to Original
* Check statistic: redirections count, lastRedirectTime


## Run
Run with Docker (deployments folder)

```shell script
cd deployments
docker-compose up -d --build
```

## Assumptions and expected environment
Project is developed in GO 1.14 an instance with go 1.14 is expected

## Configuration options
Available using first flags if not set then env variable 

```
HTTP_ADDR    -  Endpoint port (default value :8080)
```

#### Postgres related
```
DB_HOST       - Database Host. Default: localhost
DB_PORT       - Database Port. Default: 5432
DB_USER       - Database User. Default: url-shortener
DB_PASSWORD   - Database Password. Default: root
DB_NAME       - Database Name. Default: shortener
```

## Storage
Postgres is chosen for storage purposes.

### Schema definition
Tables are auto created on starting the project.  

#### Urls table.
```sql
CREATE TABLE IF NOT EXISTS url 
  ( 
     shortenedid    VARCHAR NOT NULL UNIQUE, 
     originalurl    VARCHAR NOT NULL, 
     apikey         VARCHAR NOT NULL, 
     creationtime   TIMESTAMP NOT NULL, 
     expirationdate TIMESTAMP NOT NULL 
  );
```

#### Statistics table.
```sql
CREATE TABLE IF NOT EXISTS stats 
  ( 
     shortenedid VARCHAR NOT NULL UNIQUE, 
     redirects   INTEGER DEFAULT 0, 
     visitdate   TIMESTAMP 
  ); 
```

## Known issues
* Test Coverage
* Better error handling
* Distributed caching
* SQL Schema migration

## Api Examples
Api description is located in api folder. 

### Api Definition
```http request
http://localhost:8080/swaggerui/swagger.yaml
```
### SwaggerUi
```http request
http://localhost:8080/swagger/index.html
```

## Instrumentation
Prometheus metrics:

* System load
* Latency summary per each api call (success/error)
* Counter for each api call (success/error)

```http request
http://localhost:8080/internal/metrics
```
