package home

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/annomel/go-getit/icon"
	page "github.com/annomel/go-getit/pages"

	"github.com/annomel/go-getit/tools"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the AppBar component.
type Page struct {
	activitylist   widget.List
	addressInput   component.TextField
	contextBtn     widget.Clickable
	inputAlignment text.Alignment

	widget.List
	*page.Router
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Downloads",
		Icon: icon.HomeIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	p.activitylist.Axis = layout.Vertical
	p.List.ScrollToEnd = true
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Start,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if p.contextBtn.Clicked() && p.addressInput.Editor.Len() > 0 {
					tools.Get <- p.addressInput.Editor.Text()
					p.addressInput.Editor.SetText("")
				}

				return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					return layout.Flex{Axis: layout.Vertical}.Layout(
						gtx,

						layout.Rigid(func(gtx C) D {
							return material.List(th, &p.activitylist).Layout(gtx, len(tools.Activites), func(gtx C, index int) D {

								return layout.UniformInset(5).Layout(gtx, func(gtx C) D {

									sl := material.Slider(th, &widget.Float{Value: tools.Activites[index].Progress}, 0, 1)

									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(material.H5(th, tools.Activites[index].Address).Layout),
										layout.Rigid(sl.Layout),
										layout.Rigid(material.Body1(th, tools.Activites[index].ResponseCode).Layout),
										layout.Rigid(material.Caption(th, tools.Activites[index].Status).Layout),
									)
								})

							})

						}),

						layout.Rigid(material.Body1(th, fmt.Sprint("Total:", len(tools.Activites))).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {

							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D {
									p.addressInput.Alignment = p.inputAlignment
									return p.addressInput.Layout(gtx, th, "Address")
								}),
								layout.Rigid(material.Button(th, &p.contextBtn, "Get").Layout))
						}),
					)

				})

			}),
		)
	})
}
