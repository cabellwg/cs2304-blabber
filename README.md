# `cs2304-blabber`
Semester project for CS2304 Docker Containerization. RESTful API service written in Go with PostgreSQL. Includes pre-packaged React front-end and a Traefik reverse proxy for routing.

## Details

Blabber is a mock version of Twitter (with, of course, some missing features). The purpose of the project is to use [Docker Compose](https://docs.docker.com/compose/overview/) to manage multi-container applications.

## Usage

To run the application, you must create a file `.env` in the top-level directory of the project. This file will specify the mode in which to run the project.

There are two modes, development and production. Production mode runs a lightweight production API container. Development mode uses [`fresh`](https://github.com/gravityblast/fresh) to automatically rebuild the Go project whenever the source changes by mounting a volume to the `api` container.

To specify the mode in which to run the app, set the `MODE` parameter in the `.env` file as follows:

```
MODE=[ dev | prod ]
```

Then `docker-compose up` will run the application in the specified mode.
