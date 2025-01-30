# Dummy Printer
cd frontend/config
GOOS=linux GOARCH=amd64 go build main -o lambda.go
zip dummyPrinter.zip bootstrap
chmod +x bootstrap
cd -

