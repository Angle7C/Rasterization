package core

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	m, err := NewModel("E:/MyPro/GO/3Dobj/model/cube/", "cube.obj")

	t.Log(err)
	t.Log(m.ambientxMap)
	t.Log(m.normalMap)
	t.Log(m.specularMap)
	t.Log(m.diffuseMap)
	t.Log(m.heightMap)

}
