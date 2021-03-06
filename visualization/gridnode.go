package visualization

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/parkerqueen/astar-visualization/astar"
)

//WidgetRenderer for gridNode
type gridNodeRenderer struct {
	label   *canvas.Text
	color   color.Color
	objects []fyne.CanvasObject

	gnode *gridNode
}

//A function controls the drawing, positioning and the general Layout of our gridNode
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

//A function that returns the minimum size of our gridNode
func (gnr *gridNodeRenderer) MinSize() fyne.Size {
	labelSize := gnr.label.MinSize()
	return labelSize.Add(gnr.padding())
}

//A function that updates the GUI corresponding to the changes
//in the corresponding gridNode
func (gnr *gridNodeRenderer) Refresh() {
	gnr.applyTheme()
	gnr.color = gnr.gnode.color
	gnr.label.Text = gnr.gnode.label
	gnr.Layout(gnr.gnode.Size())
	canvas.Refresh(gnr.gnode)
}

//A function that returns the background color of our gridNode
func (gnr *gridNodeRenderer) BackgroundColor() color.Color {
	return gnr.color
}

//A function that returns the objects held by the renderer
func (gnr *gridNodeRenderer) Objects() []fyne.CanvasObject {
	return gnr.objects
}

//A dummy function
func (gnr *gridNodeRenderer) Destroy() {}

//A function to calculate the padding of one gridNode
func (gnr *gridNodeRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*1, theme.Padding()*1)
}

//A function that applies the current theme to the girdNode
func (gnr *gridNodeRenderer) applyTheme() {
	gnr.label.TextSize = theme.TextSize()
	gnr.label.Color = theme.TextColor()
}

//Alias for a callback function
type callback func(astar.Node, *desktop.MouseEvent)

//Widget for gridNode
type gridNode struct {
	widget.BaseWidget

	node  astar.Node
	label string
	color color.Color

	onMouseIn   callback
	onMouseDown callback
}

//Implementing the Hoverable interface
func (gn *gridNode) MouseOut()                         {}
func (gn *gridNode) MouseMoved(ev *desktop.MouseEvent) {}
func (gn *gridNode) MouseIn(ev *desktop.MouseEvent) {
	gn.onMouseIn(gn.node, ev)
}

//Implementing the Mouseable interface
func (gn *gridNode) MouseUp(ev *desktop.MouseEvent) {}
func (gn *gridNode) MouseDown(ev *desktop.MouseEvent) {
	gn.onMouseDown(gn.node, ev)
}

//A function to create a rendere for our widget
func (gn *gridNode) CreateRenderer() fyne.WidgetRenderer {
	label := canvas.NewText(gn.label, theme.TextColor())
	label.Alignment = fyne.TextAlignCenter
	objects := []fyne.CanvasObject{label}
	return &gridNodeRenderer{label: label, color: gn.color, objects: objects, gnode: gn}
}

//A function to set the label and handle the GUI update
func (gn *gridNode) setLabel(label string) {
	gn.label = label
	gn.Refresh()
}

//A function to set the color and handle the GUI update
func (gn *gridNode) setColor(color color.Color) {
	gn.color = color
	gn.Refresh()
}

//A gridNode must only be created using this function
func newGridNode(node astar.Node, label string, color color.Color,
	onMouseDown callback, onMouseIn callback) *gridNode {

	gn := &gridNode{node: node, label: label, color: color, onMouseDown: onMouseDown,
		onMouseIn: onMouseIn}
	gn.ExtendBaseWidget(gn)
	return gn
}
