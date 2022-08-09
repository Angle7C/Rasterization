package core

import (
	"Matrix/matrix"
	"Matrix/vector"
)

//光结构
type Light struct {
	Pos    []vector.Vector3D
	Indent []vector.Vector3D
	eyePos vector.Vector3D
}

//空间结构
type Space struct {
	models      []Model
	zbuffer     []DepthMessage
	viewMat     *matrix.Matrix4
	perspectMat *matrix.Matrix4
	viewPort    *matrix.Matrix4
	light       Light
}

func NewSpace(light Light) *Space {
	world := &Space{}
	world.light = light
	world.viewMat = matrix.NewMatrix4()
	world.perspectMat = matrix.NewMatrix4()
	world.viewPort = matrix.NewMatrix4()
	return world
}

//向世界空间添加模型
func (world *Space) AddModel(models []Model) {
	world.models = models
	for _, v1 := range world.models {
		for k2, v2 := range v1.point {
			vd := *world.perspectMat.Mul(world.viewMat.Mul(v1.modelMat)).MulVector3D(&v2)
			v1.w[k2] = 1.0 / vd.W
			v1.point[k2] = *vd.Mul_lamda(v1.w[k2])

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
}

//设置Perspective 投影矩阵
func (world *Space) MakePerspective(fov, aspect, near, far float64) {
	world.perspectMat = world.perspectMat.MakePerspective(fov, aspect, near, far)
}

//设置viewMat 视图矩阵
func (world *Space) MakeView(target *vector.Vector3D) {
	world.viewMat = world.viewMat.LookAt(&light.eyePos, target, vector.NewVector3D(0, 1, 0))
}
