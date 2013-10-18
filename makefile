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
	@echo '######################################################################'; done
	@echo '############# TEST $(1) '
	@echo '######################################################################'; done
	@test/$(1).sh || { echo "test failed"; false; }
	@echo
endef
$(foreach test, $(TESTS), $(eval $(call test-template,$(test)) ) )

.PHONY: tests
tests: $(TESTS)

