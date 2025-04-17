package painter

import (
	"image"
	"image/color"

	"github.com/ProMKQ/kpi-lab3/ui"
	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

type State struct {
	BgColor color.Color
	BgRect  *image.Rectangle
	Shapes  []image.Point
}

func Repaint(state *State, t screen.Texture) {
	tRect := t.Bounds()
	t.Fill(tRect, state.BgColor, screen.Src)
	if state.BgRect != nil {
		t.Fill(*state.BgRect, color.Black, screen.Src)
	}

	for i := range state.Shapes {
		ui.DrawShape(t, tRect, state.Shapes[i])
	}
}

// WhiteFill повертає Operation, що фарбує фон у білий колір.
func WhiteFill(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BgColor = color.White
		Repaint(state, t)
	})
}

// GreenFill повертає Operation, що фарбує фон у зелений колір.
func GreenFill(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BgColor = color.RGBA{G: 0xff, A: 0xff}
		Repaint(state, t)
	})
}

// BgRect повертає Operation, що малює чорний прямокутник поверх фону.
func BgRect(state *State, x1, y1, x2, y2 int) Operation {
	return OperationFunc(func(t screen.Texture) {
		r := image.Rect(x1, y1, x2, y2)
		state.BgRect = &r
		Repaint(state, t)
	})
}

// AddShape повертає Operation, що малює нову фігуру хреста.
func AddShape(state *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.Shapes = append(state.Shapes, image.Pt(x, y))
		ui.DrawShape(t, t.Bounds(), image.Pt(x, y))
	})
}

// MoveShapes повертає Operation, що переміщує всі фігури у нові координати.
func MoveShapes(state *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		for i := range state.Shapes {
			state.Shapes[i] = image.Pt(x, y)
		}
		Repaint(state, t)
	})
}

// Reset повертає Operation, що очищує стан малюнку, залишаючи лише фон чорного кольору.
func Reset(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BgColor = color.Black
		state.BgRect = nil
		state.Shapes = nil
		t.Fill(t.Bounds(), state.BgColor, screen.Src)
	})
}
