OUT_PATH=./bin/myprog

all: run

test:
	go test ./...

build: test
	go build -o $(OUT_PATH)

run: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt

verbose: build
	$(OUT_PATH) -v data/player.txt data/rnb_trainer_1.txt