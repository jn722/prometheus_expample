

PROJ_DIR:=$(pwd)

export PATH:=$(PROJ_DIR)/bin:${PATH}


all: app

app:
	GOOS=linux CGO_ENABLED=0 go build -o bin/prom_test prom_test/src


restart: app
	docker-compose -f deploy/docker-compose.yml down
	docker-compose -f deploy/docker-compose.yml up -d