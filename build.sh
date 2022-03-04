#!/bin/bash


function prebuild() {
#    FULL_Version=`git tag -l "$APPNAME-*"|tail -n 1`
    Version=`git describe $(git rev-list --tags --max-count=1)`
    Author=`git show ${Version} | grep 'Tagger' |  cut -d ' ' -f2`
    BuildDate=$(date "+%Y-%m-%dT%H:%M:%S")
    # git describe $(git rev-list --tags --max-count=1)
#    FULL_Version=`git  branch | grep '*' | sed -e 's/\*//g' -e 's/HEAD detached at//g' -e 's/\s*//g' -e 's/[\(\)]//g'`
    if [ x$Version == "x" ]; then
        echo "use version 1.0.0"
        Version="1.0.0"
    fi
#    REVISION=`git rev-parse --short HEAD`
#    echo "prepare build:  $Version#$APPNAME"
#    REVISION=`git rev-parse --short HEAD`
    #LDFLAGS="-X main.APPVersion=${Version} -X main.REVISION=${REVISION} -X main.APPNAME=${APPNAME}"
    #LDFLAGS="-s -w -X main.APPVersion=${Version}  -X main.REVISION=${REVISION} -X main.APPNAME=${APPNAME} -extldflags '-static'"
    LDFLAGS=" -X main.Version=${Version}  -X main.Author=${Author} -X main.BuildDate=${BuildDate}"

    ##generate resource first
}
prebuild
echo "Version="${Version}
echo "Author="${Author}
echo "BuildDate="${BuildDate}
echo "LDFLAGS="${LDFLAGS}
echo "build start"
rm -rf build/*
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csadmin $PWD/main/csadmin/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csagg $PWD/main/csagg/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csfront $PWD/main/csfront/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csinterface $PWD/main/csinterface/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csmeta $PWD/main/csmeta/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/cspost $PWD/main/cspost/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csscand $PWD/main/csscand/.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/cswatch $PWD/main/cswatch/.
#go build --ldflags "$LDFLAGS" -o $PWD/build/csscand.exe $PWD/main/csscand/.
echo "build end"
