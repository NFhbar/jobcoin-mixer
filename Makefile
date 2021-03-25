HOME:= $(shell echo $$HOME)

.PHONY: install
install:
	@go install