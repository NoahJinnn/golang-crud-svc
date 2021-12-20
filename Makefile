.PHONY: swagger deploy
swagger:
	swag init --parseDependency --parseInternal
deploy:
	DOCKER_BUILDKIT=1 docker build -t dms-be . && \
	docker save dms-be:latest > ~/Desktop/dms-be.tar && \
	sshpass -p "Iot@@123" scp /home/nguyen/Desktop/dms-be.tar  sviot@iot.hcmue.space:~/iot/