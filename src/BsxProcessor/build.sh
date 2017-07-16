#! /bin/sh

PROJECT_PATH="./BsxProcessor/BsxProcessor.csproj"

dotnet restore $PROJECT_PATH
dotnet publish $PROJECT_PATH
dotnet publish $PROJECT_PATH -o ./../.build
