build-executable:
	go build -v -o youtubeSubtitles

build-image: build-executable
	docker build -t youtube-subs .

run-docker-compose: build-image
	docker-compose up
