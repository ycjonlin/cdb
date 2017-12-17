package main

import (
	"os"

	"bitbucket.org/sethlin/cdb/pkg/generator"
	"bitbucket.org/sethlin/cdb/pkg/schema"
)

func with(path string, callback func(*os.File)) {
	f, _ := os.Create(path)
	callback(f)
	f.Close()
}

func main() {
	name := os.Args[1]
	path := os.Args[2]
	s := schema.NewSchema(name)
	err := s.Parse(path)
	if err != nil {
		panic(err)
	}
	with(path+".decl.go", func(f *os.File) {
		generator.WriteDeclaration(f, s, "\n", "  ")
	})
	with(path+".marsh.bytes.go", func(f *os.File) {
		generator.WriteBytesMarshalling(f, s, "\n", "  ")
	})
	with(path+".marsh.map.go", func(f *os.File) {
		generator.WriteMapMarshalling(f, s, "\n", "  ")
	})
}
