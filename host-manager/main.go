package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const path = "/etc/hosts"

// examples:
// - sudo go run host-manager/main.go -d super.loc -add
// - sudo go run host-manager/main.go -d super.loc -rm
func main() {
	isAdd := flag.Bool("add", false, "add")
	isRm := flag.Bool("rm", false, "remove")
	d := flag.String("d", "", "domain")
	flag.Parse()

	if !*isAdd && !*isRm {
		log.Fatal(fmt.Sprintf("need specify add|rm"))
	}

	if *d == "" {
		log.Fatal(fmt.Sprintf("domain is not set"))
	}
	domain := *d

	f := openFile()
	defer f.Close()

	if *isAdd && !exists(domain) {
		add(f, domain)
	}

	if *isRm && exists(domain) {
		remove(domain)
	}

	fmt.Println("OK!")
}

func exists(domain string) bool {
	data := read()
	m, err := regexp.Match(buildStr(domain), data)
	if err != nil {
		log.Fatal(fmt.Sprintf("match err: %v", err))
	}

	return m
}

func add(f *os.File, domain string) {
	str := buildStr(domain)

	_, err := f.WriteString(str)
	if err != nil {
		log.Fatal(fmt.Sprintf("write file err: %v", err))
	}
}

func remove(domain string) {
	data := read()
	rg := regexp.MustCompile(buildStr(domain))
	res := rg.ReplaceAll(data, []byte(""))

	err := ioutil.WriteFile(path, res, fileInfo().Mode())
	if err != nil {
		log.Fatal(fmt.Sprintf("remove from file err: %v", err))
	}
}

func read() []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("read file err: %v", err))
	}

	return data
}

func buildStr(domain string) string {
	return "127.0.0.1   " + domain + "\n"
}

func fileInfo() os.FileInfo {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("fileinfo err: %v", err))
	}

	return fi
}

func openFile() *os.File {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, fileInfo().Mode())
	if err != nil {
		log.Fatal(fmt.Sprintf("open file err: %v", err))
	}

	return f
}
