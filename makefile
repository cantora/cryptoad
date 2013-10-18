.PHONY: all
all: cryptoad

cryptoad: cryptoad.go lib.go serial-toad.go serial-lib.go
	go build

serial-lib.go: lib.go
	go-bindata -out $@ $<

serial-toad.go: assets/toad.go
	go-bindata -prefix "assets/" -out $@ $< 

.PHONY: clean
clean:
	go clean
	rm -vf serial-*.go 

