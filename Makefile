default: test

# dev builds gatekeeper locally and places the compiled binaries into ./bin as
# well as $GOPATH/bin.
dev:
	@DATEBOOK_DEV=1 sh -c "$(CURDIR)/scripts/build.sh"


local_release: dev
	@sh -c "ln -sf $(CURDIR)/bin/datebook $(GOPATH)/bin/datebook"

# dev_run builds and runs gatekeeper locally greedily taking ports 8000 and 8001. This is just for development!
dev_run: dev
	@sh -c "$(CURDIR)/bin/datebook"
