package painter

import (
	"image"
	"image/color"

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
	Figures []image.Point
}

// WhiteFill повертає Operation, що фарбує фон у білий колір.
func WhiteFill(st *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		st.BgColor = color.White
	})
}

// GreenFill повертає Operation, що фарбує фон у зелений колір.
func GreenFill(st *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		st.BgColor = color.RGBA{G: 0xff, A: 0xff}
	})
}

// BgRect повертає Operation, що малює чорний прямокутник поверх фону.
func BgRect(st *State, x1, y1, x2, y2 int) Operation {
	return OperationFunc(func(t screen.Texture) {
		r := image.Rect(x1, y1, x2, y2)
		st.BgRect = &r
	})
}

// AddFigure повертає Operation, що малює нову фігуру хреста.
func AddFigure(st *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		st.Figures = append(st.Figures, image.Pt(x, y))
	})
}

// MoveFigures повертає Operation, що переміщує всі фігури у нові координати.
func MoveFigures(st *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		for i := range st.Figures {
			st.Figures[i] = image.Pt(x, y)
		}
	})
}

// Reset повертає Operation, що очищує стан малюнку, залишаючи лише фон чорного кольору.
func Reset(st *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		st.BgColor = color.Black
		st.BgRect = nil
		st.Figures = nil
	})
}
