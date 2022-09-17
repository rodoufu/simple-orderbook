test:
	go test -v --cover ./...

build:
	cd cmd/simplebook && go build

run:
	./cmd/simplebook/simplebook

clean:
	rm cmd/simplebook/simplebook || true

build_docker:
	docker build -t github.com/rodoufu/simple-orderbook:latest .

run_docker:
	docker run --rm --name simplebook -v $(pwd)/input_file.csv:/app/input_file.csv github.com/rodoufu/simple-orderbook:latest input_file.csv