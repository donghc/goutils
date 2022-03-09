#!/bin/bash

app=$1
apps=( "csscand" "csmeta" "csfront" "csinterface" "cspost" "csagg" "cswatch" "csadmin" "metric")

function prebuild() {
#    FULL_Version=`git tag -l "$APPNAME-*"|tail -n 1`
    Version=`git describe $(git rev-list --tags --max-count=1)`
    #Author=`git show ${Version} | grep 'Tagger' |  cut -d ' ' -f2`
    Author=`git config --get user.name`
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
#    GCFLAGS="-N -l"
    #GCFLAGS="'-s -w'"
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

app=$1

function build() {
app=$1
echo "building $app..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "$LDFLAGS" -o "$PWD/build/$app" "$PWD/main/$app/."
if [ "$app" == "csscand" ]; then
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build --ldflags "$LDFLAGS" -o $PWD/build/csscand.exe $PWD/main/csscand/.
fi
}

if [ "x$app" == "x" ]; then
  for app in ${apps[@]}
  do
     build $app
  done
elif [[ "${apps[@]}" =~ "${app}" ]]; then
  build $app;
elif [[ ! "${apps[@]}" =~ "${app}" ]]; then
  echo "plz input right app: ${apps[@]}"
fi

echo "build end"
