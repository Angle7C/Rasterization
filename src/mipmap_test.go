package core_test

import (
	"core"
	"image/png"
	"os"
	"strconv"
	"testing"
)

func TestNewMipMap(t *testing.T) {
	m := core.NewMipMap("..\\model\\rock\\rock.png")
	for i := 0; i < len(m.LevelImage); i++ {
		file, _ := os.OpenFile(strconv.FormatInt(int64(i), 10)+".png", os.O_RDWR|os.O_CREATE, 0766)
		png.Encode(file, m.LevelImage[i])
	}
}
