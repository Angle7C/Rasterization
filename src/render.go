package core

// import (
// 	"Matrix/matrix"
// 	"Matrix/vector"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"math"
// 	"sync"

// 	"github.com/mokiat/go-data-front/decoder/obj"
// )

// type RenderPipline interface {
// 	MakeModelMat(x, y, z float64)
// 	MakeViewMat(eye, target, up *vector.Vector3D)
// 	MakePerspectMat(fov, aspect, near, far float64)
// 	FragmentShader()
// 	MVP()
// 	ToScreen()
// 	Render()
// }
// type PointV vector.Vector3D
// type UV vector.Vector3D
// type NormalV vector.Vector3D

// type Render struct {
// 	Points      []PointV
// 	Face        []Face
// 	TextColor   []UV
// 	Normal      []NormalV
// 	W           []float64
// 	modelMat    *matrix.Matrix4
// 	viewMat     *matrix.Matrix4
// 	perspectMat *matrix.Matrix4
// 	viewPort    *matrix.Matrix4
// }

// var light Light

// func MakeLight(pos, iden, eyepos vector.Vector3D) {
// 	pos.W = 1
// 	iden.W = 1
// 	eyepos.W = 1
// 	light.Pos[0] = pos
// 	light.Indent[0] = iden
// 	light.eyePos = eyepos
// }
// func (r *Render) MVP() []PointV {
// 	mv := r.viewMat.Mul(r.modelMat)
// 	r.W = make([]float64, len(r.Points))
// 	list := make([]PointV, len(r.Points))
// 	for i := 0; i < len(list); i++ {
// 		v := r.modelMat.MulVector3D((*vector.Vector3D)(&r.Points[i]))
// 		if r.Normal != nil {
// 			n := r.modelMat.MulVector3D((*vector.Vector3D)(&r.Normal[i]))
// 			r.Normal[i] = NormalV(*n.Normal())
// 		}
// 		list[i] = PointV(*v)
// 		vd := mv.MulVector3D((*vector.Vector3D)(&r.Points[i]))
// 		vd = r.perspectMat.MulVector3D(vd)
// 		vd.X /= vd.W
// 		vd.Y /= vd.W
// 		vd.Z /= vd.W
// 		r.W[i] = 1.0 / vd.W
// 		vd.W /= vd.W
// 		r.Points[i] = PointV(*vd)
// 	}
// 	return list
// }
// func NewRenderObj(o *obj.Model) (r *Render, err error) {
// 	r = &Render{}
// 	//设置点
// 	// r.Points = make([]vector.Vector3D, len(o.Vertices))
// 	for _, v := range o.Vertices {
// 		r.Points = append(r.Points, PointV(v))
// 	}
// 	//设置纹理
// 	// r.TextColor = make([]vector.Vector3D, len(o.TexCoords))
// 	for _, v := range o.TexCoords {
// 		r.TextColor = append(r.TextColor, UV{v.U, v.V, v.W, 1})
// 	}
// 	//设置法向量
// 	// r.Normal = make([]vector.Vector3D, len(o.Normals))
// 	for _, v := range o.Normals {
// 		r.Normal = append(r.Normal, NormalV{v.X, v.Y, v.Z, 0})
// 	}
// 	// 设置面
// 	r.Face = make([]Face, len(o.Objects[0].Meshes[0].Faces))
// 	for i, v := range o.Objects[0].Meshes[0].Faces {
// 		for _, vv := range v.References {
// 			r.Face[i].PointIndex = append(r.Face[i].PointIndex, vv.VertexIndex)
// 			r.Face[i].TextCoords = append(r.Face[i].TextCoords, vv.TexCoordIndex)
// 			r.Face[i].NormalIndex = append(r.Face[i].NormalIndex, vv.NormalIndex)
// 		}
// 	}
// 	r.modelMat = matrix.NewMatrix4()
// 	r.viewMat = matrix.NewMatrix4()
// 	r.perspectMat = matrix.NewMatrix4()
// 	r.viewPort = matrix.NewMatrix4()
// 	return

// }
// func (r *Render) MakeModelMat(x, y, z float64) {
// 	r.modelMat = matrix.NewMatrix4().MulScale(
// 		x, y, z)
// }
// func (r *Render) MakeTranslation(x, y, z float64) {
// 	r.modelMat = r.modelMat.MulTranslation(x, y, z)
// }
// func (r *Render) MakeRoation(x, y, z float64) {
// 	r.modelMat = r.modelMat.MulRoTationX(x).MulRoTationY(y).MulRoTationZ(z)
// }
// func (r *Render) SetModelMat(temp *matrix.Matrix4) {
// 	r.modelMat = temp
// }
// func (r *Render) MakeViewMat(eye, target, up *vector.Vector3D) {
// 	r.viewMat = r.viewMat.LookAt(eye, target, up)
// }

// func (r *Render) MakePerspectMat(fov, aspect, near, far float64) {
// 	r.perspectMat = r.perspectMat.MakePerspective(fov, aspect, near, far)
// }
// func (r *Render) MakeViewPort(w, h float64) {
// 	r.viewPort.Set(
// 		w/2, 0, 0, w/2,
// 		0, h/2, 0, h/2,
// 		0, 0, 1, 0,
// 		0, 0, 0, 1,
// 	)
// }

