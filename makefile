TESTS 			= $(notdir $(patsubst %.sh, %, $(wildcard ./test/*_test.sh) ) )

.PHONY: all
all: cryptoad

cryptoad: cryptoad.go lib.go serial-toad.go serial-lib.go
	go build -o $@

define asset-template
serial-$(1).go: $(2)
	cat assets/template.hdr | sed 's/TEMPLATENAME/$(1)/' > $$@.tmp
	cat $$< >> $$@.tmp
	cat assets/template.ftr >> $$@.tmp
	mv $$@.tmp $$@
endef

$(eval $(call asset-template,lib,lib.go))
$(eval $(call asset-template,toad,assets/toad.go))

.PHONY: clean
clean:
	go clean
	rm -vf serial-*.go 
	rm -vf serial-*.go.tmp

define test-template
.PHONY: $(1)
$(1): test/$(1).sh cryptoad
	@test/$(1).sh > test/$(1).log 2>&1 || { echo "[-] $(1) failed:"; cat test/$(1).log; false; }
	@echo '[+] $(1) passed '
endef
$(foreach test, $(TESTS), $(eval $(call test-template,$(test)) ) )

.PHONY: tests
tests: $(TESTS)

