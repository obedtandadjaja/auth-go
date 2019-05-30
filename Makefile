SECRET_KEY       ?= "wow-much-secret-very-discreet"
ENV              ?= "development"
DB_HOST          ?= "localhost"
DB_PORT          ?= "5432"
DB_USER          ?= "obedt"
DB_PASSWORD      ?= ""
DB_NAME          ?= "auth"
TEST_DB_USER     ?= "obedt"
TEST_DB_PASSWORD ?= ""
TEST_DB_NAME     ?= "auth-test"

run:
	SECRET_KEY=$(SECRET_KEY) \
	ENV=$(ENV) \
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
  TEST_DB_USER=$(TEST_DB_USER) \
  TEST_DB_PASSWORD=$(TEST_DB_PASSWORD) \
  TEST_DB_NAME=$(TEST_DB_NAME) \
  go build && ./auth-go
