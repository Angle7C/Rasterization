package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"image/color"
	"log"
	"os"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

// author: chen
// 2022/8/8
// 判断[x*h+y]位置的zbuffer是否大于depth，并返回判断结果和相应位置的颜色
func judge(zbuffer []DepthMessage, depth float64, x, y, h int) (judge bool, color color.RGBA) {
	judge = zbuffer[x*h+y].Depth > depth
	color = zbuffer[x*h+y].Colors
	return
}

// author: chen
// 2022/8/8
// 将值限制在[min,max]之间
// 支持泛型 float，int
func clamp[T numbers](value, min, max T) T {
	if value > 255 {
		value = 255
	} else if value < 0 {
		value = 0
	}
	return value
}

// author: chen
// 2022/8/9
// 获取不同贴图下的颜色值或反射系数
func (m *MapData) getValue(l, u, v float64) *vector.Vector3D {
	if m.mip == nil {
		return m.value
	} else {
		return m.mip.NearUV(u, v)
	}
}

// author: chen
// 2022/8/9
// 构建一个模型数据
func NewModel(path, name string) (m *Model, err error) {
	var (
		material   []string
		mtlLibrary *mtl.Library
		mtlDecoder mtl.Decoder
		model      *obj.Model
	)
	m = &Model{}
	model, material, err = readObjFile(path + name)
	mtlDecoder = mtl.NewDecoder(mtl.DefaultLimits())
	//配置对应的texture
	for _, v := range material {
		mtl, _ := os.Open(path + v)
		defer mtl.Close()
		mtlLibrary, err = mtlDecoder.Decode(mtl)
		m.loadMaterial(mtlLibrary.Materials[0], path)
	}
	//设置点
	for _, v := range model.Vertices {
		m.point = append(m.point, vector.Vector3D(v))
		m.rawPoint = append(m.rawPoint, vector.Vector3D(v))
	}
	//设置uv
	for _, v := range model.TexCoords {
		m.texture = append(m.texture, vector.Vector2D{v.U, v.V, v.W})
	}
	//设置法向量
	for _, v := range model.Normals {
		m.normal = append(m.normal, *vector.NewVector3D(v.X, v.Y, v.Z))
	}
	m.modelMat = matrix.NewMatrix4()
	//设置面
	for _, v := range model.Objects[0].Meshes[0].Faces {
		if len(v.References) == 3 {
			m.face = append(m.face, v.References[0].VertexIndex, v.References[1].VertexIndex, v.References[2].VertexIndex)
		} else {
			a, b, c, d := m.rawPoint[v.References[0].VertexIndex], m.rawPoint[v.References[1].VertexIndex], m.rawPoint[v.References[2].VertexIndex], m.rawPoint[v.References[3].VertexIndex]
			ab, ac, ad, bc, cd := a.Subto(&b).Normality2(), a.Subto(&c).Normality2(), a.Subto(&d).Normality2(), b.Subto(&c).Normality2(), c.Subto(&d).Normality2()
			if ab > ac && ab > ad && ab > bc && ab > cd {
				m.face = append(m.face, v.References[0].VertexIndex, v.References[1].VertexIndex, v.References[2].VertexIndex)
				m.face = append(m.face, v.References[0].VertexIndex, v.References[1].VertexIndex, v.References[3].VertexIndex)
			} else if ac > ab && ac > ad && ac > bc && ac > cd {
				m.face = append(m.face, v.References[0].VertexIndex, v.References[2].VertexIndex, v.References[1].VertexIndex)
				m.face = append(m.face, v.References[0].VertexIndex, v.References[2].VertexIndex, v.References[3].VertexIndex)
			} else if ad > ab && ad > ac && ad > bc && ad > cd {
				m.face = append(m.face, v.References[0].VertexIndex, v.References[3].VertexIndex, v.References[2].VertexIndex)
				m.face = append(m.face, v.References[0].VertexIndex, v.References[3].VertexIndex, v.References[1].VertexIndex)
			} else if bc > ab && bc > ac && ac > ad && ac > cd {
				m.face = append(m.face, v.References[1].VertexIndex, v.References[2].VertexIndex, v.References[0].VertexIndex)
				m.face = append(m.face, v.References[1].VertexIndex, v.References[2].VertexIndex, v.References[3].VertexIndex)
			} else {
				m.face = append(m.face, v.References[2].VertexIndex, v.References[3].VertexIndex, v.References[1].VertexIndex)
				m.face = append(m.face, v.References[2].VertexIndex, v.References[3].VertexIndex, v.References[0].VertexIndex)
			}

		}
	}
	return
}

// author: chen
// 2022/8/9
// 读取对应的obj文件
func readObjFile(path string) (model *obj.Model, material []string, err error) {
	var (
		decoder   obj.Decoder
		modelFile *os.File
		// err       error
	)
	decoder = obj.NewDecoder(obj.DefaultLimits())
	modelFile, err = os.Open(path)
	if err != nil {
		log.Printf("无法打开obj文件 path:%v", path)
		return nil, nil, err
	}
	defer modelFile.Close()
	model, err = decoder.Decode(modelFile)
	if err != nil {
		log.Printf("无法解析obj文件 path:%v", path)
		return nil, nil, err
	}
	material = model.MaterialLibraries
	return model, material, nil
}

// author: chen
// 2022/8/9
// 读取对应的mtl文件的材质设置
func (model *Model) loadMaterial(m *mtl.Material, path string) {
	model.ambientxMap = &MapData{}
	model.diffuseMap = &MapData{}
	model.specularMap = &MapData{}
	model.ambientxMap.value = vector.NewVector3D(m.AmbientColor.R, m.AmbientColor.G, m.AmbientColor.B)
	model.loadTexture(path, m.AmbientTexture, "ambientx")
	model.diffuseMap.value = vector.NewVector3D(m.DiffuseColor.R, m.DiffuseColor.G, m.DiffuseColor.B)
	model.loadTexture(path, m.DiffuseTexture, "diffuse")
	model.specularMap.value = vector.NewVector3D(m.SpecularColor.R, m.SpecularColor.G, m.SpecularColor.B)
	model.loadTexture(path, m.SpecularTexture, "specular")
	model.loadTexture(path, m.BumpTexture, "normal")
}

// author: chen
// 2022/8/9
// 加载mtl下的纹理文件
func (model *Model) loadTexture(path, name, ty string) {
	switch ty {
	case "ambientx":
		model.ambientxMap.mip = NewMipMap(path + name)
	case "diffuse":
		model.diffuseMap.mip = NewMipMap(path + name)
	case "specular":
		model.specularMap.mip = NewMipMap(path + name)
	case "bump":
		model.normalMap.mip = NewMipMap(path + name)
	}
}
func (model *Model) MakeScales(x, y, z float64) {
	model.modelMat.MulScale(x, y, z)
}
func (model *Model) MakeRotation(x, y, z float64) {
	m := matrix.NewMatrix4()
	m.MulRoTationX(x).MulRoTationY(y).MulRoTationZ(z)
	model.modelMat = m.Mul(model.modelMat)
}
func (Model *Model) MakeTranslation(x, y, z float64) {
	m := matrix.NewMatrix4()
	m.MulTranslation(x, y, z)
	Model.modelMat = m.Mul(Model.modelMat)
}
func (Model *Model) getFace(index int) (i, j, k int64) {
	i = Model.face[index*3]
	j = Model.face[index*3+1]
	k = Model.face[index*3+2]
	return
}
