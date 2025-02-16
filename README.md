# Introduction
This repository is to scan, store and query json vulnerability scan files from https://github.com/velancio/vulnerability_scans

## Setup
1. Clone the GitHub repository: `git clone https://github.com/sarthakgoel1997/parse-github-files.git`

2. Download Docker: https://docs.docker.com/desktop/setup/install/mac-install/

3. Go to `Makefile` and update `PERSONAL_ACCESS_TOKEN` with your GitHub personal access token to query repositories

3. Go to the root of the repository and run: `make dev`

## Useful Docker Commands
1. `make build`: Builds the docker image

2. `make run`: Runs a docker container with the built image

3. `make stop`: Stops and deletes any running container

4. `make dev`: Builds the docker image, stops any running containers and starts up a new container

5. `make logs`: Starts up docker container logs for debugging

6. `make query-db`: Opens sqlite database in the terminal

7. `make clean`: Stops any running containers and deletes the built docker image

## Useful SQLite Commands
1. `.tables`: Displays all tables present in the database

2. `PRAGMA table_info (<table_name>)`: Displays all columns and types present in the table

3. `.mode line`: To display select query results in a readable format
