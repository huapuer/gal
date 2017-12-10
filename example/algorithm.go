package ga_composer

import (
	"gal"
	"math/rand"
	"bytes"
	"encoding/gob"
	"sort"
	"math"
	"time"
	"syscall"
	"fmt"
)

type Algorithm struct{
	gal.DefaultAlgorithm

	Pitches []int
	Frames []int

	MutateRate float64
	LengthCodeF int
	LengthCodeG int

	Sequence Items
	CodeF []float64
	CodeG []float64
	Result []float64
	LengthResult int
	LimitCodeF float64
	LimitCodeG float64
}

type Point struct{
	Type int64
	Point int64
	Value float64
}

func Gcd(x, y int) int {
	tmp := x % y
	if tmp > 0 {
		return Gcd(y, tmp)
	} else {
		return y
	}
}

type Item struct{
	Pitch int
	Frame int
	Value float64
	Touched bool
}

type Items []Item

func (is Items) Len() int{
	return len(is)
}

func (is Items) Less(i, j int) bool{
	return is[i].Value > is[j].Value
}

func (is Items) Swap(i, j int){
	is[i], is[j] = is[j], is[i]
}

func (a *Algorithm) Init(){
	a.DefaultAlgorithm.Init()

	lenSequence:=len(a.Pitches)*len(a.Frames)
	a.Sequence=make([]Item, 0, lenSequence)
	d:=1.05946309
	for i:=0;i<len(a.Pitches);i++ {
		for j := 0; j < len(a.Frames); j++ {
			a.Sequence = append(a.Sequence, Item{Pitch:a.Pitches[i], Frame:a.Frames[j], Value: math.Pow(d, float64(a.Pitches[i])) * float64(a.Frames[j])})
		}
	}
	sort.Sort(a.Sequence)

	a.LimitCodeF = a.Sequence[0].Value
	a.LimitCodeG = a.Sequence[0].Value

	a.LengthResult=a.LengthCodeF*a.LengthCodeG/Gcd(a.LengthCodeF, a.LengthCodeG)

	gob.Register(&Algorithm{})
}

func (a *Algorithm) GeneratePopulation(){
	for i:=0;i<a.PopulationSize;i++{
		algorithm:=a.Duplicate().(*Algorithm)
		algorithm.CodeF=make([]float64,a.LengthCodeF,a.LengthCodeF)
		for j:=0;j<a.LengthCodeF;j++{
			algorithm.CodeF[j]=rand.Float64()*a.LimitCodeF
		}
		algorithm.CodeG=make([]float64,a.LengthCodeG,a.LengthCodeG)
		for j:=0;j<a.LengthCodeG;j++{
			algorithm.CodeG[j]=rand.Float64()*a.LimitCodeG
		}
		/*
		algorithm:=&Algorithm{
			PunishRange:a.PunishRange,
			MutateRate:a.MutateRate,
			LengthCodeF:a.LengthCodeF,
			LengthCodeG:a.LengthCodeG,
			LimitCodeF:a.LimitCodeF,
			LimitCodeG:a.LimitCodeG,
			LengthResult:a.LengthResult,
			CodeF:CodeF,
			CodeG:CodeG,
			Result:make([]float64, a.LengthResult, a.LengthResult)}
		*/
		a.AppendCurrentGeneration(algorithm)
		//fmt.Printf("LENGTH OF ALGORITHM'S POPULATION: %d\n", algorithm.CurrentGeneration.Len())
	}
}

func (a *Algorithm) SelectChangePoint() interface{}{
	t:=rand.Int63n(2)
	switch(t){
	case 0:
		p:=rand.Int63n(int64(a.LengthCodeF))
		return Point{Type:t, Point:p, Value:a.CodeF[p]}
	case 1:
		p:=rand.Int63n(int64(a.LengthCodeG))
		return Point{Type:t, Point:p, Value:a.CodeG[p]}
	}
	panic("SelectChangePoint Failed.")
}

func (a *Algorithm) Duplicate() gal.IAlgorithm{
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(a); err != nil {
		panic(err.Error())
	}
	ret:=Algorithm{}
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&ret); err!=nil{
		panic(err.Error())
	}
	return &ret
}

func (a *Algorithm) Cross(point interface{}) bool{
	if p, ok := point.(Point); ok{
		switch(p.Type){
		case 0:
			a.CodeF[p.Point]=p.Value
		case 1:
			a.CodeG[p.Point]=p.Value
		}
	}else{
		panic("Point type assertion Failed.")
	}
	return true
}

func (a *Algorithm) Mutate() bool{
	if rand.Float64()>a.MutateRate {
		point := a.SelectChangePoint()
		if p, ok := point.(Point); ok {
			switch(p.Type){
			case 0:
				p.Value = rand.Float64() * a.LimitCodeF
			case 1:
				p.Value = rand.Float64() * a.LimitCodeG
			}
			a.Cross(p)
		} else {
			panic("Point type assertion Failed.")
		}
		return true
	}else{
		//println("NOT MUTATE.")
	}
	return false
}

func (a *Algorithm) Act(){
	a.Result=make([]float64, a.LengthResult, a.LengthResult)
	for i:=0;i<a.LengthResult;i++{
		a.Result[i]=a.CodeG[i%a.LengthCodeG]*a.CodeF[i%a.LengthCodeF]
	}
	/*
	for i:=0;i<a.LengthResult;i++{
		for j:=0;j<a.LengthCodeG;j++{
			a.Result[(i+j)%a.LengthResult]+=a.CodeG[j]*a.CodeF[(i+j)%a.LengthCodeF]
		}
	}
	*/
	/*
	for i:=0;i<a.LengthResult-a.LengthCodeG;i++{
		for j:=0;j<a.LengthCodeG;j++{
			a.Result[i+j]+=a.CodeG[j]*a.CodeF[(i+j)%a.LengthCodeF]
		}
	}
	*/
}

