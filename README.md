# eq

This program download all attachements files to a discord message and save it only if it's not already saved.

## How it work

Files are downloaded in the directory under the `SRC_DIR` environment variable and are named with an UUID.  
The bot look for his connection token in the environment variable named `TKN`.

A goroutine watch files in the directory under the `SRC_DIR` environment variable and check if the file is in the directory `DST_DIR`
environment variable.  
If a file from the `SRC_DIR` is already present in the `DST_DIR` it's removed, otherwise the file is moved.

## Environment variable

- `TKN` → The discord bot token
- `SRC_DIR` → Directory where files are downloaded
- `DST_DIR` → Directory where files are unique
- `HOST` → Host where the server listen
- `PORT` → Port the server listen on
- `BASEURL` → Base URL of the static server

## Installation

### Requirements

- go version go1.17.5
- docker

### Run the bot

- Defines environment variables
- Run `make ctn-run`, it will start the container
