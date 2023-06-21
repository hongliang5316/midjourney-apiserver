List
====

* [midjourney-apiserver](#midjourney-apiserver)
* [how to use](#how-to-use)
* [how to deploy](#how-to-deploy)
* [status](#status)
    * [imagine api](#imagine-api)
    * [upscale api](#upscale-api)
    * [describe api](#describe-api)
* [todo](#todo)
    * [variation api](#variation-api)
    * [blend api](#blend-api)
    * [progress api](#progress-api)
    * [slash command api](#slash-command-api)

# midjourney-apiserver #

`midjourney-apiserver` is an unofficial API service for `Midjourney`, designed to integrate the powerful capabilities of `Midjourney` into one's own business.

## How to use ##

`midjourney-apiserver` only provides the `grpc` protocol with a default port of `8080`.

Please refer to the [api.proto](./pkg/api/api.proto) for more information.

Here are some tools for debugging `grpc` that you can use: [awesome-grpc](https://github.com/grpc-ecosystem/awesome-grpc#tools).

If you are using `Golang`, you can use this code for testing:

```golang
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"github.com/hongliang5316/midjourney-apiserver/pkg/store"
	"github.com/hongliang5316/midjourney-apiserver/pkg/webhook"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var apiServiceClient api.APIServiceClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	apiServiceClient = api.NewAPIServiceClient(conn)

}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := &webhook.WebhookRequest{}
	json.Unmarshal(body, req)

	log.Printf("req: %+v", req)

	if req.Type == store.TypeImagine {
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
		resp, err := apiServiceClient.Imagine(context.Background(), &api.ImagineRequest{
			Prompt:  "a car",
			Webhook: "http://127.0.0.1:8000/",
		})
		if err != nil {
			panic(err)
		}

		log.Printf("%+v", resp)
	}()

	http.HandleFunc("/", webhookHandler)
	http.ListenAndServe(":8000", nil)
}
```

[List](#list)

## How to deploy ##

- Use `docker-compose` (Strongly Recommend)

```sh
mkdir -p /your/app/conf

cd /your/app

# configure conf.yml, please see ./conf/conf.yml.example
vim conf/conf.yml

# create docker-compose.yml
cat docker-compose.yml
version: '3.1'

services:
  midjourney-apiserver:
    image: hongliang5316/midjourney-apiserver:0.0.2
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
go install github.com/hongliang5316/midjourney-apiserver/cmd/midjourney-apiserver@v0.0.2

mkdir -p /your/app/conf

cd /your/app

cp `go env GOROOT`/bin/midjourney-apiserver .
mkdir conf
# configure conf.yml, please see ./conf/conf.yml.example
vim conf/conf.yml
# run
./midjourney-apiserver
```

[List](#list)

## Status ##

You can see which api are currently supported by looking at [api.proto](./pkg/api/api.proto), and if you have some ideas, feel free to submit issues.

### imagine api ###

![imagine.flow.svg](img/imagine.flow.svg)

### upscale api ###

![upscale.flow.svg](img/upscale.flow.svg)

### describe api ###

[List](#list)
