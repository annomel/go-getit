package main

import (
	"flag"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	page "github.com/annomel/go-getit/pages"
	"github.com/annomel/go-getit/pages/about"
	"github.com/annomel/go-getit/pages/home"
	"github.com/annomel/go-getit/tools"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	flag.Parse()
	go func() {

		w := app.NewWindow(app.Title("Go Get It"))

		go tools.Listen(func() { w.Invalidate() })

		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	router := page.NewRouter()
	router.Register(0, home.New(&router))

	router.Register(5, about.New(&router))

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				router.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}
