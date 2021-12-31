.PHONY: build clean run test image ctn-run

TKN ?= token
SRC_DIR ?= ./src/
DST_DIR ?= ./dst/

build:
	@CGO_ENABLED=0 go build -o ./out/eq .

clean:
	@rm -r ./out/

run: build
	@./out/./eq

test:
	@chmod -r ./files/file5 # Assure we can't read file5
	@TKN=$(TKN) SRC_DIR=$(SRC_DIR) DST_DIR=$(DST_DIR) go test ./...
	@chmod +r ./files/file5

image: build
	docker build --no-cache -t eq .

ctn-run: image
	@docker run \
		-d \
		--volume $(SRC_DIR):/srcdir/ \
		--volume $(DST_DIR):/dstdir/ \
		--env TKN=$(TKN) \
		--env SRC_DIR=/srcdir/ \
		--env DST_DIR=/dstdir/ \
		eq
