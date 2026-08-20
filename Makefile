LDFLAGS="-X main.Buildstamp=`date '+%Y-%m-%d_%I:%M:%S%p'` -X main.Version=`git describe --tags --always` -X main.Githash=`git rev-parse --short HEAD` -s -w"

build: clean
	go mod tidy
	go build -ldflags $(LDFLAGS) -o ./cfddns main.go


# install:
# 	mkdir -vp /usr/local/bin/
# 	cp output/gobeat /usr/local/bin/


clean:
	rm -rf ./cfddns