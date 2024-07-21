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

sed -i "s/{VERSION}/${VERSION}/g" ../build/rpm/reporter.spec
sed -i "s/{BUILD_NUMBER}/${BUILD_NUMBER}/g" ../build/rpm/reporter.spec

# TEST
go test -v ../...

# BUILD
export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=linux

go build -o reporter ../.

# RPM
RPMPACKAGE=go-mongo-reporter-${VERSION}
mkdir -p $RPMPACKAGE

cp reporter $RPMPACKAGE
cp ../reporter.service $RPMPACKAGE

tar -cvzf $RPMPACKAGE.tar.gz $RPMPACKAGE/

rpmdev-setuptree
cp $RPMPACKAGE.tar.gz ~/rpmbuild/SOURCES
rpmbuild -ba ../build/rpm/reporter.spec

cp ~/rpmbuild/RPMS/x86_64/$RPMPACKAGE-$BUILD_NUMBER.el9.x86_64.rpm ../build/generated-packages
rm reporter
rm -rf $RPMPACKAGE 
rm $RPMPACKAGE.tar.gz