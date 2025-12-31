all: run

build:
	docker build -t cryptoimage .

run: build
	docker run -d -p 8080:8080 --name cryptoserver cryptoimage

clean:
	docker stop cryptoserver
	docker rm cryptoserver
