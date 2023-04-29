// This program can automatic download language.json file and convert into .go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

//go:generate go run $GOFILE
//go:generate go fmt ./...
func main() {
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		readLang("en_us", f)
		return
	}
	// Pseudo code for get versionURL:
	// $manifest = {https://launchermeta.mojang.com/mc/game/version_manifest.json}
	// $latest = $manifest.latest.release
	// $version = {$manifest.versions[where .id == $latest ].url}
	// $versionURL = $version.assetIndex.url
	versionURL := "https://launchermeta.mojang.com/v1/packages/15809afbdc4489abd52a3d6c3fda0124f8f3404b/1.17.json"
	log.Print("start generating lang packages")

	resp, err := http.Get(versionURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var list struct {
		Objects map[string]struct {
			Hash string
			Size int64
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		log.Fatal(err)
	}

	tasks := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range tasks {
				v := list.Objects[i]
				if strings.HasPrefix(i, "minecraft/lang/") {
					name := i[len("minecraft/lang/") : len(i)-len(".json")]
					lang(name, v.Hash)
				}
			}
		}()
	}
	for i := range list.Objects {
		tasks <- i
	}
	close(tasks)
	wg.Wait()
}

func lang(name, hash string) {
	log.Print("generating ", name, " package")

	//download language
	LangURL := "http://resources.download.minecraft.net/" + hash[:2] + "/" + hash
	resp, err := http.Get(LangURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	readLang(name, resp.Body)
}

// read one language translation
func readLang(name string, r io.Reader) {
	var LangMap map[string]string
	err := json.NewDecoder(r).Decode(&LangMap)
	if err != nil {
		log.Fatal("unmarshal json fail: ", err)
	}
	trans(LangMap)

	pName := strings.ReplaceAll(name, "_", "-")

	// mkdir
	err = os.Mkdir(pName, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	f, err := os.OpenFile(filepath.Join(pName, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genData := struct {
		PkgName string
		Name    string
		LangMap map[string]string
	}{
		PkgName: pName,
		Name:    name,
		LangMap: LangMap,
	}

	tmpl := template.Must(template.New("code_template").Parse(
		`// Code generated by downloader.go; DO NOT EDIT. 
package {{.Name}}
{{if ne .Name "en_us"}}
import "github.com/Tnze/go-mc/chat"

func init() { chat.SetLanguage(Map) }
{{end}}
var Map = {{.LangMap | printf "%#v"}}
`))
	if err := tmpl.Execute(f, genData); err != nil {
		log.Fatal(err)
	}
}

var javaN = regexp.MustCompile(`%[0-9]\$s`)

// Java use %2$s to refer to the second arg, but Golang use %2s, so we need this
func trans(m map[string]string) {
	//replace "%[0-9]\$s" with "%[0-9]s"
	for i := range m {
		c := m[i]
		if javaN.MatchString(c) {
			m[i] = javaN.ReplaceAllStringFunc(c, func(s string) string {
				var index int
				_, err := fmt.Sscanf(s, "%%%d$s", &index)
				if err != nil {
					log.Fatal(err)
				}
				return fmt.Sprintf("%%[%d]s", index)
			})
		}
	}
}