// func (m *Render) ToScreen() {
// 	for k, v := range m.Points {
// 		v = PointV(*m.viewPort.MulVector3D((*vector.Vector3D)(&v)))
// 		m.Points[k] = v
// 	}
// }
// func (r *Render) Render(image *image.RGBA, mtl *MtlData) {
// 	var (
// 		WG      *sync.WaitGroup = &sync.WaitGroup{}
// 		RW      *sync.RWMutex   = &sync.RWMutex{}
// 		w       int             = image.Rect.Dx()
// 		h       int             = image.Rect.Dy()
// 		zbuffer []float64       = make([]float64, w*h)
// 	)
// 	for i := 0; i < len(zbuffer); i++ {
// 		zbuffer[i] = math.Inf(1)
// 	}
// 	viewPos := r.MVP()
// 	r.ToScreen()
// 	if len(r.Face) < 4 {
// 		WG.Add(1)
// 		r.FragmentShader(0, len(r.Face), zbuffer, r.W, WG, RW, image, mtl, viewPos)
// 	} else {
// 		WG.Add(4)
// 		ds := len(r.Face) / 4
// 		go r.FragmentShader(0, ds, zbuffer, r.W, WG, RW, image, mtl, viewPos)
// 		go r.FragmentShader(ds, 2*ds, zbuffer, r.W, WG, RW, image, mtl, viewPos)
// 		go r.FragmentShader(2*ds, 3*ds, zbuffer, r.W, WG, RW, image, mtl, viewPos)
// 		go r.FragmentShader(3*ds, len(r.Face), zbuffer, r.W, WG, RW, image, mtl, viewPos)
// 		WG.Wait()
// 	}
// }
// func (face *Face) renderFace(zbuffer []float64, point []PointV, TextUV []UV, RW *sync.RWMutex, image *image.RGBA) {
// }
// func (face *Face) BlinPhong(zbuffer, W []float64, point []PointV, TextUV []UV, Normal []NormalV,
// 	RW *sync.RWMutex, image *image.RGBA, mtl *MtlData, viewPos []PointV) {
// 	var (
// 		au, bu, cu UV
// 		ap, bp, cp PointV
// 		an, bn, cn NormalV
// 		av, bv, cv PointV
// 		colors     color.RGBA
// 		n          vector.Vector3D
// 		pos        vector.Vector3D
// 		xmin       float64
// 		xmax       float64
// 		ymin       float64
// 		ymax       float64
// 		w          int = image.Rect.Dx()
// 		h          int = image.Rect.Dy()
// 	)
// 	w1, w2, w3 := W[face.PointIndex[0]], W[face.PointIndex[1]], W[face.PointIndex[2]]
// 	au, bu, cu = TextUV[face.TextCoords[0]], TextUV[face.TextCoords[1]], TextUV[face.TextCoords[2]]
// 	ap, bp, cp = point[face.PointIndex[0]], point[face.PointIndex[1]], point[face.PointIndex[2]]
// 	if Normal != nil {
// 		an, bn, cn = Normal[face.NormalIndex[0]], Normal[face.NormalIndex[1]], Normal[face.NormalIndex[2]]
// 	}
// 	av, bv, cv = viewPos[face.PointIndex[0]], viewPos[face.PointIndex[1]], viewPos[face.PointIndex[2]]
// 	xmin = math.Min(ap.X, math.Min(bp.X, cp.X))
// 	xmax = math.Max(ap.X, math.Max(bp.X, cp.X))
// 	ymin = math.Min(ap.Y, math.Min(bp.Y, cp.Y))
// 	ymax = math.Max(ap.Y, math.Max(bp.Y, cp.Y))
// 	//包围盒
// 	for x := int(xmin); x >= 0 && x < w && x <= int(xmax); x++ {
// 		for y := int(ymin); y >= 0 && y < h && y <= int(ymax); y++ {
// 			if judge, t := inside(&ap, &bp, &cp, vector.NewVector3D(float64(x)+0.5, float64(y)+0.5, 1)); judge {

// 				z := 1.0 / (t.X*w1/ap.Z + t.Y*w2/bp.Z + t.Z*w3/cp.Z)
// 				//uv 三线性插值
// 				if Normal != nil {
// 					n.X = (t.X*an.X + t.Y*bn.X + t.Z*cn.X)
// 					n.Y = (t.X*an.Y + t.Y*bn.Y + t.Z*cn.Y)
// 					n.Z = (t.X*an.Z + t.Y*bn.Z + t.Z*cn.Z)
// 				}
// 				pos.X = (t.X*av.X + t.Y*bv.X + t.Z*cv.X)
// 				pos.Y = (t.X*av.Y + t.Y*bv.Y + t.Z*cv.Y)
// 				pos.Z = (t.X*av.Z + t.Y*bv.Z + t.Z*cv.Z)
// 				colors = getColor(*mtl.kd, *mtl.ks, *mtl.ka, &pos,
// 					&n, interpolateUV(x, y, ap, bp, cp, au, bu, cu, t, w1, w2, w3, mtl.TextMap), 1)
// 				RW.RLock()
// 				if zbuffer[x*h+y] >= z {
// 					RW.RUnlock()
// 					RW.Lock()
// 					image.Set(int(x), int(y), colors)
// 					zbuffer[x*h+y] = z
// 					RW.Unlock()
// 				} else {
// 					RW.RUnlock()
// 				}

// 			}
// 		}
// 	}
// }

// //uv 三线性插值

// func (r *Render) FragmentShader(s, e int, zbuffer, W []float64, wg *sync.WaitGroup, rw *sync.RWMutex, imag *image.RGBA, mtl *MtlData, viewPos []PointV) {

// 	for i := s; i < e; i++ {
// 		r.Face[i].BlinPhong(zbuffer, W, r.Points, r.TextColor, r.Normal, rw, imag, mtl, viewPos)
// 	}
// 	fmt.Printf("this process is %v-%v\n", s, e)
// 	wg.Done()
// }
