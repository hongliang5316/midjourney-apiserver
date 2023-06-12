# midjourney-apiserver #

`midjourney-apiserver` is an unofficial API service for "Midjourney", designed to integrate the powerful capabilities of `Midjourney` into one's own business.

## How to deploy ##

- Use `docker-compose` (Strongly recommend)

```sh
mkdir -p /your/app/conf

cd /your/app

# config conf.yml, please see ./conf/conf.yml.example
cat conf/conf.yml

# create docker-compose.yml
cat docker-compose.yml
version: '3.1'

services:
  midjourney-apiserver:
    image: hongliang5316/midjourney-apiserver:latest
    hostname: midjourney-apiserver
    restart: always
    volumes:
      - ./conf/conf.yml:/conf/conf.yml
    ports:
      - 8080:8080

  redis:
    image: redis:7
    hostname: redis
    restart: always
    volumes:
      - ./redis_data:/data
    command: redis-server --requirepass test

# run
docker-compse up -d

# check logs
docker-compose logs
```

- Manual installation

```sh
go install github.com/hongliang5316/midjourney-apiserver/cmd/midjourney-apiserver@latest

cd /your/app
cp `go env GOROOT`/bin/midjourney-apiserver .
mkdir conf
# config conf.yml, please see ./conf/conf.yml.example
cat conf/conf.yml
# run
./midjourney-apiserver
```

## How to use ##

`midjourney-apiserver` only provides the grpc protocol with a default port of 8080.

Please refer to the [api.proto](./pkg/api/api.proto) for more information.

Here are some tools for debugging grpc that you can use: [awesome-grpc](https://github.com/grpc-ecosystem/awesome-grpc#tools).

If you are using `Golang`, you can use this code for testing:

```golang
package main

import (
        "context"
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"

        "github.com/hongliang5316/midjourney-apiserver/internal/application"
        "github.com/hongliang5316/midjourney-apiserver/pkg/api"
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
)

var apiServiceClient api.APIServiceClient

func webhookHandler(w http.ResponseWriter, r *http.Request) {
        body, _ := ioutil.ReadAll(r.Body)
        req := &application.WebhookRequest{}
        json.Unmarshal(body, req)

        log.Printf("req: %+v", req)

        if req.Type == "Imagine" {
                resp, err := apiServiceClient.Upscale(context.Background(), &api.UpscaleRequest{
                        Index:   1,
                        TaskId:  req.TaskID,
                        Webhook: "http://127.0.0.1:8000/",
                })
                if err != nil {
                    panic(err)
                }

                log.Printf("resp: %+v", resp)
        }
}

func main() {
        go func() {
                http.HandleFunc("/", webhookHandler)
                http.ListenAndServe(":8000", nil)
        }()

        conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
                panic(err)
        }

        apiServiceClient = api.NewAPIServiceClient(conn)

        resp, err := apiServiceClient.Imagine(context.Background(), &api.ImagineRequest{
                Prompt:  "a car",
                Webhook: "http://127.0.0.1:8000/",
        })
        if err != nil {
                panic(err)
        }

        log.Printf("%+v", resp)

        select {}
}
```

## Status ##
