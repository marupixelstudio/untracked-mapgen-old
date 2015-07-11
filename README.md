### Install

go
python
pip
postgresql
postgis
node


### Prep

createdb untracked
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

### Run

go run main.go

### Generate tiles

export MAPNIK_MAP_FILE=styles/layers-contours.xml
export MAPNIK_TILE_DIR=tiles