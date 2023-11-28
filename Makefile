.PHONY:

service := link-fixer

build:
	docker-compose build --no-cache ${service}

start:
	docker-compose --compatibility up --build -d ${service}

run:
	docker-compose --compatibility up --build ${service}

clean:
	docker-compose down --volumes --remove-orphans
