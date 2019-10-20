SECRET_KEY       ?= "wow-much-secret-very-discreet"
ENV              ?= "development"
DB_HOST          ?= "localhost"
DB_PORT          ?= "5432"
DB_USER          ?= "obedtandadjaja"
DB_PASSWORD      ?= ""
DB_NAME          ?= "auth"
TEST_DB_USER     ?= "obedtandadjaja"
TEST_DB_PASSWORD ?= ""
TEST_DB_NAME     ?= "auth_test"
APP_HOST         ?= "localhost"
APP_PORT         ?= "8080"

run:
	export SECRET_KEY=$(SECRET_KEY) \
         ENV=$(ENV) \
	       DB_HOST=$(DB_HOST) \
	       DB_PORT=$(DB_PORT) \
	       DB_USER=$(DB_USER) \
	       DB_PASSWORD=$(DB_PASSWORD) \
	       DB_NAME=$(DB_NAME) \
         APP_HOST=$(APP_HOST) \
         APP_PORT=$(APP_PORT); \
  go clean; \
  go build; \
  ./auth-go

test:
	export SECRET_KEY=$(SECRET_KEY) \
         ENV=$(ENV) \
	       DB_HOST=$(DB_HOST) \
	       DB_PORT=$(DB_PORT) \
         TEST_DB_USER=$(TEST_DB_USER) \
         TEST_DB_PASSWORD=$(TEST_DB_PASSWORD) \
         TEST_DB_NAME=$(TEST_DB_NAME) \
         APP_HOST=$(APP_HOST) \
         APP_PORT=$(APP_PORT); \
  go clean; \
  go build; \
  go test
