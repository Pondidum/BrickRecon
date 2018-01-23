#! /bin/bash

TARGET="${PWD##*/}"

pushd ../.packages

npm pack "../$TARGET"

popd
