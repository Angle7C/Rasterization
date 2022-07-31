package core

import (
	"Matrix/vector"
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var num = 0

type ModelPly struct {
	Points       int
	Faces        []Face
	propertyNum  int
	propertyName map[int]string
	propertyTo   map[string]int
	propertyType map[int]int
	PointNum     int
	FaceNum      int
	Point        []Point
}
type Face struct {
	PointIndex  []int64
	TextCoords  []int64
	NormalIndex []int64
}
type Point struct {
	property []propertyer
}

func (f *Face) setIndex(str string) (err error) {
	var (
		cap int64
		// err error
		s []string
	)
	s = strings.Split(strings.TrimSpace(str), " ")
	cap, err = strconv.ParseInt(s[0], 0, 0)
	if err != nil {
		log.Panicln(err)
		return
	}
	f.PointIndex = make([]int64, cap)
	f.TextCoords = make([]int64, cap)
	f.NormalIndex = make([]int64, cap)

	for i := 1; i < len(s); i++ {
		f.PointIndex[i-1], err = strconv.ParseInt(s[i], 0, 0)

		if err != nil {
			log.Panicln(err)
			return
		}
	}
	return

}
func NewPoint(cap int) *Point {
	p := &Point{}
	p.property = make([]propertyer, cap)
	return p
}
func (p *Point) AddProperty(t, index int, name, v string) {
	var (
		temp propertyer
	)
	if t == 1 {
		temp = &propertyFloat{}
	} else {
		temp = &propertyInt{}
	}
	temp.SetProperty(name, v)
	p.property[index] = temp
}

type propertyInt struct {
	Value int64
	Name  string
}
type propertyFloat struct {
	Value float64
	Name  string
}
type propertyer interface {
	SetProperty(name, v string)
	GetValue() float64
	SetValue(value float64)
}
type Pointer interface {
	GetVector(map[string]int) (vector *vector.Vector3D, err error)
	GetColor(map[string]int) (vector *vector.Vector3D, err error)
	GetArray(para ...string) (arr []float64)
	SetVector(m map[string]int, v *vector.Vector3D)
}

func (p *Point) GetColor(m map[string]int) (v *vector.Vector3D, err error) {
	var arr []float64
	arr, err = p.GetArray(m, "r", "g", "b")
	if err != nil {
		return
	}
	v = vector.NewVector3D(arr[0], arr[1], arr[2])
	err = nil
	return
}
func (p *Point) GetVector(m map[string]int) (v *vector.Vector3D, err error) {
	var arr []float64
	arr, err = p.GetArray(m, "x", "y", "z")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	v = vector.NewVector3D(arr[0], arr[1], arr[2])
	v.W = 1
	err = nil
	return
}
func (p *Point) SetVector(m map[string]int, v *vector.Vector3D) {
	i, j, k := m["x"], m["y"], m["z"]
	p.property[i].SetValue(v.X)
	p.property[j].SetValue(v.Y)
	p.property[k].SetValue(v.Z)
}
func (p *Point) GetArray(m map[string]int, para ...string) (arr []float64, err error) {
	arr = make([]float64, len(para))
	for k, v := range para {
		if m[v] == 0 {
			arr = nil
			err = errors.New("没有这个属性")
			return
		}
		arr[k] = p.property[m[v]-1].GetValue()
	}
	err = nil
	return
}
func (p *propertyInt) SetProperty(name, v string) {
	p.Value, _ = strconv.ParseInt(v, 0, 0)
	p.Name = name
}
func (p *propertyInt) GetValue() float64 {
	return float64(p.Value)
}
func (p *propertyInt) SetValue(v float64) {
	p.Value = int64(v)
}
func (p *propertyFloat) SetProperty(name, v string) {
	p.Value, _ = strconv.ParseFloat(v, 64)
	p.Name = name
}
func (p *propertyFloat) GetValue() float64 {
	return p.Value
}
func (p *propertyFloat) SetValue(v float64) {
	p.Value = v
}
func NewModelPly(path string) (modelPly *ModelPly, err error) {
	var (
		f *os.File
		r *bufio.Scanner
	)
	modelPly = &ModelPly{}
	f, err = os.OpenFile(path, os.O_RDONLY, fs.FileMode(0777))
	if err != nil {
		log.Fatalln("read ply error")
		return modelPly, errors.New("read ply error")

	}
	defer f.Close()
	r = bufio.NewScanner(f)
	err = modelPly.setHeader(r)
	if err != nil {
		log.Printf("set Head Fail:%v", err.Error())
	}
	err = modelPly.setPoint(r)
	if err != nil {
		log.Printf("set Point Fail:%v", err.Error())
	}
	err = modelPly.setFace(r)
	if err != nil {
		log.Printf("set Face Fail:%v", err.Error())
	}
	err = nil
	return
}
func (m *ModelPly) setHeader(r *bufio.Scanner) (err error) {
	var (
		// strs []string
		reg        *regexp.Regexp
		headString string = ""
		temp       string = ""
	)
	for r.Scan() {
		temp = r.Text()
		if temp == "ply" {
			continue
		}
		if temp == "end_header" {
			break
		}
		headString += temp + "\n"
	}
	reg, err = regexp.Compile("element (vertex|face) [0-9]*")
	if err != nil {
		log.Fatalln("read head element V|F fail:", err)
		return
	}
	s := reg.FindAllString(headString, 2)
	fmt.Printf("%s\n%s\n", s[0], s[1])
	//读取顶点数量
	m.PointNum, err = strconv.Atoi(strings.Split(s[0], " ")[2])
	if err != nil {
		log.Fatalln("vertexNum error:", err)
		return
	}
	//读取面数量
	m.FaceNum, err = strconv.Atoi(strings.Split(s[1], " ")[2])
	if err != nil {
		log.Fatalln("faceNum error:", err)
		return
	}
	//读取顶点属性数量
	reg, err = regexp.Compile("property (int|float|double|uchar)* [a-z0-9]*")
	if err != nil {
		log.Fatalln("property read error:", err)
		return
	}
	s = reg.FindAllString(headString, -1)
	m.propertyName = make(map[int]string)
	m.propertyType = make(map[int]int)
	m.propertyTo = make(map[string]int)
	m.propertyNum = len(s)
	for i, v := range s {
		str := strings.Split(v, " ")
		m.propertyName[i] = str[2]
		m.propertyTo[str[2]] = i + 1
		if str[1] == "float" {
			m.propertyType[i] = 1
		} else {
			m.propertyType[i] = 2
		}
	}
	err = nil
	return
}
func (m *ModelPly) setPoint(r *bufio.Scanner) (err error) {
	var (
		p      *Point
		str    string
		vertex []string
	)
	m.Point = make([]Point, m.PointNum)
	for i := 0; i < m.PointNum && r.Scan(); i++ {
		str = r.Text()
		p = NewPoint(m.propertyNum)
		vertex = strings.Split(strings.TrimSpace(str), " ")
		for k, v := range vertex {
			p.AddProperty(m.propertyType[k], k, m.propertyName[k], v)
		}
		m.Point[i] = *p
	}
	err = nil
	return
}
func (m *ModelPly) setFace(r *bufio.Scanner) (err error) {
	var (
		f   *Face
		str string
	)
	m.Faces = make([]Face, m.FaceNum)
	for i := 0; i < m.FaceNum && r.Scan(); i++ {
		str = r.Text()
		f = &Face{}
		err = f.setIndex(str)
		if err != nil {
			return
		}
		m.Faces[i] = *f
	}
	err = nil
	return
}

func (m *ModelPly) arrayVector() (vrr []vector.Vector3D, err error) {
	for i := 0; i < m.PointNum; i++ {
		vd, err := m.Point[i].GetVector(m.propertyTo)
		if err != nil {
			return nil, err
		}
		vrr = append(vrr, *vd)
	}
	err = nil
	return
}
func (m *ModelPly) vectorArray(vrr []vector.Vector3D) {}
