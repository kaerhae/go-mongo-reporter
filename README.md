# go-mongo-reporter
Simple Reporting Tool. 

Backend designed to use MongoDB. Backend written in Golang.

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
`apt-get install /path/to/package.deb`

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

