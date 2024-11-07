.PHONY: all backend client setup init

all: backend client

out_dir = dist

backend:
	@echo "building server"
	cd backend; go build -ldflags "-s -w" -o ../$(out_dir)/backend/

client:
	@echo "building client"
	cd client; npm install; npm run build

init:
	@echo "creating \"backend/config.yaml\""
	@cp -n backend/example-config.yaml backend/config.yaml

setup:
	@echo "running setup"
	cd setup; go run .