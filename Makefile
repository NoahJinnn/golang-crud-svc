.PHONY: swagger deploy
# Remember to adjust env file and main.go before make
APP_NAME=dms-be
VERSION=latest
IMAGE_NAME=$(APP_NAME):$(VERSION)
LOCAL_DIR=/home/nguyen/Desktop/$(APP_NAME).tar
SERVER_PW=Iot@@123
SERVER_USER=sviot
SERVER_HOST=iot.hcmue.space
SERVER_DIR=~/iot/

swagger:
	swag init --parseDependency --parseInternal
build:
	DOCKER_BUILDKIT=1 docker build -t dms-be .
deploy: build
	docker save $(IMAGE_NAME) > $(LOCAL_DIR) && \
	sshpass -p "$(SERVER_PW)" scp $(LOCAL_DIR)  $(SERVER_USER)@$(SERVER_HOST):$(SERVER_DIR)