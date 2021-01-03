package main

// Simulate plant growth using a discrete Lindenmayer function.
//
// Example:
//    go run meristem.go
//    convert img/result*.png -loop 0 animation.gif

import (
	"fmt"
	"github.com/fogleman/gg"
	"math"
	"math/cmplx"
	"strconv"
)

// An L-system is defined as a set of |Rules|. A Rule maps one set of symbols to another.
type Rules map[string]string

// Define rules for a plant-like system.
func PlantFractal() Rules {
	r := make(Rules)
	r["X"] = "F-[[X]+X]+F[+FX]-X"
	r["F"] = "FF"
	return r
}

// Compute the next system state according to Rules and the current state.
func (r Rules) Mutate(state string) string {
	var new_state string
	// Loop through symbols in the |state| string. If there is a rule for the symbol, apply it. Else, let the symbol pass through.
	for _, c := range state {
		if _, has_rule := r[string(c)]; has_rule {
			new_state += r[string(c)]
		} else {
			new_state += string(c)
		}
	}
	return new_state
}

type Mutator interface {
	Mutate(s string) string
}

// Mutate the state exactly |num_steps| times.
func Simulate(state string, num_steps int, m Mutator) string {
	for i := 0; i < num_steps; i++ {
		state = m.Mutate(state)
	}
	return state
}

// Each |Branch| has a |phase| (an angle) and an absolute location |xy| represented in the complex plane.
type Branch struct {
	phase float64
	xy    complex128
}

// Calculate new |xy| location using current phase and given distance.
func Forward(b Branch, distance float64) Branch {
	var b_new Branch
	b_new.phase = b.phase
	b_new.xy = b.xy + cmplx.Rect(distance, b.phase)
	return b_new
}

// Turn the current branch by |delta_phase|.
func Turn(b Branch, delta_phase float64) Branch {
	var b_new Branch
	b_new.phase = b.phase + delta_phase
	// Keep phase within [-Pi, Pi].
	if b_new.phase > math.Pi {
		b_new.phase = b_new.phase - 2*math.Pi
	}
	if b_new.phase < math.Pi {
		b_new.phase = b_new.phase + 2*math.Pi
	}
	b_new.xy = b.xy
	return b_new
}

// Interpret symbols in |s| as instructions for a 2D vector drawing. Save each
// subsequent instruction as an image.
func Render(state string, branch_length float64, branch_width float64, phase_diff float64, phase_init float64, canvas_size int, img_path string) {
	d := gg.NewContext(canvas_size, canvas_size)
	d.SetRGB(205, 133, 63)
	d.SetLineWidth(branch_width)
	// Create a stack of branch-off points |b_stack|. |b| always points to the top of the stack.
	var b_stack []Branch
	var b *Branch
	var b_new Branch
	// Initialize the stack.
	b_stack = append(b_stack, Branch{phase_init, complex(float64(canvas_size)/2, float64(canvas_size))})
	b = &b_stack[len(b_stack)-1]

	// Interpret each character of the input string as a drawing action.
	for idx, c := range state {
		if c == 'F' {
			// Move forward from |b| to |b_new| and draw a line.
			b_new = Forward(*b, branch_length)
			d.DrawLine(real(b.xy), imag(b.xy), real(b_new.xy), imag(b_new.xy))
			d.Stroke()
			*b = b_new
			d.SavePNG(img_path + "/" + strconv.Itoa(idx) + ".png")
		} else if c == 'G' {
			// Move forward without drawing anything.
			b_new = Forward(*b, branch_length)
			*b = b_new
		} else if c == '-' {
			// Turn left.
			*b = Turn(*b, -1*phase_diff)
		} else if c == '+' {
			// Turn right.
			*b = Turn(*b, phase_diff)
		} else if c == '[' {
			// Store the current branch state.
			b_save := *b
			b_stack = append(b_stack, b_save)
			b = &b_stack[len(b_stack)-1]
		} else if c == ']' {
			// Return to the last saved branch state.
			b_stack = b_stack[:len(b_stack)-1]
			b = &b_stack[len(b_stack)-1]
		}
	}
	fmt.Println("Saved images to: " + img_path)
}

func main() {
	rule := PlantFractal()
	final_state := Simulate("X", 4, rule)
	Render(final_state, 15, 1, math.Pi*2*25/365, -math.Pi/2, 500, "img")
}
