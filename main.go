package main

import (
	dmca "github.com/ipfs/refs/dmca"
	shell "github.com/whyrusleeping/ipfs-shell"

	"flag"
	"fmt"
	"log"
)

func addSkeleton(sh *shell.Shell, dmca string) (string, error) {
	hdenylists, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}
	hdenylists, err = sh.PatchLink(hdenylists, "dmca", dmca, true)
	if err != nil {
		return "", err
	}
	hlists, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}
	hlists, err = sh.PatchLink(hlists, "denylists", hdenylists, true)
	if err != nil {
		return "", err
	}
	hroot, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}
	hroot, err = sh.PatchLink(hroot, "lists", hlists, true)
	if err != nil {
		return "", err
	}
	return hroot, nil
}

func addPrevious(sh *shell.Shell, h string, c string) (string, error) {
	cur, err := sh.ResolvePath(c)
	if err != nil {
		log.Fatalf("could not resolve current version: %s", err)
	}
	prev, err := sh.ResolvePath(c + "/previous")
	if err == nil {
		h, err = sh.PatchLink(h, "previous", prev, true)
		if err != nil {
			return "", err
		}
	}

	if h != cur {
		h, err = sh.PatchLink(h, "previous", cur, true)
		if err != nil {
			return "", err
		}
	}

	return h, nil
}

func main() {
	u := flag.String("uri", "127.0.0.1:5001", "the IPFS API endpoint to use")
	c := flag.String("current", "/ipns/refs.ipfs.io", "add links to previous versions")
	flag.Parse()

	sh := shell.NewShell(*u)

	h, err := dmca.AddDenylist(sh)
	if err != nil {
		log.Fatalf("dmca.AddDenylist: %s", err)
	}

	if len(*c) > 0 {
		p, err := addPrevious(sh, h, *c+"/lists/denylists/dmca")
		if err != nil {
			log.Printf("could not add previous link: %s\n", err)
		} else {
			h = p
		}
	}

	h, err = addSkeleton(sh, h)
	if err != nil {
		log.Fatalf("main.addSkeleton: %s", err)
	}

	fmt.Println(h)
}
