test:
	go test -coverprofile=coverage.out ./... -cover -bench=
	go tool cover -html=coverage.out -o coverage.html
