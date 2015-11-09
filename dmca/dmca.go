package dmca

import (
	"bufio"
	"fmt"
	shell "github.com/whyrusleeping/ipfs-shell"
	"html/template"
	"io"
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

	kbytes, err := Asset(kpath)
	if err != nil {
		return nil, nil, err
	}
	s := bufio.NewScanner(strings.NewReader(string(kbytes)))
	for s.Scan() {
		keys = append(keys, s.Text())
	}

	nbytes, err := Asset(npath)
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
		defer nwriter.Close()
	}()

	return keys, nreader, nil
}

func addDenylist(srcpath string, sh *shell.Shell) (string, error) {
	dirs, err := AssetDir(srcpath)
	if err != nil {
		return "", err
	}

	h, err := sh.NewObject("unixfs-dir")
	if err != nil {
		return "", err
	}

	for _, dirname := range dirs {
		keys, nreader, err := keysAndNotice(srcpath + "/" + dirname)
		if err != nil {
			return "", err
		}

		n, err := sh.Add(nreader)
		if err != nil {
			return "", err
		}

		for i, k := range keys {
			n, err = sh.PatchLink(n, fmt.Sprintf("object-%d", i), k, true)
			if err != nil {
				return "", err
			}
		}

		h, err = sh.PatchLink(h, dirname, n, true)
		if err != nil {
			return "", err
		}
	}

	return h, nil
}

func AddDenylist(sh *shell.Shell) (string, error) {
	noticeTemplate = template.Must(template.New("notice").Parse(string(noticeBytes)))

	h, err := addDenylist("dmca/notices", sh)
	if err != nil {
		return "", err
	}

	return h, nil
}
