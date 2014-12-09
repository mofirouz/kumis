KUMIS_VERSION="0.1.1"

dev: 
	go run src/*.go

all: build package

deps:
	go get github.com/Shopify/sarama
	go get github.com/go-martini/martini
	go get github.com/samuel/go-zookeeper/zk

compile:
	go build -v -o out/kumis src/*.go

build:
	make deps
	make compile

clean:
	# remove to go get fresh new ones in every build
	go clean
	rm -rf out

distclean:
	make clean
	rm -rf out

package: build
	mkdir -p out/package/
	mv out/kumis out/package/kumis
	cp -r static out/package/static
	make version
	cd out/package; zip -v -r ../kumis-$(KUMIS_VERSION).zip *

version:
	git log -n 1 --decorate --pretty=oneline > out/package/version.txt
	git log -n 1 --format="%aD" >> out/package/version.txt

release: package
	echo "About to release Kumis v" + $(KUMIS_VERSION)
	git tag $(KUMIS_VERSION)
	git push origin $(KUMIS_VERSION)

	make version
	rm out/kumis-$(KUMIS_VERSION).zip
	cd out/package; zip -v -r ../kumis-$(KUMIS_VERSION).zip *
