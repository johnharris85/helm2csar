.PHONY: build assets all gen

all: assets build

build:
	go build -v -o ./h2c .

assets:
	cd cmd && rice embed-go

gen: all
	./h2c generate examples/charts/fluent-bit
