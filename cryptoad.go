package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"bytes"
	"path/filepath"
	"io"
	"io/ioutil"
	"os/exec"
	"archive/zip"
	"errors"
	"github.com/howeyc/gopass"
	"github.com/cryptobox/gocryptobox/strongbox"
)

type data_pair struct {
	name string
	data []byte
}

func asset_archive(assets []data_pair) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, el := range assets {
		var f io.Writer
		f, err = w.Create(filepath.Join("assets", el.name))
		if err != nil {
			return
		}
		_, err = f.Write(el.data)
		if err != nil {
			return
		}
	}

	err = w.Close()
	if err != nil {
		return
	}

	return
}

func append_all(src *bytes.Buffer, dst_path string) error {
	dst, err := os.OpenFile(dst_path, os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func run_cmd(cmd *exec.Cmd) (output []byte, err error) {
	log(2, "run command: %s\n", cmd)
	output, err = cmd.CombinedOutput()
	if err != nil {
		format := "failed to run command '%s'. error = %s; output = '%s'"
		err = errors.New(fmt.Sprintf(format, cmd, err, string(output)))
	}

	return
}

func go_cmd(dst, opsys, arch string) *exec.Cmd {
	cmd := exec.Command("go", "build", "-o", dst)
	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("GOOS=%s", opsys),
		fmt.Sprintf("GOARCH=%s", arch),
	)
	
	return cmd
}

func summon_toad(dst, opsys, arch, dir, name string, salt []byte, box []byte) error {
	source_files := []data_pair{
		{"toad.go", toad_go()},
		{"lib.go", lib_go()},
	}

	assets := []data_pair{
		{"box", box},
		{"name", []byte(name)},
		{"salt", salt},
	}

	for _, el := range source_files {
		fp := filepath.Join(dir, el.name)
		log(2, "write file %s\n", fp)
		err := ioutil.WriteFile(fp, el.data, 0440)
		if err != nil {
			return err
		}
	}
	
	os.Chdir(dir)

	_, err := run_cmd(go_cmd(dst, opsys, arch))
	if err != nil {
		return err
	}

	zip_buf, err := asset_archive(assets)
	if err != nil {
		return err
	}

	log(2, "append zip to binary\n")
	if err := append_all(zip_buf, dst); err != nil {
		return err
	}

	log(2, "clean up temporary files\n")
	for _, file := range source_files {
		os.Remove(file.name)
	}
	os.Remove(dir)

	return nil
}

func check_dependencies() error {
	deps := []string{"go"}

	for _, el := range deps {
		_, err := exec.LookPath(el)
		if err != nil {
			return errors.New(el)
		}
	}

	return nil
}

func go_env(name string) (val string, err error) {
	v, err := run_cmd(exec.Command("go", "env", name)); 
	val = strings.TrimSpace(string(v))
	return
}

func go_arch() (arch string, err error) {
	arch, err = go_env("GOARCH")
	return
}

func go_opsys() (opsys string, err error) {
	opsys, err = go_env("GOOS")
	return
}

func go_root() (dir string, err error) {
	dir, err = go_env("GOROOT")
	return
}

func platform_info(f func(x string) (string, bool)) (result []string, err error) {
	goroot, err := go_root()
	if err != nil {
		err = errors.New(fmt.Sprintf("couldnt find go root: %s", err))
		return
	}

	ls, err := ioutil.ReadDir(path.Join(goroot, "bin"))
	if err != nil {
		err = errors.New(fmt.Sprintf("couldnt ls go root bin: %s", err))
		return
	}

	set := make(map[string]bool)
	for _, el := range ls {
		if el.IsDir() {
			k, ok := f(el.Name())
			if !ok {
				continue
			}

			_, found := set[k]
			if !found {
				set[k] = true
				result = append(result, k)
			}
		}
	}

	return 
}

func available_archs() (result []string, err error) {
	f := func(x string) (string, bool) {
		arr := strings.Split(x, "_")
		if len(arr) > 1 {
			return arr[1], true
		}
		return "", false
	}

	return platform_info(f)
}

func available_oses() (result []string, err error) {
	f := func(x string) (string, bool) {
		return strings.Split(x, "_")[0], true
	}

	return platform_info(f)
}

func get_passwd() (pass []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to get password")
		}
	}()

	i := 0
	for {
		fmt.Printf("enter password: ")
		pass = gopass.GetPasswd()
		fmt.Printf("re-enter password: ")
		confirm := gopass.GetPasswd()
		if bytes.Equal(pass, confirm) {
			break
		}

		if i < 2 {
			fmt.Printf("passwords did not match. try again\n")
		} else {
			err = errors.New("too many password failures")
			return
		}
		i += 1
	}

	return
}

