# GoNCS: Go Networking Container Suite

GoNCS is a suite of modular networking functions built in Docker containers. The containers include a reverse proxy, dns, logging manager, and a dashboard. These containers communicate with each other through network requests in order to stay independent of one another and stay isolated on the network.

### Prerequisites

Once you have cloned the project, run `go get` to retrieve all the dependencies you will need for development.

```
git clone https://github.com/mattackard/project-1.git
cd $GOPATH/src/github.com/mattackard/project-1/
go get ./...
```

Docker build will automatically get the dependencies for your containers during the build process, but having the dependencies locally will provide your IDE with code sensing and prevent any warnings.

### Installing

To get started working with the source code, you will need to create a .env file in the root project directory containing your DNS container's address for requests, and your authentication credentials being used in the reverse proxy. Once the env has been created you can run `make all` to build and run the docker containers.

```
make all flag=docker-compose_build_flag_here
```

This command will run `docker-compose build` to build all the containers using a build stage to create a lightweight container with only the go binary. It will then run `docker-compose up` to run all your containers.

Once your docker containers are running, you will see log files are generating outlining the communication between each of the containers.

## Deployment

To deploy the containers onto a live system, create an image containing only the go binary. From there follow the steps outlined in the [docker documentation](https://docs.docker.com/compose/production/) to deploy the docker compose application.

## Built With

- [Docker-Compose](https://docs.docker.com/compose/) - Containerization
- [Procstats](github.com/segmentio/stats/procstats) - Process statistics
