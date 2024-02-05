#
# There are fewer things I hate more than gnu make
#

# where to find my targets
CMD_DIR := ./cmd
# where to put my targets
BIN_DIR := ./bin
TARGETS := echo uniques
export SHELL := /bin/bash # I use fish in macos and zsh on linux

.PHONY: build clean

build:
	@for target in ${TARGETS}; do                               \
	 go build -v -o ${BIN_DIR}/$${target} ${CMD_DIR}/$${target}; \
    done

clean:
	@rm -vrf ${BIN_DIR}
