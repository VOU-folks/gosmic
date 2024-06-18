build-api:
	go build -o bin/gosmic-api-server cmd/api/main.go
	chmod +x bin/gosmic-api-server

build-scripts:
	go build -o bin/gosmic-pbf-downloader cmd/scripts/gosmic-pbf-downloader/main.go
	go build -o bin/gosmic-pbf-importer cmd/scripts/gosmic-pbf-importer/main.go
	chmod +x bin/gosmic-pbf-downloader
	chmod +x bin/gosmic-pbf-importer

build: build-api build-scripts

run: build
	./bin/gosmic-server

test:
	go test -v ./...

create-linux-user:
	useradd -d /etc/gosmic -s /usr/sbin/nologin gosmic

install: build-api build-scripts
	go mod download
	go mod tidy

	mkdir -p /usr/local/gosmic/bin
	cp bin/gosmic-api-server /usr/local/gosmic/bin/gosmic-api-server
	cp bin/gosmic-pbf-downloader /usr/local/gosmic/bin/gosmic-pbf-downloader
	cp bin/gosmic-pbf-importer /usr/local/gosmic/bin/gosmic-pbf-importer

	mkdir -p /etc/gosmic
	cp config.yaml /etc/gosmic/config.yaml

	mkdir -p /var/log/gosmic

	cp gosmic-api-server.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable gosmic-api-server.service
	systemctl restart gosmic-api-server.service

configure-nginx:
	cp certs/gosmic.io.crt /etc/gosmic/
	cp certs/gosmic.io.key /etc/gosmic/
	cp nginx-vhost.conf /etc/gosmic/nginx-vhost.conf
	rm -f /etc/nginx/sites-enabled/gosmic
	ln -s /etc/gosmic/nginx-vhost.conf /etc/nginx/sites-enabled/gosmic
	chown -R gosmic:gosmic /etc/gosmic