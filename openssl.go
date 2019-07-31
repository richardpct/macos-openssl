// openssl package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"log"
	"os"
	"os/exec"
	"path"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "openssl"
	vers     = "1.0.2s"
	ext      = "tar.gz"
	url      = "https://www.openssl.org/source"
	hashType = "sha256"
	hash     = "cabd5c9492825ce5bd23f3c3aeed6a97f8142f606d893df216411f07d1abab96"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while configuring ...")
	cmd := exec.Command("/usr/bin/perl",
		"./Configure",
		"darwin64-x86_64-cc",
		"--prefix="+*destdir,
		"--openssldir="+*destdir+"/etc/openssl")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func build() {
	fmt.Println("Waiting while compiling ...")
	cmd := exec.Command("make", "depend")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}

	cmd = exec.Command("make", "-j" + pkgsrc.Ncpu)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func install() {
	fmt.Println("Waiting while installing ...")
	cmd := exec.Command("make", "install")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
	install()
}
