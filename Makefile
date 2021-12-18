.PHONY: build clean run test

build:
	go build -o ./out/eq .

clean:
	rm -r ./out/

run: build
	./out/./eq

test:
	@chmod -r ./files/file5 # Assure we can't read file5
	go test ./...
