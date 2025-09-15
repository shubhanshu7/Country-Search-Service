
test:
	go test ./... -count=1

cover:
	go test ./... -coverprofile=coverage.out

race:
	go test ./... -race -count=1
