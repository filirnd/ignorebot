build:
	go build
build-arm:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ignorebot-arm
	mkdir -p deploy/arm
	mv ignorebot-arm deploy/arm/ignorebot
	cp config.yaml deploy/arm
	cp start.sh deploy/arm
	chmod +x deploy/arm/start.sh

