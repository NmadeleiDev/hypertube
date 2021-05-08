#!/usr/bin/make

include .env
export

.DEFAULT_GOAL := help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo "\n  Allowed for overriding next properties:\n\n\
		Usage example:\n\
	    	make all"

f=cover.out

git-push: ## full git push
	git commit -am "$(m)" && git push

build: ## build all containers (docker compose)
	docker-compose build

#up: build-front up-back
up: up-back

re: down up

clean:
	rm -f files/*
	rm -rf data/*

up-back:
	docker-compose up --build -d

up-front: ## build & start the project (docker-compose)
	cd src/frontend && npm start

init-front:
	cd src/frontend && npm i

build-front: init-front
	cd src/frontend && npm run build && cp -R build/* ../nginx/static/

init-up: up-back init-front up-front

down: ## stop the project (docker-compose)
	docker-compose down

push:
	git add docker-compose.yml .gitignore nginx/config/* nginx/Dockerfile main_backend/* media_backend/* .env* README.md && git commit -m "minor fixed" && git push

all: build-front up

.DEFAULT_GOAL := all
