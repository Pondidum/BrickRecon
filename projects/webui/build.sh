#! /bin/sh

ln ../variables.dev.json src/variables.dev.json

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

npm run build --prefix $DIR
