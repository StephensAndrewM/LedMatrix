package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const BASE_DIR = "icons/weather/"

func main() {
	fs, err := ioutil.ReadDir(BASE_DIR)
	if err != nil {
		panic(err)
	}
	out, err := os.Create("weathericons.go")
	if err != nil {
		panic(err)
	}
	out.Write([]byte("package main \n\nvar weatherIcons = map[string]string{\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".png") {
			out.Write([]byte("`" + strings.TrimSuffix(f.Name(), ".png") + "`:`"))
			f, err := os.Open(BASE_DIR + f.Name())
			if err != nil {
				panic(err)
			}
			encoder := base64.NewEncoder(base64.StdEncoding, out)
			_, err = io.Copy(encoder, f)
			if err != nil {
				panic(err)
			}
			encoder.Close()
			out.Write([]byte("`,\n"))
		}
	}
	out.Write([]byte("}\n"))
	out.Close()
}
