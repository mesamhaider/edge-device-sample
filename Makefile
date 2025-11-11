run-sim:
	./sim/device-simulator-mac-amd64 --port $(PORT)

run-server:
	go run ./cmd

run-docker:
	docker compose up

stop-docker:
	docker compose down --volumes --rmi all

build-docker:
	docker compose build

build-and-run-docker:
	docker compose build
	docker compose up -d