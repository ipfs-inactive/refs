package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Notice struct {
	Uri  string   `json:"uri"`
	Keys []string `json:"keys"`
}

var repo = "/home/lars/workspace/ipfs/dmca"
var repoDest = "./dmca"
var denylist []Notice

func fetchDenylist() ([]Notice, error) {
	exists := true
	_, err := os.Stat(repoDest)
	if err != nil && os.IsNotExist(err) {
		exists = false
	} else if err != nil {
		return nil, err
	}

	if exists {
		runGit(&exec.Cmd{
			Path: "/usr/bin/git",
			Args: []string{"/usr/bin/git", "fetch", "-av", "--progress"},
			Dir:  repoDest,
		})
		runGit(&exec.Cmd{
			Path: "/usr/bin/git",
			Args: []string{"/usr/bin/git", "reset", "--hard", "origin/master"},
			Dir:  repoDest,
		})
	} else {
		runGit(&exec.Cmd{
			Path: "/usr/bin/git",
			Args: []string{"/usr/bin/git", "clone", "-v", "--progress", repo, repoDest},
		})
	}

	// TODO: make use of ioutil from here on

	repoDir, err := os.Open(repoDest)
	if err != nil {
		return nil, err
	}

	dirs, err := repoDir.Readdir(0)
	if err != nil {
		return nil, err
	}

	dnylist := []Notice{}

	for _, dir := range dirs {
		dirName := strings.Join([]string{repoDir.Name(), dir.Name()}, "/")
		keysName := strings.Join([]string{dirName, "keys"}, "/")
		noticeName := strings.Join([]string{dirName, "notice.md"}, "/")

		_, err := os.Stat(keysName)
		_, err2 := os.Stat(noticeName)
		if err != nil || err2 != nil {
			log.Printf("fetch: skip %s", dirName)
			continue
		}

		notice := Notice{
			Uri:  fmt.Sprintf("http://dmca.ipfs.io/%s", dir.Name()),
			Keys: []string{},
		}

		b, err := ioutil.ReadFile(keysName)
		if err != nil {
			log.Printf("fetch: %s read error: %s", keysName, err)
			continue
		}
		scan := bufio.NewScanner(strings.NewReader(string(b)))
		for scan.Scan() {
			notice.Keys = append(notice.Keys, scan.Text())
		}

		dnylist = append(dnylist, notice)
	}

	return dnylist, nil
}

func runGit(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		log.Printf("fetch git: %s\n", scan.Text())
	}

	cmd.Wait()
	return nil
}

func main() {
	go func() {
		dl, err := fetchDenylist()
		if err != nil {
			log.Printf("fetch error: %s\n", err)
		} else {
			numKeys := 0
			for _, notice := range dl {
				numKeys = numKeys + len(notice.Keys)
			}
			log.Printf("fetch: %d notices, %d keys", len(dl), numKeys)
			denylist = dl
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", `application/json`)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(denylist)
		log.Printf("http req: %s\n", r.RequestURI)
	})

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("http error: %s\n", err)
	}
}
