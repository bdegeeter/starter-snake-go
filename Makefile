
.PHONY: play-local
play-local:
	battlesnake play -W 11 -H 11 --name first-post --url http://localhost:8080 -g solo --viewmap

.PHONY: serve-local
serve-local:
	go run main.go

.PHONY: vendor
vendor:
	go mod vendor
