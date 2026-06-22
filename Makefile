OUT_PATH=./bin/myprog

all: run

test:
	go test ./...

build: test
	go build -o $(OUT_PATH)

run: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt -v 

brief: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt

rain: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt -v -w 1

sun: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt -v -w 2

sandstorm: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt -v -w 3

hail: build
	$(OUT_PATH) data/player.txt data/rnb_trainer_1.txt -v -w 4