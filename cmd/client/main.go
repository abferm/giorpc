// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/abferm/giorpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	disable = flag.Bool("disable", false, "disable all widgets")
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	encodeInput.SetText(encodeInitialValue)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rpcclient = giorpc.NewGiorpcClient(conn)

	go func() {
		w := new(app.Window)
		w.Option(app.Size(unit.Dp(800), unit.Dp(700)), app.Title("GIORPC"))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	events := make(chan event.Event)
	acks := make(chan struct{})

	go func() {
		for {
			ev := w.Event()
			events <- ev
			<-acks
			if _, ok := ev.(app.DestroyEvent); ok {
				return
			}
		}
	}()

	var ops op.Ops
	for {
		e := <-events
		switch e := e.(type) {
		case app.DestroyEvent:
			acks <- struct{}{}
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if *disable {
				gtx = gtx.Disabled()
			}
			client(gtx, th)
			e.Frame(gtx.Ops)
		}
		acks <- struct{}{}

	}
}

var (
	encodeInput       = new(widget.Editor)
	encodeOutput      = ""
	encodeError       = ""
	encodeButton      = new(widget.Clickable)
	radioButtonsGroup = new(widget.Enum)
	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	rpcclient giorpc.GiorpcClient
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func client(gtx layout.Context, th *material.Theme) layout.Dimensions {
	widgets := []layout.Widget{
		func(gtx C) D {
			l := material.H3(th, "Simple GRPC Client GUI")
			return l.Layout(gtx)
		},
		func(gtx C) D {
			l := material.H6(th, "Text to Encode:")
			return l.Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(200))
			e := material.Editor(th, encodeInput, "Hint")
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},
		func(gtx C) D {
			items := []layout.FlexChild{}
			for i := giorpc.Encoding_ENCODING_BASE32_STANDARD; i <= giorpc.Encoding_ENCODING_BASE64_URL_SAFE; i++ {
				name, ok := giorpc.Encoding_name[int32(i)]
				if !ok {
					continue
				}
				items = append(items, layout.Rigid(material.RadioButton(th, radioButtonsGroup, name, strings.TrimPrefix(name, "ENCODING_")).Layout))
			}
			return layout.Flex{}.Layout(gtx, items...)
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						for encodeButton.Clicked(gtx) {
							encoding, ok := giorpc.Encoding_value[radioButtonsGroup.Value]
							if !ok {
								encodeError = fmt.Sprintf("illegal encoding selection %q", radioButtonsGroup.Value)
								break
							}
							in := giorpc.EncodeRequest{
								Encoding: giorpc.Encoding(encoding),
								Decoded:  encodeInput.Text(),
							}
							resp, err := rpcclient.Encode(context.Background(), &in)
							if err != nil {
								encodeError = err.Error()
								break
							}
							encodeError = ""
							encodeOutput = resp.Encoded
						}
						items := []layout.FlexChild{layout.Rigid(material.Button(th, encodeButton, "Encode").Layout)}
						if len(encodeError) > 0 {
							e := material.Body1(th, encodeError)
							e.Color = color.NRGBA{R: 0xff, A: 0xff}
							items = append(items, layout.Rigid(e.Layout))
						}
						return layout.Flex{}.Layout(gtx, items...)
					})
				}),
			)
		},
		func(gtx C) D {
			l := material.H6(th, "Encoded Text:")
			return l.Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(200))
			e := material.Body1(th, encodeOutput)
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},
	}

	return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

const encodeInitialValue = `Sample text`
