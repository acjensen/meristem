package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"math"
	"math/cmplx"
	"strconv"
)

type Rules map[string]string

type Mutator interface {
	Mutate(s string) string
}

func PlantFractal() Rules {
	r := make(Rules)
	r["X"] = "F-[[X]+X]+F[+FX]-X"
	r["F"] = "FF"
	return r
}

func (r Rules) Mutate(s string) string {
	var s_new string
	for _, c := range s {
		if _, has_rule := r[string(c)]; has_rule {
			s_new += r[string(c)]
		} else {
			s_new += string(c)
		}
	}
	return s_new
}

func Simulate(s string, num_steps int, m Mutator) string {
	for i := 0; i < num_steps; i++ {
		s = m.Mutate(s)
	}
	return s
}

func NextCoord(x int, y int, a int) (int, int) {
	return 1, 1
}

type Branch struct {
	phase float64
	xy    complex128
}

func Forward(b Branch, length float64) Branch {
	var b_new Branch
	b_new.phase = b.phase
	b_new.xy = b.xy + cmplx.Rect(length, b.phase)
	return b_new
}

func Turn(b Branch, delta_phase float64) Branch {
	var b_new Branch
	b_new.phase = b.phase + delta_phase
	// Keep phase within [-Pi, Pi]
	if b_new.phase > math.Pi {
		b_new.phase = b_new.phase - 2*math.Pi
	}
	if b_new.phase < math.Pi {
		b_new.phase = b_new.phase + 2*math.Pi
	}
	b_new.xy = b.xy
	return b_new
}

func Render(s string, Branch_width float64, phase_diff float64, Branch_length float64, phase_init float64) {
	init_size := 500.0
	d := gg.NewContext(int(init_size), int(init_size))
	d.SetRGB(20, 100, 30)
	d.SetLineWidth(Branch_width)
	var stack []Branch
	// |b| pointer always points to the top of the stack.
	var b *Branch
	var b_new Branch
	stack = append(stack, Branch{phase_init, complex(init_size/2.0, init_size/2.0)})
	b = &stack[len(stack)-1]

	for idx, c := range s {
		if c == 'F' {
			b_new = Forward(*b, Branch_length)
			d.DrawLine(real(b.xy), imag(b.xy), real(b_new.xy), imag(b_new.xy))
			d.Stroke()
			*b = b_new
			d.SavePNG("img/result" + strconv.Itoa(idx) + ".png")
		} else if c == 'G' {
			b_new = Forward(*b, Branch_length)
			*b = b_new
		} else if c == '-' {
			*b = Turn(*b, -1*phase_diff)
		} else if c == '+' {
			*b = Turn(*b, phase_diff)
		} else if c == '[' {
			b_save := *b
			stack = append(stack, b_save)
			b = &stack[len(stack)-1]
		} else if c == ']' {
			stack = stack[:len(stack)-1]
			b = &stack[len(stack)-1]
		}
	}
}

func main() {
	r := PlantFractal()
	ans := Simulate("X", 3, r)
	Render(ans, 1, math.Pi*2*25/365, 15, 0)
	fmt.Println("done")
}
