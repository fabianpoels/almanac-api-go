# Almanac API

The almanac api is written in Go and uses Gin as an API framework.

## Setting up a dev environment

The project is setup to use [VSCode devcontainers](https://code.visualstudio.com/docs/devcontainers/containers), so the only requirements are VSCode + Docker.

Following services are contained within the devcontainer. Their respective ports are also exposed on the container, so make sure they are not in use (or update `docker-compose.yml`)
* Go/Gin API :3000
* Mongodb database :27017
* Redis cache :6379

### 1 Set environment variables

Copy the example `cp .env.example .env` and update the required values.

### 2 Configure gitlab SSH key

The project needs access to [https://gitlab.com/almanac-app/models](https://gitlab.com/almanac-app/models) during the go modules installation. Make sure you have configured ssh access, as the docker container uses ssh passthrough for downloading the repo.

### 3 Build the devcontainer

Use the VSCode command to build the devcontainer `Dev Containers: Rebuild and reopen in container`

### 4 Start the API

Inside the container, run `air` to start the API with automatic code reloading enabled.

### 5 Import data

For starting with a filled database, import some data through the host machine:
`mongorestore --uri="mongodb://localhost:27017" --gzip --archive="$1" --nsFrom="almanac-go.*" --nsTo="almanac-go-prod.*" --drop`
