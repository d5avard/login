all:
	make clean
	make build 
	make run
test:
	ENV=TEST WD=$(PWD) go test -v ./...

cover:
	ENV=TEST WD=$(PWD) go test -v ./... -coverprofile=cover.txt

build:
	go build -o login

run:
	./login

clean:
	rm -f login