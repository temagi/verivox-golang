### A very simple implementation of test task with golang (instead of Kotlin / JVM)

## Requirements

### For local run required:

* Golang 1.21+ (installed with proper GOHOME environment variable setted)

### for run in container:   

* Docker

## How to run

### Local run

```
go test
```

### Build docker image and run in container:

```
docker build -t verivox-go .
docker run -it verivox-go
```

### Run in docker container without building image:

```
docker run -it -v  .:/app -w /app golang:1.21-alpine go test
```

