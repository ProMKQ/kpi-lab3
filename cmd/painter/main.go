package main

import (
	"net/http"

	"github.com/ProMKQ/kpi-lab3/painter"
	"github.com/ProMKQ/kpi-lab3/painter/lang"
	"github.com/ProMKQ/kpi-lab3/ui"
)

func main() {
	var (
		pv     ui.Visualizer
		opLoop painter.Loop
		state  painter.State
		parser = lang.NewParser(&state)
	)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
