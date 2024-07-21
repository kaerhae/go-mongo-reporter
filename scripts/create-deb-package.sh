#!/bin/bash

if [[ -z $VERSION ]]; then
    echo $VERSION
    echo "VERSION environment variable is not set"
    exit 1
fi

if [[ -z "$BUILD_NUMBER" ]]; then
    echo "BUILD_NUMBER environment variable is not set"
    exit 1
fi

# TEST
go test -v ../...

# BUILD
export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=linux

go build -o reporter ../.

# DEB
ARCH=amd64
DEBPACKAGE=reporter_${VERSION}-${BUILD_NUMBER}_${ARCH}
mkdir -p $DEBPACKAGE

mkdir -p $DEBPACKAGE/usr/local/bin
mkdir -p $DEBPACKAGE/etc/systemd/system
mkdir -p $DEBPACKAGE/etc/opt/reporter/config


cp reporter $DEBPACKAGE/usr/local/bin
cp ../reporter.service $DEBPACKAGE/etc/systemd/system
cp -r ../build/deb/DEBIAN $DEBPACKAGE

dpkg --build $DEBPACKAGE
echo ""
echo "-------READING .deb PACKAGE INFO----------"
dpkg --info $DEBPACKAGE.deb
echo ""
echo "-------READING .deb PACKAGE CONTENTS----------"
dpkg --contents $DEBPACKAGE.deb

mkdir -p ../build/generated-packages
mv -f $DEBPACKAGE.deb ../build/generated-packages


# CLEANUP

rm -rf $DEBPACKAGE
rm -rf usr/local/bin
rm -rf etc/systemd/system
rm -rf etc/opt/reporter/config
rm reporter