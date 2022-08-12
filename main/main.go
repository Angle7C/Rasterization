package main

import (
	"Matrix/vector"
	"core"
)

func main() {
	var (
		w     int        = 1024
		h     int        = 1024
		light core.Light = *core.NewLight(vector.NewVector3D(0, 200, 300))
		space core.Space
	)
	light.AddLight(*vector.NewVector3D(200, 200, 200), *vector.NewVector3D(500, 500, 500))
	// light.AddLight(*vector.NewVector3D(0, 200, 00), *vector.NewVector3D(1000, 1000, 1000))
	space = *core.NewSpace(light)
	space.MakeView(&light.EyePos, vector.NewVector3D(100, 100, 0))
	space.MakePerspective(90, 1, -1, -1000)
	space.MakeViewPort(float64(w), float64(h))
	a, err := core.NewModel("../model/cube/", "cube.obj")
	if err != nil {
		panic(err)
	}
	a.MakeScales(50, 50, 50)
	a.MakeTranslation(0, 0, -50)
	b, err := core.NewModel("E:/MyPro/C/tinyrenderer/obj/", "floor.obj")
	if err != nil {
		panic(err)
	}
	b.MakeScales(100, 100, 100)
	b.MakeRotation(0, 0, 0)
	b.MakeTranslation(0, 0, -100)
	space.AddModel([]core.Model{*b, *a})
	space.Render("render.png", w, h)

}
