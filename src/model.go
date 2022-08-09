package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/mokiat/go-data-front/decoder/mtl"
)

type numbers interface {
	float32 | float64 | uint8 | int
}
type Modeler interface {
	loadMaterial(m *mtl.Material, path string)
	loadTexture(path, name, ty string)
}
type DepthMessage struct {
	Depth  float64
	Colors color.RGBA
}
type Model struct {
	diffuseMap  *MapData
	normalMap   *MapData
	heightMap   *MapData
	ambientxMap *MapData
	specularMap *MapData
	point       []vector.Vector3D
	rawPoint    []vector.Vector3D
	face        []int64
	w           []float64
	normal      []vector.Vector3D
	texture     []vector.Vector2D
	modelMat    *matrix.Matrix4
}
type MapData struct {
	value *vector.Vector3D
	mip   *MipMap
}

type MipMaper interface {
	UV(l, u, v float64)
	GetRGBA()
}
type MipMap struct {
	Name       string
	Num        int
	w          int
	h          int
	LevelImage []image.Image
}

func init() {
	f, err := os.OpenFile("E:/MyPro/GO/3Dobj/src/dariy.log", os.O_APPEND, 0777)

	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}
