package main

import (
    "os"
    "strings"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "text/template"
    "regexp"
)

var (
    onlyLetters = regexp.MustCompile("[^A-Za-z]")
)

type TestSuite struct {
    SpecFiles []SpecFile
}

type SpecFile struct {
    Overview string
    Tests    []MustacheTest
}

type MustacheTest struct {
    Name      string
    Desc      string
    Data      interface{}
    Template  string
    Partials  map[string]string
    Expected  string
    TestName  string
}

func (mt *MustacheTest) HasData() bool {
    return len(mt.Data.(map[string]interface{})) > 0
}

func (mt *MustacheTest) HasPartials() bool {
    return len(mt.Partials) > 0
}

func main() {
    specs := make([]SpecFile, 0)

    files, _ := filepath.Glob("spec/specs/*.json")
    for _, fileName := range files {

        // skip optional specs
        baseName := filepath.Base(fileName)
        if !strings.HasPrefix(baseName, "~") {

            // parse the json file into a spec
            spec := SpecFile{}
            file, _ := ioutil.ReadFile(fileName)
            err := json.Unmarshal(file, &spec)
            if(err != nil) { panic(err) }

            // titlecase and drop the file extension
            ext := filepath.Ext(fileName)
            specName := strings.Title(baseName[0:len(baseName) - len(ext)])

            // Add each test in this spec file
            for i := 0; i < len(spec.Tests); i += 1 {
                spec.Tests[i].TestName = "Test" + specName + onlyLetters.ReplaceAllString(spec.Tests[i].Name, "")
            }

            specs = append(specs, spec)
        }
    }

    suite := TestSuite{SpecFiles: specs}
    tmpl.Execute(os.Stdout, suite)
}

var tmpl = template.Must(template.New("main").Parse(
`package mustache

import (
    "testing"
    "io/ioutil"
    "os"
)

{{range .SpecFiles}}



// -----------------------------------------------------------------------------
{{range .Tests}}
// {{.Desc}}
func {{.TestName}}(t *testing.T) { {{if .HasPartials}}{{range $k, $v := .Partials}}
    ioutil.WriteFile("{{$k}}", []byte({{printf "%#v" $v}}), 0666)
    defer os.Remove("{{$k}}")
    {{end}}{{end}}
    template := {{printf "%#v" .Template}}{{if .HasData}}
    data     := {{printf "%#v" .Data}}{{end}}
    expected := {{printf "%#v" .Expected}}{{if .HasData}}
    actual   := Render(template, data){{else}}
    actual   := Render(template){{end}}

    if actual != expected {
        t.Errorf("returned %#v, expected %#v", actual, expected)
    }
}
{{end}}{{end}}`))
