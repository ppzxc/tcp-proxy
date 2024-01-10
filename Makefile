restart: build down up log

restart-envoy:
	docker compose down -v envoy
	docker compose up --build -d envoy
	docker compose logs -f 

build:
	go build -C ./echo ./
	go build -C ./client -o ./

up:
	docker compose up --build -d

down:
	docker compose down -v

log:
	docker compose logs -f
