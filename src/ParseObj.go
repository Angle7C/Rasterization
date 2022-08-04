package core

import (
	"Matrix/vector"
	"os"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

type MtlData struct {
	TextMap *MipMap
	ka      *vector.Vector3D
	kd      *vector.Vector3D
	ks      *vector.Vector3D
}

func Parse(path, name string) (model *obj.Model, mtls *MtlData) {
	decoder := obj.NewDecoder(obj.DefaultLimits())
	file1, _ := os.Open(path + name)
	defer file1.Close()
	model, _ = decoder.Decode(file1)
	file2, _ := os.Open(path + model.MaterialLibraries[0])
	defer file2.Close()
	l, _ := mtl.NewDecoder(mtl.DefaultLimits()).Decode(file2)
	m := l.Materials[0]
	mm := NewMipMap(path + m.DiffuseTexture)
	mtls.TextMap = mm
	mtls.ka = vector.NewVector3D(l.Materials[0].AmbientColor.R, l.Materials[0].AmbientColor.G, l.Materials[0].AmbientColor.B)
	mtls.kd = vector.NewVector3D(l.Materials[0].DiffuseColor.R, l.Materials[0].DiffuseColor.G, l.Materials[0].DiffuseColor.B)
	mtls.ks = vector.NewVector3D(l.Materials[0].SpecularColor.R, l.Materials[0].SpecularColor.G, l.Materials[0].SpecularColor.B)

	return
}
func ParseNoMtl(path, name string) (*obj.Model, *MipMap) {
	decoder := obj.NewDecoder(obj.DefaultLimits())
	file1, _ := os.Open(path + name)
	defer file1.Close()
	model, _ := decoder.Decode(file1)
	// mm := NewMipMap(path + "floor_diffuse.tga")
	// mm := NewMipMap(path + "crate_1.jpg")
	// mm := NewMipMap(path + "rock.png")
	mm := NewMipMap(path + "spot_texture.png")

	return model, mm
}
