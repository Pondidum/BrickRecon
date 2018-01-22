#! /bin/sh

npm run build

pushd ./build/app

cp ../../package.json .
sed -i "s/file:..\//file:..\/..\/..\//g" package.json

npm install --only prod

popd
