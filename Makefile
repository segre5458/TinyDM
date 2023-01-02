DISTFILE=emptty

# test:
# 	@echo "Testing..."
# 	@go test -coverprofile cover.out ${TAGS_ARGS} ./...
# 	@echo "Done"

clean:
	@echo "Cleaning..."
	@rm -f dist/${DISTFILE}
	@rm -f dist/emptty.1.gz
	@rm -rf dist
	@echo "Done"

build:
	@echo "Building${TAGS_ARGS}..."
	@mkdir -p dist
	# @go build ${TAGS_ARGS} -o dist/${DISTFILE} -ldflags "-X github.com/tvrzna/emptty/src.buildVersion=${BUILD_VERSION}" ${GOVCS}
	# @gzip -cn res/emptty.1 > dist/emptty.1.gz
	@echo "Done"

install:
	@echo "Installing..."
	@install -DZs dist/${DISTFILE} -m 755 -t ${DESTDIR}/usr/bin
	@echo "Done"

install-manual:
	@echo "Installing manual..."
	@install -D dist/emptty.1.gz -t ${DESTDIR}/usr/share/man/man1
	@echo "Done"

install-pam:
	@echo "Installing pam file..."
	@install -DZ res/pam -m 644 -T ${DESTDIR}/etc/pam.d/${DISTFILE}
	@echo "Done"

install-all: install install-manual install-pam