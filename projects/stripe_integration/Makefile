.PHONY: run

# APP_NAME is used as a naming convention for resources to the local environment
APP_NAME := stripe-integration
# APP_PATH is the project/app directory in the container
APP_PATH := /${APP_NAME}

COMPOSE := docker-compose -f docker-compose.yml

# Commands for running docker compose
RUN_COMPOSE = $(COMPOSE) run --rm --service-ports -w $(APP_PATH) $(MOUNT_VOLUME) go
TEST_COMPOSE = $(COMPOSE) run --rm -w $(APP_PATH) $(MOUNT_VOLUME) go

setup: db migrate

run: MOUNT_VOLUME =  -v $(shell pwd):$(APP_PATH)
run: setup
	@$(RUN_COMPOSE) env $(shell cat .env | egrep -v '^#'| xargs) go run -mod=vendor cmd/serverd/*.go

# db runs the db service defined in the compose file
db:
	$(COMPOSE) run --rm alpine sh -c "nc -vz db 5432" && exit 0 || $(COMPOSE) up -d db
# migrate runs the db-migrate service defined in the compose file
migrate: MOUNT_VOLUME =  -v $(shell pwd)/data/migrations:/migrations
migrate:
	$(COMPOSE) run --rm $(MOUNT_VOLUME) db-migrate \
	sh -c 'sleep 5;./migrate -path /migrations -database $$DATABASE_URL up'

# down stops and removes all containers and resources associated to docker-compose.yml
down:
	$(COMPOSE) down -v

update-vendor:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor