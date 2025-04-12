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

	sz          size.Event
	crossCenter *image.Point // нове: центр хрестика
}

func (pw *Visualizer) Update(t screen.Texture) {
	panic("implement me")
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	driver.Main(pw.run)
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

	pw.scr = s

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

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e)
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		return e.To == lifecycle.StageDead
	case key.Event:
		return e.Code == key.CodeEscape
	}
	return false
}

func (pw *Visualizer) handleEvent(e any) {
	switch e := e.(type) {
	case size.Event:
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button == mouse.ButtonRight && e.Direction == mouse.DirPress {
			center := image.Pt(int(e.X), int(e.Y))
			pw.crossCenter = &center
			pw.w.Send(paint.Event{}) // тригер оновлення вікна
		}

	case paint.Event:
		pw.draw()
		pw.w.Publish()
	}
}

func (pw *Visualizer) draw() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src)

	var center image.Point
	if pw.crossCenter != nil {
		center = *pw.crossCenter
	} else {
		center = image.Point{
			X: pw.sz.Bounds().Dx() / 2,
			Y: pw.sz.Bounds().Dy() / 2,
		}
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
