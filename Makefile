OUT_PATH=./out

all: run

test:
	go test ./...

build: test
	go build -o $(OUT_PATH)

run: build
	$(OUT_PATH) data/player.txt data/rnb_calvin.txt