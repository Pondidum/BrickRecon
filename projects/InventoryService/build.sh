#! /bin/sh

npm run build

pushd ./build/app

cp ../../package.json .
npm install --only prod

popd
