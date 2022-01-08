.PHONY: swagger deploy
# Remember to adjust env file and main.go before make
swagger:
	swag init --parseDependency --parseInternal
build:
	DOCKER_BUILDKIT=1 docker build -t dms-be .
deploy: build
	docker save dms-be:latest > ~/Desktop/dms-be.tar && \
	sshpass -p "Iot@@123" scp /home/nguyen/Desktop/dms-be.tar  sviot@iot.hcmue.space:~/iot/