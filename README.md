# OpenStreetMap API Server 
![GoLang CI/CD](https://github.com/VOU-folks/osm-api/actions/workflows/golang-ci-cd.yml/badge.svg)

## Overview
This is an OpenStreetMap HTTP API Server that provides endpoints to query OpenStreetMap data.

## Features
- [x] HTTP API Server
- [x] Endpoints to query OpenStreetMap data
- [x] Download OpenStreetMap PBF files
- [x] Import OpenStreetMap PBF files

## Applications
1. [osm-api-server](api/main.go): OpenStreetMap HTTP API Server
2. [osm-pbf-downloader](scripts/osm-pbf-downloader/main.go): OpenStreetMap PBF Downloader
3. [osm-pbf-importer](scripts/osm-pbf-importer/main.go): OpenStreetMap PBF Importer

(1) **osm-api-server** is the main application for this project. It is an HTTP API server that provides endpoints to query OpenStreetMap data.

(2) **osm-pbf-downloader** is a utility application to download OpenStreetMap PBF files to storage folder.

(3) **osm-pbf-importer** is a utility application to import OpenStreetMap PBF files into supported by `osm-api-server` database.
