#!/bin/sh

env_setup(){
    PROTO_PATH=$GOPATH/src/auservices/api/proto
    API_PATH=$GOPATH/src/auservices/api/
    TMP_PATH=$GOPATH/src/auservices/tmp
    FRONTEND_API_PATH=~/Workspace/UmaWorkspace/au/services/api
}

env_setup

rm -Rf $TMP_PATH/*
protoc api/proto/*.proto -I. --go_out=plugins,import_path=api:$TMP_PATH

ret_val=$?

if [ $ret_val -ne 0 ]; then
    echo FAIL
    exit 1
fi

rm -f $API_PATH/*.pb.go
cp $TMP_PATH/api/proto/*.go $API_PATH
cp -f $PROTO_PATH/*.proto $FRONTEND_API_PATH

echo OK

