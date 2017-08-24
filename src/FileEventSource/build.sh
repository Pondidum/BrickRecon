#! /bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_PATH="$DIR/FileEventSource/FileEventSource.csproj"

dotnet restore $PROJECT_PATH
dotnet publish $PROJECT_PATH
dotnet publish $PROJECT_PATH -o ./../.build
