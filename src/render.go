package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"fmt"
	"image"
	"math"
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
	MVP()
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

func (r *Render) MVP() {
	mvp := r.perspectMat.Mul(r.viewMat.Mul(r.modelMat))
	for i := 0; i < len(r.Points); i++ {
		vd := mvp.MulVector3D((*vector.Vector3D)(&r.Points[i]))
		vd.X /= vd.W
		vd.Y /= vd.W
		vd.Z /= vd.W
		vd.W /= vd.W
		r.Points[i] = PointV(*vd)
	}
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
	for k, v := range m.Points {
		v = PointV(*m.viewPort.MulVector3D((*vector.Vector3D)(&v)))
		m.Points[k] = v
	}
}
func (r *Render) Render(image *image.RGBA, mtl *MipMap) {
	var (
		WG      *sync.WaitGroup = &sync.WaitGroup{}
		RW      *sync.RWMutex   = &sync.RWMutex{}
		w       int             = image.Rect.Dx()
		h       int             = image.Rect.Dy()
		zbuffer []float64       = make([]float64, w*h)
	)
	for i := 0; i < len(zbuffer); i++ {
		zbuffer[i] = math.Inf(1)
	}
	r.MVP()
	r.ToScreen()
	if len(r.Face) < 4 {
		r.FragmentShader(0, len(r.Face), zbuffer, WG, RW, image, mtl)
	} else {
		// WG.Add(4)
		// ds := len(r.Face) / 4
		go r.FragmentShader(0, 2, zbuffer, WG, RW, image, mtl)
		// go r.FragmentShader(ds, 2*ds, zbuffer, WG, RW, image, mtl)
		// go r.FragmentShader(2*ds, 3*ds, zbuffer, WG, RW, image, mtl)
		// go r.FragmentShader(3*ds, len(r.Face), zbuffer, WG, RW, image, mtl)
		// WG.Wait()
	}
}
func (face *Face) renderFace(zbuffer []float64, point []PointV, TextUV []UV, RW *sync.RWMutex, image *image.RGBA) {
	// var (
	// 	au, bu, cu UV
	// 	ap, bp, cp PointV
	// 	xmin       float64
	// 	xmax       float64
	// 	ymin       float64
	// 	ymax       float64
	// )
	// au, bu, cu = TextUV[face.PointIndex[0]], TextUV[face.PointIndex[1]], TextUV[face.PointIndex[2]]
	// ap, bp, cp = point[face.PointIndex[0]], point[face.PointIndex[1]], point[face.PointIndex[2]]
	// xmin = math.Min(ap.X, math.Min(bp.X, cp.X))
	// xmax = math.Max(ap.X, math.Max(bp.X, cp.X))
	// ymin = math.Min(ap.Y, math.Min(bp.Y, cp.Y))
	// ymax = math.Max(ap.Y, math.Max(bp.Y, cp.Y))
	// for x := xmin; x < xmax; x++ {
	// 	for y := ymin; y < ymax; y++ {

	// 	}
	// }

}
func (face *Face) BlinPhong(zbuffer []float64, point []PointV, TextUV []UV, RW *sync.RWMutex, image *image.RGBA, mtl *MipMap) {
	var (
		au, bu, cu UV
		ap, bp, cp PointV
		xmin       float64
		xmax       float64
		ymin       float64
		ymax       float64
		w          int = image.Rect.Dx()
		h          int = image.Rect.Dy()
		u          float64
		v          float64
	)
	au, bu, cu = TextUV[face.PointIndex[0]], TextUV[face.PointIndex[1]], TextUV[face.PointIndex[2]]
	ap, bp, cp = point[face.PointIndex[0]], point[face.PointIndex[1]], point[face.PointIndex[2]]
	xmin = math.Min(ap.X, math.Min(bp.X, cp.X))
	xmax = math.Max(ap.X, math.Max(bp.X, cp.X))
	ymin = math.Min(ap.Y, math.Min(bp.Y, cp.Y))
	ymax = math.Max(ap.Y, math.Max(bp.Y, cp.Y))
	for x := int(xmin); x >= 0 && x < w && x < int(xmax); x++ {
		for y := int(ymin); y >= 0 && y < h && y < int(ymax); y++ {
			if judge, t := inside(&ap, &bp, &cp, vector.NewVector3D(float64(x)+0.5, float64(y)+0.5, 1)); judge {
				z := 1.0 / (t.X/ap.Z + t.Y/bp.Z + t.Z/cp.Z)
				u = (t.X/cp.Z*au.X + t.Y/cp.Z*bu.X + t.Z/cp.Z*cu.X) * z
				v = (t.X/cp.Z*au.Y + t.Y/cp.Z*bu.Y + t.Z/cp.Z*cu.Y) * z
				RW.RLock()
				if zbuffer[x*h+y] >= z {
					RW.RUnlock()
					RW.Lock()
					image.Set(int(x), int(y), mtl.UV(0, u, v))
					zbuffer[x*h+y] = z
					RW.Unlock()
				} else {
					RW.RUnlock()
				}

			}
		}
	}
}
func GeometryShader() {

}
func (p *Point) VertexShader() {

}
func (r *Render) FragmentShader(s, e int, zbuffer []float64, wg *sync.WaitGroup, rw *sync.RWMutex, imag *image.RGBA, mtl *MipMap) {

	for i := s; i < e; i++ {
		r.Face[i].BlinPhong(zbuffer, r.Points, r.TextColor, rw, imag, mtl)
	}
	fmt.Printf("this process is %v-%v\n", s, e)
	// wg.Done()
}
func inside(a, b, c *PointV, p *vector.Vector3D) (judge bool, u *vector.Vector3D) {
	var (
		i float64 = (-(p.X-b.X)*(c.Y-b.Y) + (p.Y-b.Y)*(c.X-b.X)) / (-(a.X-b.X)*(c.Y-b.Y) + (a.Y-b.Y)*(c.X-b.X))
		j float64 = (-(p.X-c.X)*(a.Y-c.Y) + (p.Y-c.Y)*(a.X-c.X)) / (-(b.X-c.X)*(a.Y-c.Y) + (b.Y-c.Y)*(a.X-c.X))
		k float64 = 1 - i - j
	)
	judge = (i >= 0 && j >= 0 && k >= 0 && i <= 1 && j <= 1 && k <= 1)
	u = vector.NewVector3D(i, j, k)
	return
}
