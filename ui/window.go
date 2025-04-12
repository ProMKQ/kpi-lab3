package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	scr  screen.Screen
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	pw.scr = s // Зберігаємо screen.Screen для створення текстури

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event:
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button == mouse.ButtonRight && e.Direction == mouse.DirPress {
			center := image.Pt(int(e.X), int(e.Y))
			tex, err := pw.scr.NewTexture(pw.sz.Bounds().Size())
			if err != nil {
				log.Println("failed to create texture:", err)
				return
			}
			drawCross(tex, pw.sz, center)
			pw.Update(tex)
		}

	case paint.Event:
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src)

	center := image.Point{
		X: pw.sz.Bounds().Dx() / 2,
		Y: pw.sz.Bounds().Dy() / 2,
	}

	drawCross(pw.w, pw.sz, center)

	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func drawCross(target interface {
	Fill(r image.Rectangle, c color.Color, op draw.Op)
}, sz size.Event, center image.Point) {
	crossSize := sz.Bounds().Dx() / 2
	thickness := crossSize / 3
	blue := color.RGBA{0, 0, 255, 255}

	vRect := image.Rect(
		center.X-thickness/2,
		center.Y-crossSize/2,
		center.X+thickness/2,
		center.Y+crossSize/2,
	)

	hRect := image.Rect(
		center.X-crossSize/2,
		center.Y-thickness/2,
		center.X+crossSize/2,
		center.Y+thickness/2,
	)

	target.Fill(vRect, blue, draw.Src)
	target.Fill(hRect, blue, draw.Src)
}
