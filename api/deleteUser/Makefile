run:
	@echo "\n*** Run locally ***\n"
	$(MAKE) build
	@echo "Starting"
	sam local start-api -p 1234

deployDev:
	@echo "\n*** Building for DEV ***\n"
	$(MAKE) build
	$(MAKE) uploadDev

uploadDev:
	@echo "Uploading new version to Dev..."
	aws lambda update-function-code \
		--region eu-north-1 \
		--function-name back-test-delete-user \
		--zip-file fileb://aws/bootstrap.zip

build:
	@echo "Building ..."
	[ -d aws ] || mkdir -p aws
	GOARCH=amd64 GOOS=linux go build -o aws/bootstrap main.go
	zip -j aws/bootstrap.zip aws/bootstrap
	@echo "Build OK!"
