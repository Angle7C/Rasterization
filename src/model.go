package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sync"

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
	Depth  []float64
	Colors []color.RGBA
	rw     *sync.RWMutex
}
type Vertex struct {
	normIndex    int64
	pointIndex   int64
	textureIndex int64
}
type Model struct {
	diffuseMap  *MapData
	normalMap   *MapData
	heightMap   *MapData
	ambientxMap *MapData
	specularMap *MapData
	point       []vector.Vector3D
	rawPoint    []vector.Vector3D
	face        []Vertex
	faceNum     int64
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
func NewZbuffer(x, y int) *DepthMessage {
	temp := &DepthMessage{}
	temp.Colors = make([]color.RGBA, x*y)
	temp.Depth = make([]float64, x*y)
	temp.rw = new(sync.RWMutex)
	for k, _ := range temp.Colors {
		temp.Depth[k] = math.MaxFloat64
	}
	return temp
}
