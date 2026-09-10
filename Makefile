help: ## показывает справку по доступным командам
	@printf "\033[33m%s:\033[0m\n" 'Commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_\/-]+:.*?## / {printf "  \033[32m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## запустить приложение 
	docker compose -f ./deploy/docker/docker-compose.yaml \
	--env-file deploy/docker/.env \
	up -d

down: ## остановить приложение 
	docker compose -f ./deploy/docker/docker-compose.yaml \
	--env-file deploy/docker/.env  \
	down

down/all:  ## остановить приложение c удалением всех ресурсов
	docker compose -f ./deploy/docker/docker-compose.yaml \
	--env-file deploy/docker/.env  \
	down \
	-v --rmi all



