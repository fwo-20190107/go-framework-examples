
PHONY: demo
demo:
	docker compose -f ./docker/compose.yaml up

PHONY: stop
stop:
	docker compose -f ./docker/compose.yaml stop

PHONY: clean
clean:
	docker container rm db-mysql app

PHONY: develop
develop:
	docker compose -f ./docker/compose.yaml run --service-ports --name dev-mysql -d mysql
	air