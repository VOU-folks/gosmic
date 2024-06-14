build-api:
	go build -o bin/osm-api-server cmd/api/main.go
	chmod +x bin/osm-api-server

build-scripts:
	go build -o bin/osm-pbf-downloader cmd/scripts/osm-pbf-downloader/main.go
	go build -o bin/osm-pbf-importer cmd/scripts/osm-pbf-importer/main.go
	chmod +x bin/osm-pbf-downloader
	chmod +x bin/osm-pbf-importer

build: build-api build-scripts

run: build
	./bin/osm-api-server

test:
	go test -v ./...

install: build-api build-scripts
	go mod download
	go mod tidy

	mkdir -p /usr/local/osm-api/bin
	cp bin/osm-api-server /usr/local/osm-api/bin/osm-api-server
	cp bin/osm-pbf-downloader /usr/local/osm-api/bin/osm-pbf-downloader
	cp bin/osm-pbf-importer /usr/local/osm-api/bin/osm-pbf-importer

	mkdir -p /etc/osm-api
	cp config.yaml /etc/osm-api/config.yaml

	mkdir -p /var/log/osm-api

	cp osm-api-server.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable osm-api-server.service
	systemctl restart osm-api-server.service
