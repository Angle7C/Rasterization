package core

import (
	"os"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

func Parse(path, name string) (*obj.Model, *MipMap) {
	decoder := obj.NewDecoder(obj.DefaultLimits())
	file1, _ := os.Open(path + name)

	defer file1.Close()
	model, _ := decoder.Decode(file1)
	file2, _ := os.Open(path + model.MaterialLibraries[0])
	defer file2.Close()
	l, _ := mtl.NewDecoder(mtl.DefaultLimits()).Decode(file2)
	m := l.Materials[0]
	mm := NewMipMap(path + m.DiffuseTexture)
	return model, mm
}
