package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"math"
	"math/cmplx"
)

type rules map[string]string

type mutator interface {
	mutate(s string) string
}

func koch() rules {
	r := make(rules)
	r["F"] = "F+F-F-F+F"
	return r
}

func tri() rules {
	r := make(rules)
	r["F"] = "F-G+F+G-F"
	r["G"] = "GG"
	return r
}

func plant_fractal() rules {
	r := make(rules)
	r["X"] = "F-[[X]+X]+F[+FX]-X"
	r["F"] = "FF"
	return r
}

func (r rules) mutate(s string) string {
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

func simulate(s string, num_steps int, m mutator) string {
	for i := 0; i < num_steps; i++ {
		s = m.mutate(s)
	}
	return s
}

func next_coord(x int, y int, a int) (int, int) {
	return 1, 1
}

type branch struct {
	phase float64
	xy    complex128
}

func forward(b branch, length float64) branch {
	var b_new branch
	b_new.phase = b.phase
	b_new.xy = b.xy + cmplx.Rect(length, b.phase)
	return b_new
}

func turn(b branch, delta_phase float64) branch {
	var b_new branch
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

func render(s string, branch_width float64, phase_diff float64, branch_length float64, phase_init float64) {
	init_size := 10000.0
	d := gg.NewContext(int(init_size), int(init_size))
	d.SetRGB(20, 100, 30)
	d.SetLineWidth(branch_width)
	var stack []branch
	// |b| pointer always points to the top of the stack.
	var b *branch
	var b_new branch
	stack = append(stack, branch{phase_init, complex(init_size/2.0, init_size/2.0)})
	b = &stack[len(stack)-1]

	for _, c := range s {
		if c == 'F' {
			b_new = forward(*b, branch_length)
			d.DrawLine(real(b.xy), imag(b.xy), real(b_new.xy), imag(b_new.xy))
			d.Stroke()
			// Update |b|.
			*b = b_new
		} else if c == 'G' {
			b_new = forward(*b, branch_length)
			*b = b_new
		} else if c == '-' {
			*b = turn(*b, -1*phase_diff)
		} else if c == '+' {
			*b = turn(*b, phase_diff)
		} else if c == '[' {
			b_save := *b
			stack = append(stack, b_save)
			b = &stack[len(stack)-1]
		} else if c == ']' {
			stack = stack[:len(stack)-1]
			// Update |b|.
			b = &stack[len(stack)-1]
		}
	}
	d.SavePNG("result.png")
}

func main() {
	r := plant_fractal()
	ans := simulate("X", 6, r)
	render(ans, 1, math.Pi*2*25/365, 15, 0)
	fmt.Println("done")
}
