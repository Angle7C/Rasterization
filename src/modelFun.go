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

	m.faceNum = 0
	//设置面
	for _, v := range model.Objects[0].Meshes[0].Faces {

		if len(v.References) == 3 {
			m.faceNum++
			temp := Vertex{}
			for i := 0; i <= 3; i++ {
				temp.normIndex = v.References[i].NormalIndex
				temp.pointIndex = v.References[i].VertexIndex
				temp.textureIndex = v.References[i].TexCoordIndex
				m.face = append(m.face, temp)
			}
		} else {

		}
	}
	m.modelMat = matrix.NewMatrix4()

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

// author: chen
// 2022/8/9
//放缩矩阵
func (model *Model) MakeScales(x, y, z float64) {
	model.modelMat.MulScale(x, y, z)
}

// author: chen
// 2022/8/9
// 旋转矩阵
func (model *Model) MakeRotation(x, y, z float64) {
	m := matrix.NewMatrix4()
	m.MulRoTationX(x).MulRoTationY(y).MulRoTationZ(z)
	model.modelMat = m.Mul(model.modelMat)
}

// author: chen
// 2022/8/9
// 平移矩阵
func (Model *Model) MakeTranslation(x, y, z float64) {
	m := matrix.NewMatrix4()
	m.MulTranslation(x, y, z)
	Model.modelMat = m.Mul(Model.modelMat)
}

// author: chen
// 2022/8/9
// 获取对应面的相关信息
func (Model *Model) getFace(index int64) (i, j, k *Vertex) {
	i = &Model.face[index*3]
	j = &Model.face[index*3+1]
	k = &Model.face[index*3+2]
	return
}

// author: chen
// 2022/8/9
// 渲染程序
func (Model *Model) Render(zbuffer DepthMessage, wg *sync.WaitGroup, image *image.RGBA, h int) {
	var (
		index int64
	)
	for index = 0; index < Model.faceNum; index++ {
		i, j, k := Model.getFace(index)
		Model.renderTriangle(i, j, k, zbuffer, h, light)

	}
}