func sds(code []float64, limit float64) float64{
	sum:=0.0
	for _,v:=range code{
		sum+=v
	}
	avg:=sum/float64(len(code))
	sd:=0.0
	for _,v:=range code{
		sd+=(v-avg)*(v-avg)
	}
	sd=math.Sqrt(sd/float64(len(code)))
	return sd/limit
}

func (a *Algorithm) Fitness(){
	for _,v:=range a.Sequence{
		v.Touched=false
	}

	a.FitnessScore = 0
	tSum:=0.0
	for i:=0;i<a.LengthResult;i++{
		if a.Result[i] <= a.Sequence[0].Value && a.Result[i] >= a.Sequence[len(a.Sequence)-1].Value{
			for j := 0; j < len(a.Sequence) - 1; j++ {
				if a.Sequence[j].Value >= a.Result[i] && a.Sequence[j+1].Value <= a.Result[i]{
					d:=math.Min(a.Sequence[j].Value-a.Result[i], a.Result[i]-a.Sequence[j+1].Value)
					a.FitnessScore+=1.0-d/(a.Sequence[j].Value-a.Sequence[j+1].Value)
					if d==a.Sequence[j].Value-a.Result[i]{
						if !a.Sequence[j].Touched{
							tSum+=1.0
							a.Sequence[j].Touched=true
						}
					}else if d==a.Result[i]-a.Sequence[j+1].Value{
						if !a.Sequence[j+1].Touched{
							tSum+=1.0
							a.Sequence[j+1].Touched=true
						}
					}
				}
			}
		}
	}
	a.FitnessScore/=float64(a.LengthResult)
	a.FitnessScore=0.5*a.FitnessScore+0.5*tSum/float64(len(a.Sequence))
	if a.FitnessScore < 0{
		panic("Fitness Score less than 0.")
	}
}

func (a *Algorithm) OutputResult(){
	for i:=0;i<a.LengthResult;i++{
		if a.Result[i]>a.Sequence[0].Value{
			fmt.Printf("PITCH:%d FRAME:%d\n", a.Sequence[0].Pitch, a.Sequence[0].Frame)
		}else {
			for j := 0; j < len(a.Sequence) - 1; j++ {
				if a.Sequence[j].Value > a.Result[i] && a.Sequence[j+1].Value < a.Result[i]{
					if a.Sequence[j].Value-a.Result[i] < a.Result[i]-a.Sequence[j+1].Value{
						fmt.Printf("PITCH:%d FRAME:%d\n", a.Sequence[j].Pitch, a.Sequence[j].Frame)
					}else{
						fmt.Printf("PITCH:%d FRAME:%d\n", a.Sequence[j+1].Pitch, a.Sequence[j+1].Frame)
					}
				}
			}
		}
	}

	println("code F:")
	for i:=0;i<len(a.CodeF);i++{
		println(a.CodeF[i])
	}
	println("code G:")
	for i:=0;i<len(a.CodeG);i++{
		println(a.CodeG[i])
	}

	midi := syscall.NewLazyDLL("Midi.dll")
	openMIDI := midi.NewProc("openMIDI")
	openNote := midi.NewProc("openNote")
	closeNote := midi.NewProc("closeNote")
	closeMIDI := midi.NewProc("closeMIDI")

	openMIDI.Call()
	ins:=[]int{0,0x00003C90,0x00004090,0x00004390}

	pf:=100

	for i:=0;i<a.LengthResult;i++{
		if a.Result[i]>a.Sequence[0].Value{
			openNote.Call(uintptr(ins[a.Sequence[0].Pitch]))
			time.Sleep(time.Duration(a.Sequence[0].Frame*pf)*time.Millisecond)
			closeNote.Call(uintptr(ins[a.Sequence[0].Pitch]))
		}else {
			for j := 0; j < len(a.Sequence) - 1; j++ {
				if a.Sequence[j].Value > a.Result[i] && a.Sequence[j+1].Value < a.Result[i]{
					if a.Sequence[j].Value-a.Result[i] < a.Result[i]-a.Sequence[j+1].Value{
						openNote.Call(uintptr(ins[a.Sequence[j].Pitch]))
						time.Sleep(time.Duration(a.Sequence[j].Frame*pf)*time.Millisecond)
						closeNote.Call(uintptr(ins[a.Sequence[j].Pitch]))
					}else{
						openNote.Call(uintptr(ins[a.Sequence[j+1].Pitch]))
						time.Sleep(time.Duration(a.Sequence[j+1].Frame*pf)*time.Millisecond)
						closeNote.Call(uintptr(ins[a.Sequence[j+1].Pitch]))
					}
				}
			}
		}
	}

	closeMIDI.Call()
}

func (a* Algorithm) Play(code Items){
	midi := syscall.NewLazyDLL("Midi.dll")
	openMIDI := midi.NewProc("openMIDI")
	openNote := midi.NewProc("openNote")
	closeNote := midi.NewProc("closeNote")
	closeMIDI := midi.NewProc("closeMIDI")

	openMIDI.Call()

	pf:=100

	for _, v:=range code{
		if v.Pitch>=0 {
			openNote.Call(uintptr(0x00003C90))
		}
		time.Sleep(time.Duration(v.Frame*pf)*time.Millisecond)
		if v.Pitch>=0 {
			closeNote.Call(uintptr(0x00003C90))
		}
	}

	closeMIDI.Call()
}