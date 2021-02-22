GO := go

all:

test:
	cd cmap && go test

clean:
	find . -name '*~' -delete

.PHONY: all clean test
