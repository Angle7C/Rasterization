package core

import (
	"Matrix/matrix"
	"Matrix/vector"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

//光结构
type Light struct {
	Pos    []vector.Vector3D
	Indent []vector.Vector3D
	EyePos vector.Vector3D
}

//空间结构
type Space struct {
	models      []Model
	zbuffer     *DepthMessage
	viewMat     *matrix.Matrix4
	perspectMat *matrix.Matrix4
	viewPort    *matrix.Matrix4
	light       Light
}

func NewLight(eye *vector.Vector3D) *Light {
	light := &Light{}
	light.EyePos = *eye
	return light
}
func NewSpace(light Light) *Space {
	world := &Space{}
	world.light = light
	world.viewMat = matrix.NewMatrix4()
	world.perspectMat = matrix.NewMatrix4()
	world.viewPort = matrix.NewMatrix4()
	return world
}
func (light *Light) AddLight(pos, ide vector.Vector3D) {
	light.Pos = append(light.Pos, pos)
	light.Indent = append(light.Indent, ide)
}

//向世界空间添加模型
func (world *Space) AddModel(models []Model) {
	world.models = append(world.models, models...)

	for k1, v1 := range world.models {
		world.models[k1].w = make([]float64, len(v1.point))
		for k2, v2 := range v1.point {
			vd := *world.perspectMat.Mul(world.viewMat.Mul(v1.modelMat)).MulVector3D(&v2)
			world.models[k1].w[k2] = 1.0 / vd.W
			vd = *vd.Mul_lamda(world.models[k1].w[k2])
			vd = *world.viewPort.MulVector3D(&vd)
			world.models[k1].point[k2] = vd
		}
	}
}

//设置viewPort
func (world *Space) MakeViewPort(w, h float64) {
	world.viewPort.Set(
		w/2, 0, 0, w/2,
		0, h/2, 0, h/2,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	world.zbuffer = NewZbuffer(int(w), int(h))
}

//设置Perspective 投影矩阵
func (world *Space) MakePerspective(fov, aspect, near, far float64) {
	world.perspectMat = world.perspectMat.MakePerspective(fov, aspect, near, far)
}

//设置viewMat 视图矩阵
func (world *Space) MakeView(eyePos, target *vector.Vector3D) {
	world.viewMat = world.viewMat.LookAt(eyePos, target, vector.NewVector3D(0, 1, 0))
}
func (world *Space) Render(out string, w, h int) {
	var (
		wg *sync.WaitGroup = &sync.WaitGroup{}
	)
	f, _ := os.OpenFile(out, os.O_CREATE, 0777)
	n := image.NewRGBA(image.Rect(0, 0, w, h))
	world.zbuffer = NewZbuffer(w, h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			world.zbuffer.Colors[x*h+y] = color.RGBA{0, 0, 0, 255}
		}
	}

	fmt.Printf("world.zbuffer: %p\n", world.zbuffer)
	wg.Add(len(world.models))
	for k, _ := range world.models {
		go world.models[k].Render(world.zbuffer, wg, n, h, world.light, world.perspectMat, world.viewPort)

	}
	wg.Wait()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {

			n.SetRGBA(x, y, world.zbuffer.Colors[x*h+y])

		}
	}
	png.Encode(f, n)
}
func (zbuffer *DepthMessage) judgeZbuffer(z float64, x, y, h int) bool {
	zbuffer.rw.Lock()
	defer zbuffer.rw.Unlock()
	return zbuffer.Depth[x*h+y] > z
}
func (zbuffer *DepthMessage) setZbuffer(x, y, h int, z float64, rgb color.RGBA) {
	zbuffer.rw.Lock()
	defer zbuffer.rw.Unlock()
	zbuffer.Depth[x*h+y] = z
	zbuffer.Colors[x*h+y] = rgb
}
