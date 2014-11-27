KUMIS_VERSION="0.1.0"
KUMIS_KAFKA := $(shell echo $(KUMIS_KAFKA))
KUMIS_ZK := $(shell echo $(KUMIS_ZK))

dev: 
	go run src/*.go --kafka $(KUMIS_KAFKA) --zk $(KUMIS_ZK)

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

version:
	git log -n 1 --decorate --pretty=oneline > out/package/version.txt
	git log -n 1 --format="%aD" >> out/package/version.txt

release: package
	echo "About to release Kumis v" + $(KUMIS_VERSION)

	# change version number here for git tagging
	git tag $(KUMIS_VERSION)

	#redo version file
	make version

	cd out/package; zip -v -r ../kumis-$(KUMIS_VERSION).zip *
	
	# push to git first
	git push origin $(KUMIS_VERSION)
	# change version number here for posting the package to Nexus
	#mvn deploy:deploy-file -DgroupId=com.mofirouz -DartifactId=kumis -Dversion=$(KUMIS_VERSION) -Dfile=kumis.zip -DrepositoryId= -Durl=
