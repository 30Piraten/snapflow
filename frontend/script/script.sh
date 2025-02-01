#!/bin/bash

# Navigate to the directory containing lambda.go (Important!)
cd ../lambda

# Compile Lambda function and name the output "bootstrap"
GOOS=linux GOARCH=amd64 go build -o bootstrap lambda.go

# Create the zip file including the "bootstrap" binary
zip dummyprinter.zip bootstrap

# Verify the zip file contents
zipinfo -e dummyprinter.zip

# Make "bootstrap" executable (Important!)
chmod +x bootstrap

# Verify the zip file contents and permissions (crucial!)
# zipinfo dummyprinter.zip

echo "Lambda deployment package (dummyprinter.zip) created."
