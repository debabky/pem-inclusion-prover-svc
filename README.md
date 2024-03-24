# PEM Inclusion Prover Service

## Description

A service that generates PEM Merkle trees and inclusion proofs.

## Install

  ```
  git clone github.com/debabky/pem-inclusion-prover-svc
  cd pem-inclusion-prover-svc
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/). Here is how you can start it:
```
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

## Running from docker 
  
use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t github.com/debabky/pem-inclusion-prover-svc .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/debabky/pem-inclusion-prover-svc
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `run service` command

## IPFS

The service utilizes IPFS to save Merkle trees. Just provide a valid url to IPFS in the `config.yaml`.
