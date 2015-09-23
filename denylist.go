package main

import (
	"bufio"
	"flag"
	"fmt"
	shell "github.com/whyrusleeping/ipfs-shell"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type noticeData struct {
	Body string
	Keys []string
}

var noticeTemplate *template.Template
var noticeBytes = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8" />
	<title>Unavailable for Legal Reasons</title>
</head>
<body>
	{{ .Body }}
	<h2>Affected Objects</h2>
	<ul>
	{{ range .Keys }}
		<li><a href="/ipfs/{{ . }}">/ipfs/{{ . }}</a></li>
	{{ end }}
	</ul>
</body>
</html>
`

func addDenylist(srcpath string, sh *shell.Shell) (string, error) {
	srcdir, err := os.Open(srcpath)
	if err != nil {
		return "", err
	}

	ndirs, err := srcdir.Readdir(0)
	if err != nil {
		return "", err
	}

	lhash, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}

	for _, dir := range ndirs {
		if !dir.IsDir() || strings.HasPrefix(dir.Name(), ".") {
			continue
		}

		dpath := strings.Join([]string{srcdir.Name(), dir.Name()}, "/")
		kpath := strings.Join([]string{dpath, "keys"}, "/")
		npath := strings.Join([]string{dpath, "notice.md"}, "/")

		kbytes, err := ioutil.ReadFile(kpath)
		if err != nil {
			return "", err
		}
		keys := []string{}
		s := bufio.NewScanner(strings.NewReader(string(kbytes)))
		for s.Scan() {
			keys = append(keys, s.Text())
		}

		nbytes, err := ioutil.ReadFile(npath)
		if err != nil {
			return "", err
		}
		ndata := &noticeData{
			Keys: keys,
			Body: string(nbytes),
		}
		nreader, nwriter := io.Pipe()
		go func() {
			noticeTemplate.Execute(nwriter, ndata)
			nwriter.Close()
		}()
		nhash, err := sh.Add(nreader)
		if err != nil {
			return "", err
		}

		for _, k := range keys {
			dhash, err := sh.NewObject("unixfs-dir")
			if err != nil {
				return "", err
			}

			dhash, err = sh.PatchLink(dhash, "notice", nhash, true)
			if err != nil {
				return "", err
			}

			dhash, err = sh.PatchLink(dhash, "object", k, true)
			if err != nil {
				return "", err
			}

			link := fmt.Sprintf("%s-%s", dir.Name(), k)
			lhash, err = sh.PatchLink(lhash, link, dhash, true)
			if err != nil {
				return "", err
			}
		}
	}

	return lhash, nil
}

func main() {
	u := flag.String("uri", "127.0.0.1:5001", "the IPFS API endpoint to use")
	// p := flag.Bool("pin", false, "pin after adding")
	flag.Parse()

	noticeTemplate = template.Must(template.New("notice").Parse(string(noticeBytes)))

	sh := shell.NewShell(*u)

	h, err := addDenylist("./", sh)
	if err != nil {
		log.Fatalf("denylist failed: %s\n", err)
	}

	fmt.Println(h)
}
