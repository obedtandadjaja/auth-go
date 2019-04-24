SECRET_KEY  ?= "wow-much-secret-very-discreet"
ENV         ?= "development"
DB_HOST     ?= "localhost"
DB_PORT     ?= "5432"
DB_USER     ?= "obedt"
DB_PASSWORD ?= ""
DB_NAME     ?= "auth"

run:
	SECRET_KEY=$(SECRET_KEY) \
	ENV=$(ENV) \
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
	go run *.go
