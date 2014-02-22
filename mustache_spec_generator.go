package main

import (
    "os"
    "fmt"
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
    Name     string
    Overview string
    Tests    []MustacheTest
}

func (sf *SpecFile) FunctionName() string {
    str := onlyLetters.ReplaceAllString(sf.Name, "")
    return strings.Title(str[0:len(str)-4])
}

type MustacheTest struct {
    Name      string
    Desc      string
    Data      interface{}
    Template  string
    Partials  map[string]string
    Expected  string
}

func (mt *MustacheTest) FunctionName() string {
    return onlyLetters.ReplaceAllString(mt.Name, "")
}

func (mt *MustacheTest) TemplateStr() string {
    return fmt.Sprintf("%#v", mt.Template)
}

func (mt *MustacheTest) ExpectedStr() string {
    return fmt.Sprintf("%#v", mt.Expected)
}

func (mt *MustacheTest) DataStr() string {
    data, _ := json.Marshal(mt.Data)
    return fmt.Sprintf("%#v", string(data))
}


func main() {
    files, _ := filepath.Glob("spec/specs/*.json")
    specs := make([]SpecFile, len(files) - 1)

    for i, fileName := range files {

        // lambdas are not supported
        baseName := filepath.Base(fileName)
        if baseName != "~lambdas.json" {

            file, _ := ioutil.ReadFile(fileName)
            err := json.Unmarshal(file, &specs[i])
            if(err != nil) { panic(err) }
            specs[i].Name = baseName
        }
    }

    suite := TestSuite{SpecFiles: specs}
    tmpl.Execute(os.Stdout, suite)
}

var tmpl = template.Must(template.New("main").Parse(
`package mustache

import (
    "testing"
    "encoding/json"
)

func mustDecodeJson (str string) interface{} {
    var data interface{}
    err := json.Unmarshal([]byte(str), &data)
    if(err != nil) { panic(err) }
    return data
}

{{range .SpecFiles}}
{{$SpecFunctionName := .FunctionName}}
// -----------------------------------------------------------------------------

{{range .Tests}}
// {{.Desc}}
func Test{{$SpecFunctionName}}{{.FunctionName}}(t *testing.T) {
    template := {{.TemplateStr}}
    data     := mustDecodeJson({{.DataStr}})
    actual   := Render(template, data)
    expected := {{.ExpectedStr }}

    if actual != expected {
        t.Errorf("returned %#v, expected %#v", actual, expected)
    }
}
{{end}}
{{end}}`))
