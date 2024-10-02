makefile_dir	:= $(abspath $(shell pwd))
m				:= "updates"

.PHONY: list bootstrap init build clean

list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

purge:
	rm -rf $(makefile_dir)/tmp

bootstrap:
	mkdir -p $$HOME/.local/bin \
		&& wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.13/tailwindcss-linux-x64 -o $$HOME/.local/bin/tailwindcss \
		&& chmod 755 $$HOME/.local/bin/tailwindcss
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest

generate:
	go generate ./...
	templ generate

dev:
	air

test: generate
	go test -v ./...

commit: test
	git add . || true
	git commit -m "$(m)" || true

push: commit
	git push origin master

pull:
	git fetch
	git pull origin master

tag:
	git tag -a $(tag) -m "$(tag)"
	git push origin $(tag)

tag-list:
	git fetch --tags
	git tag --list | sort -V

publish: generate test
	@if ack replace go.mod ;then echo 'Remove the "replace" line from the go.mod file'; exit 1; fi
	make commit m=$(m)
	make tag tag=$(m)
