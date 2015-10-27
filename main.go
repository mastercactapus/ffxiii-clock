package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "strconv"

type Face struct {
	Len int
	Numbers []int
	Active []bool
	Solution []int
	solve chan []int
}

func NewFace(n []int) *Face {
	f := new(Face)
	f.Len = len(n)
	f.Numbers = make([]int,f.Len)
	f.Active = make([]bool,f.Len)
	for i,v:=range n {
		f.Numbers[i]=v
		f.Active[i] = true
	}
	f.Solution = make([]int,0,f.Len)
	f.solve = make(chan []int,100)
	return f
}

func (f *Face) Copy() *Face {
	nf:=new(Face)
	nf.Numbers = make([]int,f.Len)
	nf.Active = make([]bool,f.Len)
	nf.Solution = make([]int,0,f.Len)
	nf.Len = f.Len
	for i :=0;i<f.Len;i++ {
		nf.Numbers[i] = f.Numbers[i]
		nf.Active[i] = f.Active[i]
		nf.solve = f.solve
	}
	for _,v:=range f.Solution {
		nf.Solution = append(nf.Solution, v)
	}
	return nf	
}

func (f *Face) CheckSolve() bool {
	for _,a:=range f.Active {
		if a {
			return false
		}
	}
	return true
}

func (f *Face) TakeStep(n int) {
	if n >= f.Len {
		return
	}
	if !f.Active[n] {
		return
	}
	f.Active[n] = false
	f.Solution = append(f.Solution, n)
	if f.CheckSolve() {
		f.solve<-f.Solution
		return
	}
	hand1:= n + f.Numbers[n]
	hand2:=n - f.Numbers[n]
	if hand2 < 0 {
		hand2+=f.Len
	}
	if hand1 >= f.Len {
		hand1-=f.Len
	}
	go f.Copy().TakeStep(hand1)
	go f.Copy().TakeStep(hand2)
}

func (f *Face) Solve() []int {
	for i:=0;i<f.Len;i++ {
		go f.Copy().TakeStep(i)
	}
	return <-f.solve
}

func main() {
	fmt.Println("Enter the face numbers starting with the top, blank after the last")
	n:=make([]int, 0, 100)
	r:=bufio.NewReader(os.Stdin)
	i:=0
	for {
		i++
		fmt.Printf("Position #%d: ", i)
		text,err:=r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if strings.TrimSpace(text) == "" {
			break
		}
		val,err:=strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			panic(err)
		}
		n=append(n,val)
	}

	fmt.Println("")
	f:=NewFace(n)
	s:=f.Solve()
	for i,v:=range s {
		fmt.Printf("Step #%d -- Position #%d\n",i+1,v+1)
	}
}
