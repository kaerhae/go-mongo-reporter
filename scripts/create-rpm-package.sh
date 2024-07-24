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

go build -o ../bin/reporter/reporter ../cmd/reporter/main.go

# RPM PACKAGE
RPMPACKAGE=reporter-${VERSION}
mkdir -p $RPMPACKAGE
cp ../bin/reporter/reporter $RPMPACKAGE
cp ../reporter.service $RPMPACKAGE
tar -cvzf $RPMPACKAGE.tar.gz $RPMPACKAGE/

# RPMBUILD
rpmdev-setuptree
cp $RPMPACKAGE.tar.gz ~/rpmbuild/SOURCES
rpmbuild -ba ../build/rpm/reporter.spec
if [ $? -ne 0 ]; then { echo "Build failed, aborting." ; exit 1; } fi

cp ~/rpmbuild/RPMS/x86_64/$RPMPACKAGE-$BUILD_NUMBER.el9.x86_64.rpm ../build/generated-packages
# CLEANUP
rm ../bin/reporter/reporter
rm -rf $RPMPACKAGE 
rm $RPMPACKAGE.tar.gz

cd ../build/generated-packages
echo ""
echo "Build succesful, new RPM package created to: $PWD"
echo ""