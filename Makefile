build-executable:
	go build -v -o youtubeSubtitles

build-image:
	docker build -t youtube-subs .

run-docker-compose: build-image
	docker-compose up