func run(verbosity int, input, output, opsys, arch, pass string) {
	var pw []byte

	log_level(verbosity)

	abs_dst, err := filepath.Abs(output)
	if err != nil {
		err_exit("failed to get absolute file path for %s: %s", output, err)
	}

	log(2, "reading input file\n")
	msg_data, err := ioutil.ReadFile(input)
	if err != nil {
		err_exit("error opening input file: '%s'\n", err)
	}

	if len(pass) < 1 {
		pw, err = get_passwd()
		if err != nil {
			err_exit(err.Error())
		}
	} else {
		pw = []byte(pass)
	}

	if len(pw) < 8 {
		err_exit("password is really too short. you can do better than that.")
	} else if len(pw) < 10 {
		log(0, "warning: password is weak\n")
	}

	log(1, "generate key from password\n")
	salt, key, err := gen_key(pw)
	if err != nil {
		err_exit("failed to generate key")
	}
	//TODO: maybe erase password from memory here?
	
	log(1, "seal message\n")
	box, ok := strongbox.Seal(msg_data, key)
	if !ok {
		err_exit("failed to seal message")
	}

	log(1, "message sealed! size = %d bytes\n", len(box))
	dir, err := ioutil.TempDir("", "cryptoad-summon")
	if err != nil {
		err_exit("failed to create tmp dir: %s", err)
	}

	log(1, "summon toad...\n")
	log(2, "using tmp dir %s\n", dir)
	if err := summon_toad(abs_dst, opsys, arch, dir, path.Base(input), salt, box); err != nil {
		err_exit("%s", err)
	}
	log(1, "finished.\n")
}

func main() {
	const pass_desc      = "       the password with which to encrypt. if not specified, you\n" + 
	                       "                   " + 
	                       "will be prompted for one. make it a good one " +
	                       "whydontcha!?"
	const verbosity_desc = "           verbosity level. higher level is more verbose"
	const arch_desc      = "  target architecture. one of %s"
	const opsys_desc     = "    target OS. one of %s"
	const usage_fmt      = "usage: %s INPUT OUTPUT\n" + 
	                       "  INPUT:           the file to encrypt\n" + 
	                       "  OUTPUT:          your newly summoned, self-decrypting toad friend /.0 _0}\n"
	var pass string
	var verbosity int

	if err := check_dependencies(); err != nil {
		err_exit("command '%s' required to be in executable PATH", err)
	}

	arches, err := available_archs()
	if err != nil {
		err_exit("failed to find list of architectures: %s", err)
	}

	oses, err := available_oses()
	if err != nil {
		err_exit("failed to find list of OSes: %s", err)
	}

	arch, err := go_arch()
	if err != nil {
		err_exit("failed to get default go arch: %s", err)
	}

	opsys, err := go_opsys()
	if err != nil {
		err_exit("failed to get default go os: %s", err)
	}

	flag.StringVar(&pass, "pass", "", pass_desc)
	flag.IntVar(&verbosity, "v", 0, verbosity_desc)
	flag.StringVar(&arch, "arch", arch, fmt.Sprintf(arch_desc, strings.Join(arches, ", ")))
	flag.StringVar(&opsys, "os", opsys, fmt.Sprintf(opsys_desc, strings.Join(oses, ", ")))
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage_fmt, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) != 2 {
		log(0, "two positional arguments expected. got %s\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	run(verbosity, flag.Args()[0], flag.Args()[1], 
			opsys, arch, pass)
}