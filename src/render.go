package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"image"
	"sync"

	"github.com/mokiat/go-data-front/decoder/obj"
)

type RenderPipline interface {
	MakeModelMat(x, y, z float64)
	MakeViewMat(eye, target, up *vector.Vector3D)
	MakePerspectMat(fov, aspect, near, far float64)
	GeometryShader()
	VertexShader()
	FragmentShader()
	MV()
	P()
	ToScreen()
	Render()
}
type PointV vector.Vector3D
type UV vector.Vector3D
type NormalV vector.Vector3D
type Render struct {
	Points      []PointV
	Face        []Face
	TextColor   []UV
	Normal      []NormalV
	modelMat    *matrix.Matrix4
	viewMat     *matrix.Matrix4
	perspectMat *matrix.Matrix4
	viewPort    *matrix.Matrix4
}

func NewRenderObj(o *obj.Model) (r *Render, err error) {
	r = &Render{}
	//设置点
	// r.Points = make([]vector.Vector3D, len(o.Vertices))
	for _, v := range o.Vertices {
		r.Points = append(r.Points, PointV(v))
	}
	//设置纹理
	// r.TextColor = make([]vector.Vector3D, len(o.TexCoords))
	for _, v := range o.TexCoords {
		r.TextColor = append(r.TextColor, UV{v.U, v.V, v.W, 1})
	}
	//设置法向量
	// r.Normal = make([]vector.Vector3D, len(o.Normals))
	for _, v := range o.Normals {
		r.Normal = append(r.Normal, NormalV{v.X, v.Y, v.Z, 0})
	}
	// 设置面
	r.Face = make([]Face, len(o.Objects[0].Meshes[0].Faces))
	for i, v := range o.Objects[0].Meshes[0].Faces {
		for _, vv := range v.References {
			r.Face[i].PointIndex = append(r.Face[i].PointIndex, vv.VertexIndex)
			r.Face[i].TextCoords = append(r.Face[i].TextCoords, vv.TexCoordIndex)
			r.Face[i].NormalIndex = append(r.Face[i].NormalIndex, vv.NormalIndex)
		}
	}
	r.modelMat = matrix.NewMatrix4()
	r.viewMat = matrix.NewMatrix4()
	r.perspectMat = matrix.NewMatrix4()
	r.viewPort = matrix.NewMatrix4()
	return

}
func (r *Render) MakeModelMat(x, y, z float64) {
	r.modelMat = matrix.NewMatrix4().MulScale(
		x, y, z)
}
func (r *Render) SetModelMat(temp *matrix.Matrix4) {
	r.modelMat = temp
}
func (r *Render) MakeViewMat(eye, target, up *vector.Vector3D) {
	r.viewMat = r.viewMat.LookAt(eye, target, up)
}

func (r *Render) MakePerspectMat(fov, aspect, near, far float64) {
	r.perspectMat = r.perspectMat.MakePerspective(fov, aspect, near, far)
}
func (r *Render) MakeViewPort(w, h float64) {
	r.viewPort.Set(
		w/2, 0, 0, w/2,
		0, h/2, 0, h/2,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

func (m *Render) ToScreen() {
	// for k, v := range m.Points {
	// 	v = *m.viewPort.MulVector3D(&v)
	// 	m.Points[k] = v
	// // }
}
func (r *Render) Render(image *image.RGBA) {
}
func (face *Face) renderFace(zbuffer []float64, point []vector.Vector3D, Color []vector.Vector3D, RW *sync.RWMutex, image *image.RGBA) {

}
func (face *Face) BlinPhong(zbuffer []float64, point []vector.Vector3D, Color []vector.Vector3D, RW *sync.RWMutex, image *image.RGBA) {

}
func GeometryShader() {

}
func (p *Point) VertexShader() {

}
func FragmentShader() {

}
func inside(a, b, c, p *vector.Vector3D) (judge bool, u *vector.Vector3D) {
	var (
		i float64 = (-(p.X-b.X)*(c.Y-b.Y) + (p.Y-b.Y)*(c.X-b.X)) / (-(a.X-b.X)*(c.Y-b.Y) + (a.Y-b.Y)*(c.X-b.X))
		j float64 = (-(p.X-c.X)*(a.Y-c.Y) + (p.Y-c.Y)*(a.X-c.X)) / (-(b.X-c.X)*(a.Y-c.Y) + (b.Y-c.Y)*(a.X-c.X))
		k float64 = 1 - i - j
	)
	judge = (i >= 0 && j >= 0 && k >= 0 && i <= 1 && j <= 1 && k <= 1)
	u = vector.NewVector3D(i, j, k)
	return
}
