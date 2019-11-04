FROM golang:latest

LABEL maintainer="Obed Tandadjaja <obed.tandadjaja@gmail.com>"

ENV APP_HOME /auth-go
RUN mkdir $APP_HOME
WORKDIR $APP_HOME
