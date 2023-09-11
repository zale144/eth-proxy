APP?=eth-proxy
HTTP_PORT?=8080
ETH_NETWORK_URL?=https://mainnet.infura.io/v3/b3bd9456e8d44150b963248668023317 # (or 23e0bcc507bf4f74aa685e94f8bff9dc)

clean:
	rm -f ${APP}

build: clean
	go build -o ${APP}

run: build
	HTTP_PORT=${HTTP_PORT} ETH_NETWORK_URL=${ETH_NETWORK_URL} ./${APP}

test:
	go test -v -race ./...

docker-build:
	docker build -t ${APP} .

docker-run:
	docker run --name ${APP} -p ${HTTP_PORT}:${HTTP_PORT} -e HTTP_PORT=${HTTP_PORT} -e ETH_NETWORK_URL=${ETH_NETWORK_URL} ${APP}

docker-stop:
	docker stop ${APP}
