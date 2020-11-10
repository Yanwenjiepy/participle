linux64:
	GOOS=linux GOARCH=amd64 go build -o cmd/participle-linux-amd64 cmd/participle.go

windows64:
	GOOS=windows GOARCH=amd64 go build -o cmd/participle-windows-amd64.exe cmd/participle.go
