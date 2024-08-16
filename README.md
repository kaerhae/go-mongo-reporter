# go-mongo-reporter
Simple Reporting Tool. 

Backend designed to use MongoDB. Backend written in Golang. Project contains reporter app, migration tool, and also unit/integration tests and building tools. Backend is designed as REST API, which is built to support MongoDB. Exhaustive API description can be found  [here](/api/README.md)


## Build and Installation

Reporter always needs .env file to launch. In development, .env file can be located to the project root. For Debian/Fedora packages, environment file should be located on 
`/etc/opt/reporter/config/` and it should be named as `reporter.env`. 

Environment variables are following:

```
IP_ADDR= 
PORT=
SECRET_KEY=
MONGO_USER=
MONGO_PASS=
MONGO_IP=
MONGO_PORT=
DATABASE=
```
### Debian
.deb package can be created locally with:
```bash
export VERSION=x.x.x
export BUILD_VERSION=x
cd scripts && ./create-deb-package.sh
```

Alternatively, build can be done with Makefile command:
`make build-deb`

To install package, run:
`dpkg -i /path/to/package.deb`

To remove package, run:
`apt remove reporter`

### Fedora
.rpm package can be created locally with:
```bash
export VERSION=x.x.x
export BUILD_VERSION=x
cd scripts && ./create-rpm-package.sh
```
Alternatively, build can be done with Makefile command:
`make build-rpm`

To install package, run:
`dnf install /path/to/package.rpm`

To remov package, run:
`dnf remove reporter`


### Docker

Repository contains docker-compose.yml, which builds database and reporter service with initial migrations. Before running, check ./docker-compose.yml file and configure suitable environment variables. To compose dockerfiles, run:
`docker-compose up -d`

### Kubernetes

This project has deployment files to Kubernetes, which are tested on a local Minikube environment. For further instructions, read more [here](kubernetes/README.md)
