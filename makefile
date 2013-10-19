TESTS 			= $(notdir $(patsubst %.sh, %, $(wildcard ./test/*_test.sh) ) )

.PHONY: all
all: cryptoad

cryptoad: cryptoad.go lib.go serial-toad.go serial-lib.go
	go build -o $@

serial-lib.go: lib.go
	go-bindata -out $@ $<

serial-toad.go: assets/toad.go
	go-bindata -prefix "assets/" -out $@ $< 

.PHONY: clean
clean:
	go clean
	rm -vf serial-*.go 


define test-template
.PHONY: $(1)
$(1): test/$(1).sh cryptoad
	@test/$(1).sh > test/$(1).log 2>&1 || { echo "[-] $(1) failed:"; cat test/$(1).log; false; }
	@echo '[+] $(1) passed '
endef
$(foreach test, $(TESTS), $(eval $(call test-template,$(test)) ) )

.PHONY: tests
tests: $(TESTS)

