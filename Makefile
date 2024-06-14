build:
	go build -o bin/api cmd/api/main.go

run: build
	./bin/api

install:
	go mod download
	go mod tidy
	mkdir -p /etc/osm-api
	cp config.yaml /etc/osm-api/config.yaml
	mkdir -p /var/log/osm-api
	go build -o /usr/local/bin/osm-api cmd/api/main.go
	cp osm-api.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable osm-api.service
	systemctl restart osm-api.service
