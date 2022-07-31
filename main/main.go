package main

import (
	"core"
	"fmt"
)

func main() {
	// var (
	// 	w   int = 400
	// 	h   int = 400
	// 	err error
	// 	f   *os.File
	// 	m   *core.ModelPly
	// 	i   *image.RGBA
	// 	r   *core.Render
	// )
	// m, err = core.NewModelPly("..\\model\\Triangle.ply")
	// if err != nil {
	// 	log.Fatalf("new model fail")
	// 	return
	// }
	// r, err = core.NewRenderPly(m)
	// if err != nil {
	// 	log.Fatalf("new Render fail")
	// 	return

	// }
	// m = nil
	// i = image.NewRGBA(image.Rect(0, 0, w, h))
	// r.MakeModelMat(10, 10, 10)
	// r.MakeViewMat(
	// 	vector.NewVector3D(20, 20, 20),
	// 	vector.NewVector3D(0, 0, 0),
	// 	vector.NewVector3D(0, 1, 0))
	// r.MakePerspectMat(90, 1, 10, 100)
	// r.MakeViewPort(float64(i.Rect.Dx()), float64(i.Rect.Dy()))
	// r.MVP()
	// r.ToScreen()
	// f, err = os.Create("..\\model\\modl.png")
	// if err != nil {
	// 	log.Fatalf("png Encode fail")
	// 	return
	// }
	// r.Render(i)
	// png.Encode(f, i)
	m, mm := core.Parse("../model/cube/", "cube.obj")
	r, _ := core.NewRenderObj(m)
	fmt.Printf("r: %v\n", r)
	fmt.Printf("mm: %v\n", mm)
}
