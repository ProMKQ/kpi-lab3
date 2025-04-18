package painter_test

import (
	"bytes"
	"github.com/ProMKQ/kpi-lab3/painter"
	"image"
	"image/color"
	"testing"

	"github.com/ProMKQ/kpi-lab3/painter/lang"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
)

type mockReceiver struct {
	lastColor color.Color
}

func (r *mockReceiver) Update(t screen.Texture) {
	if mt, ok := t.(*mockTexture); ok && len(mt.Colors) > 0 {
		r.lastColor = mt.Colors[len(mt.Colors)-1]
	}
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}
func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}
func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}
func (m *mockTexture) Size() image.Point {
	return image.Pt(400, 400)
}
func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}

func TestBgColorCommands(t *testing.T) {
	var state painter.State
	parser := lang.NewParser(&state)

	ops, err := parser.Parse(bytes.NewBufferString("green\nwhite\nupdate"))
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	rc := &mockReceiver{}
	var loop painter.Loop
	loop.Receiver = rc

	loop.Start(mockScreen{})
	loop.Post(painter.OperationList(ops))
	loop.StopAndWait()

	if state.BgColor != color.White {
		t.Errorf("expected BgColor white, got %v", state.BgColor)
	}
}

func TestBgRectCommand(t *testing.T) {
	var state painter.State
	parser := lang.NewParser(&state)

	ops, err := parser.Parse(bytes.NewBufferString("bgrect 0.25 0.25 0.75 0.75\nupdate"))
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	rc := &mockReceiver{}
	var loop painter.Loop
	loop.Receiver = rc

	loop.Start(mockScreen{})
	loop.Post(painter.OperationList(ops))
	loop.StopAndWait()

	if state.BgRect == nil || state.BgRect.Min.X != 100 || state.BgRect.Min.Y != 100 {
		t.Errorf("unexpected BgRect: %+v", state.BgRect)
	}
}

func TestShapeCommands(t *testing.T) {
	var state painter.State
	parser := lang.NewParser(&state)

	ops, err := parser.Parse(bytes.NewBufferString("figure 0.0 0.0\nmove 1.0 1.0\nupdate"))
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	rc := &mockReceiver{}
	var loop painter.Loop
	loop.Receiver = rc

	loop.Start(mockScreen{})
	loop.Post(painter.OperationList(ops))
	loop.StopAndWait()

	want := image.Pt(400, 400)
	if len(state.Shapes) != 1 || state.Shapes[0] != want {
		t.Errorf("expected moved shape to %v, got %v", want, state.Shapes)
	}
}

func TestResetCommand(t *testing.T) {
	var state painter.State
	parser := lang.NewParser(&state)

	ops, err := parser.Parse(bytes.NewBufferString("white\nfigure 0.5 0.5\nbgrect 0.1 0.1 0.2 0.2\nreset\nupdate"))
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	rc := &mockReceiver{}
	var loop painter.Loop
	loop.Receiver = rc

	loop.Start(mockScreen{})
	loop.Post(painter.OperationList(ops))
	loop.StopAndWait()

	if state.BgColor != color.Black {
		t.Errorf("expected BgColor black after reset, got %v", state.BgColor)
	}
	if state.BgRect != nil {
		t.Errorf("expected BgRect nil after reset, got %v", state.BgRect)
	}
	if len(state.Shapes) != 0 {
		t.Errorf("expected Shapes empty after reset, got %v", state.Shapes)
	}
}

func TestInvalidCommands(t *testing.T) {
	cmds := []string{
		"figure 0.1",
		"move a b",
		"cmd",
	}

	for _, cmd := range cmds {
		t.Run("invalid: "+cmd, func(t *testing.T) {
			var state painter.State
			parser := lang.NewParser(&state)
			_, err := parser.Parse(bytes.NewBufferString(cmd + "\n"))
			if err == nil {
				t.Errorf("expected error for cmd: %q", cmd)
			}
		})
	}
}
