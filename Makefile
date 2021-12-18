.PHONY: test

test:
	@chmod -r ./files/file5 # Assure we can't read file5
	go test ./...
