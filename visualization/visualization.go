package visualization

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/parkerqueen/a-star-go/astar"
)

const gridRows = 30
const gridCols = 30

type pstatus int

const (
	drawingWalls pstatus = iota
	choosingSource
	choosingDestination
	running
)

type painter struct {
	app            fyne.App
	win            fyne.Window
	nodesContainer *fyne.Container

	grid        *astar.Grid
	walls       map[astar.Node]bool
	source      astar.Node
	destination astar.Node

	status pstatus
}

func (artist *painter) Paint(node astar.Node, opened bool, closed bool) {
	if node == artist.source || node == artist.destination {
		return
	}

	index := int(node.R*artist.grid.Cols + node.C)
	gridNode := artist.nodesContainer.Objects[index].(*gridNode)

	var color color.RGBA
	if closed {
		color = closedNodeCol
	} else if opened {
		color = openedNodeCol
	}

	gridNode.setColor(color)
	time.Sleep(10 * time.Millisecond)
}

func (artist *painter) initContainer() *fyne.Container {
	artist.nodesContainer = fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for r := uint(0); r < gridRows; r++ {
		for c := uint(0); c < gridCols; c++ {
			artist.nodesContainer.AddObject(artist.newGridNodeWithCB(astar.Node{R: r, C: c}, simpleNodeCol))
		}
	}

	actionsController := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(),
		widget.NewButton("RUN", func() {}), layout.NewSpacer())

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), artist.nodesContainer, actionsController)
}

func (artist *painter) init() {
	artist.source = astar.Node{R: 0, C: 0}
	artist.destination = astar.Node{R: 10, C: 14}
	artist.grid = &astar.Grid{Rows: gridRows, Cols: gridCols, Artist: artist}

	artist.app = app.New()
	artist.win = artist.app.NewWindow("A* Path Visualization")
	artist.win.Resize(fyne.Size{Width: 1200, Height: 800})
	artist.win.SetFixedSize(true)

	artist.win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeySpace {
			artist.grid.AStarSearch(artist.source, artist.destination)
		}
	})

	artist.win.SetContent(artist.initContainer())
	artist.win.Show()
	artist.app.Run()
}

func (artist *painter) newGridNodeWithCB(node astar.Node, color color.RGBA) *gridNode {
	return newGridNode(node, "", color, artist.onClickCB, artist.onMouseMoveCB)
}

func (artist *painter) onClickCB(node astar.Node, ev *desktop.MouseEvent)     {}
func (artist *painter) onMouseMoveCB(node astar.Node, ev *desktop.MouseEvent) {}

//AStarVisualization performs the A Star path finding animation
//on a Grid provided by the astar package
func AStarVisualization() {
	artist := painter{}
	artist.init()
}
