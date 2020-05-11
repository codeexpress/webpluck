all:
	GOOS=windows GOARCH=amd64 go build -o binaries/webpluck.exe main.go logger.go webpluck.go
	GOOS=linux   GOARCH=amd64 go build -o binaries/webpluck     main.go logger.go webpluck.go
	GOOS=darwin  GOARCH=amd64 go build -o binaries/webpluck_osx main.go logger.go webpluck.go


