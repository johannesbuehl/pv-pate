.PHONY: all backend client setup

all: backend client

backend:
	@echo "building server"
	cd backend; go build -ldflags "-s -w"

client:
	@echo "building client"
	cd client; npm install; npm run build

setup:
	@echo "running setup"
	cd setup; go run .