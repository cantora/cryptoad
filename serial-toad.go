package main

func toad_go() []byte {
	return []byte(`
package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"github.com/howeyc/gopass"
	"github.com/cryptobox/gocryptobox/strongbox"
)

func asset(name string) []byte {
	value, err := get_asset(name)
	if err != nil {
		err_exit("failed to extract '%s' asset: %s", name, err)
	}

	return value
}

func main() {
	const pass_desc        = "the password with which to decrypt. " +
	                         "if not specified password will be prompted"
	const verbosity_desc   = "the level of verbosity. higher is more verbose"
	var out string
	var verbosity int
	var pass string
	var pw []byte

	default_name := string(asset("name"))

	flag.StringVar(&out, "out", default_name, "output path")
	flag.StringVar(&pass, "pass", "", pass_desc)
	flag.IntVar(&verbosity, "v", 0, verbosity_desc)
	flag.Parse()
	log_level(verbosity)

	if len(pass) < 1 {
		fmt.Printf("enter password: ")
		pw = gopass.GetPasswd()
	} else {
		pw = []byte(pass)
	}

	key := get_key(pw, asset("salt"))
	msg, ok := strongbox.Open(asset("box"), key)
	if !ok {
		err_exit("wrong password. failed to decrypt box")
	}

	if err := ioutil.WriteFile(out, msg, 0440); err != nil {
		err_exit("failed to write output file: %s", err)
	}

}
`)
}
