#!/usr/bin/make

include .env
export

.DEFAULT_GOAL := help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo "\n  Allowed for overriding next properties:\n\n\
		Usage example:\n\
	    	make run"

f=cover.out


git-push: ## full git push
	git commit -am "$(m)" && git push

build: ## build all containers (docker compose)
	docker-compose build

up: ## build & start the project (docker-compose)
	docker-compose up --build -d
	cd src/frontend && npm start

init-up: ## build & start the project (docker-compose)
	docker-compose up --build
	cd src/frontend && npm i && npm start

down: ## stop the project (docker-compose)
	docker-compose down

node-up:
	cd ./media_server && node index.js

push:
	git add docker-compose.yml .gitignore nginx/config/* nginx/Dockerfile main_backend/* media_backend/* .env* README.md && git commit -m "minor fixed" && git push
