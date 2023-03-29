
PHONY: start
start:
	docker compose -f ./docker/compose.yaml up

PHONY: stop
stop:
	docker compose -f ./docker/compose.yaml stop

PHONY: clean
clean:
	docker container rm db-mysql app