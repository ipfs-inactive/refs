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
	"regexp"
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

func keysAndNotice(dir string) ([]string, io.Reader, error) {
	kpath := dir + "/keys"
	npath := dir + "/notice.md"
	keys := []string{}

	kbytes, err := ioutil.ReadFile(kpath)
	if err != nil {
		return nil, nil, err
	}
	s := bufio.NewScanner(strings.NewReader(string(kbytes)))
	for s.Scan() {
		keys = append(keys, s.Text())
	}

	nbytes, err := ioutil.ReadFile(npath)
	if err != nil {
		return nil, nil, err
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

	return keys, nreader, nil
}

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
		omit := regexp.MustCompile(`\A(\.|node_modules|versions)`)
		if !dir.IsDir() || omit.Match([]byte(dir.Name())) {
			continue
		}

		keys, nreader, err := keysAndNotice(srcdir.Name() + "/" + dir.Name())
		if err != nil {
			return "", err
		}

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

// Creates a unixfs structure of notice/object tuples,
// e.g. /ipfs/Qmdenylist/2015-09-24-Qmobject/notice (the rendered notice)
// and /ipfs/Qmdenylist/2015-09-24-Qmobject/object (link to Qmobject)
//
// Each object listed in a keys file will get its own tuple with a link
// name of the form <dirname>-<key>.
func main() {
	u := flag.String("uri", "127.0.0.1:5001", "the IPFS API endpoint to use")
	flag.Parse()

	noticeTemplate = template.Must(template.New("notice").Parse(string(noticeBytes)))

	sh := shell.NewShell(*u)

	h, err := addDenylist("./", sh)
	if err != nil {
		log.Fatalf("denylist failed: %s\n", err)
	}

	fmt.Println(h)
}
