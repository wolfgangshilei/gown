.PHONY: release-ios

RELEASE_VERSION=0.1.0

release-ios:
	mkdir -p releases
	gomobile bind -target=ios .
	tar czvf release-wordnet-go-ios-${RELEASE_VERSION}.tgz Wordnet.framework
	rm -rf Wordnet.framework
	mv release-wordnet-go-ios-${RELEASE_VERSION}.tgz releases/
