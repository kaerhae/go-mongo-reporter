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
DATABASE_URI=
SECRET_KEY=
DATABASE=
```
### Debian
.deb package can be created locally with:
```bash
export VERSION=x.x.x
export BUILD_VERSION=x
cd scripts && ./create-deb-package.sh
```

To install package, run:
`apt-get install /path/to/package.deb`

### Fedora
.rpm package can be created locally with:
```bash
export VERSION=x.x.x
export BUILD_VERSION=x
cd scripts && ./create-rpm-package.sh
```

To install package, run:
`dnf install /path/to/package.rpm`

