#!/usr/bin/env bash
runDir=$(dirname $0)
cd "$runDir"
binName=$(basename `pwd`)

#for linux
echo "GOOS=linux go build -o release/bin/$binName main.go"
GOOS=linux go build -o "release/bin/$binName" main.go

#for mac
#echo "GOOS=darwin go build -o release/bin/$binName main.go"
#GOOS=darwin go build -o "release/bin/$binName" main.go

#for windows
#echo "GOOS=windows go build -o release/bin/$binName main.go"
#GOOS=windows go build -o "release/bin/$binName" main.go

