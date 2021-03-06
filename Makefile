ARCH?=amd64
VERSION?=v1.1.0
PREFIX?=1314520999


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) GO111MODULE=on go build -a -o faucet cmd/main.go
docker-container:
	docker build -t $(PREFIX)/faucet:$(VERSION) .
test:
	go run cmd/main.go  --label="app.kubernetes.io/name=bee" --container=bee --webhook="https://bin.webhookrelay.com/v1/webhooks/"  --command="curl -s localhost:1635/addresses"
clean:
	rm -f faucet