// author: chen
// 2022/8/9
// 渲染一个三角形
func (Model *Model) renderTriangle(a, b, c *Vertex, zbuffer DepthMessage, h int, light Light) {
	var (
		ap, bp, cp, p *vector.Vector3D
		an, bn, cn, n *vector.Vector3D
		auv, buv, cuv *vector.Vector2D
		ar, br, cr, r *vector.Vector3D
		aw, bw, cw    float64
		rgb           color.RGBA
	)
	ap, bp, cp = Model.getFacePoint(a.pointIndex, b.pointIndex, c.pointIndex)
	an, bn, cn = Model.getFaceNormal(a.normIndex, b.normIndex, c.normIndex)
	auv, buv, cuv = Model.getFaceTexture(a.textureIndex, b.textureIndex, c.textureIndex)
	aw, bw, cw = Model.getFaceWeight(a.pointIndex, b.pointIndex, c.pointIndex)
	ar, br, cr = Model.getFaceRawPoint(a.pointIndex, b.pointIndex, c.pointIndex)
	box := genertorBox(ap, bp, cp)
	for x := box.Min.X; x <= box.Max.X; x++ {
		for y := box.Min.Y; y <= box.Max.Y; y++ {
			p = vector.NewVector3D(float64(x)+0.5, float64(y)+0.5, 1)
			if judge, t := inside(ap, bp, cp, p); judge {
				p.Z = 1.0 / (t.X*aw/ap.Z + t.Y*bw/bp.Z + t.Z*cw/cp.Z)
				if zbuffer.judgeZbuffer(p.Z, x, y, h) {
					n = interpolation3D(an, bn, cn, t)
					r = interpolation3D(ar, br, cr, t)
					rgb = interpolateUV(p, ap, bp, cp, auv, buv, cuv, t, aw, bw, cw, Model.diffuseMap.mip)
					rgb = getColor(
						Model.diffuseMap.value,
						Model.specularMap.value,
						Model.ambientxMap.value,
						r, n, rgb, light, 10,
					)
					zbuffer.setZbuffer(x, y, h, rgb)
				}
			}
		}
	}
}
func (Model *Model) getFacePoint(i, j, k int64) (ap, bp, cp *vector.Vector3D) {
	ap, bp, cp = &Model.point[i], &Model.point[j], &Model.point[k]
	return
}
func (Model *Model) getFaceNormal(i, j, k int64) (an, bn, cn *vector.Vector3D) {
	an, bn, cn = &Model.normal[i], &Model.normal[j], &Model.normal[k]
	return
}
func (Model *Model) getFaceTexture(i, j, k int64) (auv, buv, cuv *vector.Vector2D) {
	return &Model.texture[i], &Model.texture[j], &Model.texture[k]
}
func genertorBox(ap, bp, cp *vector.Vector3D) image.Rectangle {
	xmin := math.Min(ap.X, math.Min(bp.X, cp.X))
	xmax := math.Max(ap.X, math.Max(bp.X, cp.X))
	ymin := math.Min(ap.Y, math.Min(bp.Y, cp.Y))
	ymax := math.Max(ap.Y, math.Max(bp.Y, cp.Y))
	box := image.Rect(int(xmin), int(ymin), int(xmax), int(ymax))
	return box
}
func inside(a, b, c, p *vector.Vector3D) (judge bool, u *vector.Vector3D) {
	var (
		i float64 = (-(p.X-b.X)*(c.Y-b.Y) + (p.Y-b.Y)*(c.X-b.X)) / (-(a.X-b.X)*(c.Y-b.Y) + (a.Y-b.Y)*(c.X-b.X))
		j float64 = (-(p.X-c.X)*(a.Y-c.Y) + (p.Y-c.Y)*(a.X-c.X)) / (-(b.X-c.X)*(a.Y-c.Y) + (b.Y-c.Y)*(a.X-c.X))
		k float64 = 1 - i - j
	)
	judge = (i > 0 && j > 0 && k > 0 && i <= 1 && j <= 1 && k <= 1)
	p.W = 1

	u = vector.NewVector3D(i, j, k)
	return
}
func (Model *Model) getFaceWeight(i, j, k int64) (aw, bw, cw float64) {
	return Model.w[i], Model.w[j], Model.w[k]
}
func (Model *Model) getFaceRawPoint(i, j, k int64) (ar, br, cr *vector.Vector3D) {
	return &Model.rawPoint[i], &Model.rawPoint[j], &Model.rawPoint[k]
}
func interpolation3D(a, b, c, t *vector.Vector3D) *vector.Vector3D {
	temp := matrix.NewMatrix4()
	temp.Set(
		a.X, b.X, c.X, 0,
		a.Y, b.Y, c.Y, 0,
		a.Z, b.Z, c.Z, 0,
		0, 0, 0, 1,
	)
	m := temp.MulVector3D(t)
	return m
}
func interpolateUV(p *vector.Vector3D, ap, bp, cp *vector.Vector3D, au, bu, cu *vector.Vector2D, t *vector.Vector3D,
	w1, w2, w3 float64, textMap *MipMap) color.RGBA {
	u := (t.X*w1/ap.Z*au.X + t.Y*w2/bp.Z*bu.X + t.Z*w3/cp.Z*cu.X) * p.Z
	v := (t.X*w1/ap.Z*au.Y + t.Y*w2/bp.Z*bu.Y + t.Z*w3/cp.Z*cu.Y) * p.Z
	p.X++
	_, t1 := inside(ap, bp, cp, p)
	p.X--
	z1 := 1.0 / (t1.X*w1/ap.Z + t1.Y*w2/bp.Z + t1.Z*w3/cp.Z)
	u1 := (t1.X*w1/ap.Z*au.X + t1.Y*w2/bp.Z*bu.X + t1.Z*w3/cp.Z*cu.X) * z1
	v1 := (t1.X*w1/ap.Z*au.Y + t1.Y*w2/bp.Z*bu.Y + t1.Z*w3/cp.Z*cu.Y) * z1
	p.Y++
	_, t2 := inside(ap, bp, cp, p)
	p.Y--
	z2 := 1.0 / (t2.X*w1/ap.Z + t2.Y*w2/bp.Z + t2.Z*w3/cp.Z)
	u2 := (t2.X*w1/ap.Z*au.X + t2.Y*w2/bp.Z*bu.X + t2.Z*w3/cp.Z*cu.X) * z2
	v2 := (t2.X*w1/ap.Z*au.Y + t2.Y*w2/bp.Z*bu.Y + t2.Z*w3/cp.Z*cu.Y) * z2
	r := textMap.UV(u, v, u1, v1, u2, v2)
	return r
}
func getColor(kd, ks, ka *vector.Vector3D, Pos, Normal *vector.Vector3D, c color.RGBA, light Light, p float64) color.RGBA {
	n := Normal.Normal()
	var (
		vd *vector.Vector3D = vector.NewVector3D(0, 0, 0)
	)
	for i := 0; i < len(light.Pos); i++ {
		length := (light.Pos[i].Subto(Pos)).Normality2()
		v := light.eyePos.Subto(Pos).Normal()
		l := light.Pos[i].Subto(Pos).Normal()
		h := l.Add(v).Normal()
		rate := math.Max(n.Dot(l), 0) * 100 / length
		r := float64(c.R) * rate * light.Indent[0].X
		g := float64(c.G) * rate * light.Indent[0].Y
		b := float64(c.B) * rate * light.Indent[0].Z

		Ld := vector.NewVector3D(r, g, b)
		rate = math.Pow(math.Max(0, v.Dot(h)), p) * 10000 / length
		r = ks.X * rate * float64(c.R)
		g = ks.Y * rate * float64(c.G)
		b = ks.Z * rate * float64(c.B)
		Ls := vector.NewVector3D(r, g, b)
		r = ka.X * float64(c.R)
		g = ka.Y * float64(c.G)
		b = ka.Z * float64(c.B)
		La := vector.NewVector3D(r, g, b)
		vd = vd.Add(Ld.Add(Ls).Add(La))
	}
	return color.RGBA{uint8(clamp(vd.X, 0, 255)), uint8(clamp(vd.Y, 0, 255)), uint8(clamp(vd.Z, 0, 255)), uint8(clamp(vd.W, 0, 255))}
}
