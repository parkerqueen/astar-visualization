package visualization

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type gridNodeRenderer struct {
	label   *canvas.Text
	color   color.Color
	objects []fyne.CanvasObject

	gnode *gridNode
}

func (gnr *gridNodeRenderer) Layout(size fyne.Size) {
	padding := gnr.padding()
	contentSize := size.Subtract(padding)
	contentOffset := fyne.NewPos(padding.Width/2, padding.Height/2)

	labelSize := gnr.label.MinSize()
	labelWidth := labelSize.Width
	labelHeight := labelSize.Height
	labelOffset := fyne.NewPos(
		(contentSize.Width+labelWidth)/2-labelSize.Width,
		(contentSize.Height-labelHeight)/2)

	gnr.label.Resize(labelSize)
	gnr.label.Move(contentOffset.Add(labelOffset))
}

func (gnr *gridNodeRenderer) MinSize() fyne.Size {
	labelSize := gnr.label.MinSize()
	return labelSize.Add(gnr.padding())
}

func (gnr *gridNodeRenderer) Refresh() {
	gnr.applyTheme()
	gnr.label.Text = gnr.gnode.Label
	gnr.Layout(gnr.gnode.Size())
	canvas.Refresh(gnr.gnode)
}

func (gnr *gridNodeRenderer) BackgroundColor() color.Color {
	return gnr.color
}

func (gnr *gridNodeRenderer) Objects() []fyne.CanvasObject {
	return gnr.objects
}

func (gnr *gridNodeRenderer) Destroy() {}

func (gnr *gridNodeRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

func (gnr *gridNodeRenderer) applyTheme() {
	gnr.label.TextSize = theme.TextSize()
	gnr.label.Color = theme.TextColor()
}

type gridNode struct {
	widget.BaseWidget
	Label string
	Color color.Color
}

func (gn *gridNode) setLabel(label string) {
	gn.Label = label
	gn.Refresh()
}

func (gn *gridNode) setColor(color color.Color) {
	gn.Color = color
	gn.Refresh()
}

func (gn *gridNode) CreateRenderer() fyne.WidgetRenderer {
	label := canvas.NewText(gn.Label, theme.TextColor())
	label.Alignment = fyne.TextAlignCenter
	objects := []fyne.CanvasObject{label}
	return &gridNodeRenderer{label: label, color: gn.Color, objects: objects, gnode: gn}
}

func newGridNode(label string, color color.Color) *gridNode {
	gn := &gridNode{Label: label, Color: color}
	gn.ExtendBaseWidget(gn)
	return gn
}
