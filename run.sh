#!/bin/bash

echo "Starting Holger Hahn Website..."
echo "Building templates..."
templ generate

echo "Starting server on http://localhost:8080"
go run main.go