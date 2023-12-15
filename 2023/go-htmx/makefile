build:
	go mod download
	templ generate
	go build

clean:
	rm -f adventofcode
	@echo "rm templ generated files"
	@find . -name '*_templ.go' -delete