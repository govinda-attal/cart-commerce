init:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega/...

install: init
	rm -rf ./vendor
	dep ensure

test: install
	ginkgo -r

build: 
	rm -rf ./dist
	mkdir dist
	mkdir dist/config
	mkdir dist/rules
	mkdir dist/api
	GOOS=linux GOARCH=amd64 go build -o ./dist/cartcom .
	cp ./test/fixtures/app-config.yml ./dist/config/app-config.yml
	cp ./api/* ./dist/api/
	cp ./rules/* ./dist/rules/

serve:
	cd dist && ./cartcom

docker-compose-start:
	docker-compose up #-d db adminer swagger-ui

clean:
	rm ./dist/ -rf

pack:
	docker build -t gattal/cart-commerce:latest .

upload:
	docker push gattal/cart-commerce:latest	

ship: init test pack upload clean