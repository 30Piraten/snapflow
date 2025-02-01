# #!/bin/bash

# Navigate to the directory containing lambda.go (Important!)
cd ../lambda

# Clean up
rm -f bootstrap dummyprinter.zip

# Compile Lambda function and name the output "bootstrap"
GOOS=linux GOARCH=amd64 go CGO_ENABLED=0 build -o bootstrap lambda.go

# Make "bootstrap" executable (Important!)
chmod +x bootstrap

# Create the zip file including the "bootstrap" binary
zip -X dummyprinter.zip bootstrap

echo "Lambda deployment package (dummyprinter.zip) created."