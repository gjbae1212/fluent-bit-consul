#!/bin/bash
set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT

export GO111MODULE=on
cd `dirname $0`
CURRENT=`pwd`

function build
{
   go build -buildmode=c-shared -o consul.so .
   local osname=`go env | grep GOOS | awk -F "=" '{print $2}' | sed 's/\"//g'`
   local archname=`go env | grep GOARCH | awk -F "=" '{print $2}' | sed 's/\"//g'`
   mkdir -p $CURRENT/bin/${osname}_${archname} || true
   mv consul.so consul.h $CURRENT/bin/${osname}_${archname}/
}

function build_linux
{
   local plugin=`docker ps | grep consul-plugin | wc -l`
   if [ $plugin -eq 1 ]
   then
      docker kill consul-plugin
   fi
   local osname=linux
   local archname=amd64
   mkdir -p $CURRENT/bin/${osname}_${archname} || true
   docker build --no-cache -t consul-plugin:latest -f Dockerfile .
   docker run -it --rm -d --name consul-plugin consul-plugin:latest /bin/bash
   docker cp consul-plugin:/go/consul/consul.so $CURRENT/bin/${osname}_${archname}/
   docker cp consul-plugin:/go/consul/consul.h $CURRENT/bin/${osname}_${archname}/
   docker kill consul-plugin
}

CMD=$1
shift
$CMD $*