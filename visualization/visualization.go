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

type vstatus int

const (
	walling vstatus = iota
	choosingSource
	choosingDestination
	running
)

type visualization struct {
	app              fyne.App
	win              fyne.Window
	nodesContainer   *fyne.Container
	actionsContainer *fyne.Container

	grid        *astar.Grid
	source      astar.Node
	destination astar.Node

	status vstatus
}

func (vis *visualization) init() {
	vis.source = astar.Node{R: gridRows, C: gridCols}
	vis.destination = astar.Node{R: gridRows, C: gridCols}
	vis.grid = &astar.Grid{Rows: gridRows, Cols: gridCols, Walls: make(map[astar.Node]bool),
		Artist: vis}

	vis.app = app.New()
	vis.win = vis.app.NewWindow("A* Path Visualization")
	vis.win.Resize(fyne.Size{Width: 1200, Height: 800})
	vis.win.SetFixedSize(true)

	vis.win.SetContent(vis.initContainer())
	vis.win.Show()
	vis.app.Run()
}

func (vis *visualization) initContainer() *fyne.Container {
	vis.nodesContainer = fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for r := uint(0); r < gridRows; r++ {
		for c := uint(0); c < gridCols; c++ {
			vis.nodesContainer.AddObject(vis.newGridNode(astar.Node{R: r, C: c}, simpleNodeCol))
		}
	}

	runAction := widget.NewButton("RUN", func() {
		if !vis.sourceSet() || !vis.destinationSet() {
			return
		}

		vis.status = running
		path := vis.grid.AStarSearch(vis.source, vis.destination)
		for _, node := range path {
			vis.getGridNode(node).setColor(pathedNodeCol)
		}
	})

	sourceAction := widget.NewButton("CHOOSE SOURCE", func() {
		vis.status = choosingSource
	})

	destinationAction := widget.NewButton("CHOOSE DESTINATION", func() {
		vis.status = choosingDestination
	})

	vis.actionsContainer = fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(),
		sourceAction, destinationAction, runAction, layout.NewSpacer())

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), vis.nodesContainer, vis.actionsContainer)
}

func (vis *visualization) newGridNode(node astar.Node, color color.RGBA) *gridNode {
	return newGridNode(node, "", color, vis.onMouseDownCB, vis.onMouseInCB)
}

func (vis *visualization) Paint(node astar.Node, opened bool, closed bool) {
	if node == vis.source || node == vis.destination {
		return
	}
	gridNode := vis.getGridNode(node)

	var color color.RGBA
	if closed {
		color = closedNodeCol
	} else if opened {
		color = openedNodeCol
	}

	gridNode.setColor(color)
	time.Sleep(10 * time.Millisecond)
}

func (vis *visualization) getGridNode(node astar.Node) *gridNode {
	index := int(node.R*vis.grid.Cols + node.C)
	return vis.nodesContainer.Objects[index].(*gridNode)
}

func (vis *visualization) sourceSet() bool {
	return vis.source.R != gridRows
}

func (vis *visualization) destinationSet() bool {
	return vis.destination.R != gridRows
}

func (vis *visualization) onMouseInCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == walling && ev.Button == desktop.LeftMouseButton &&
		node != vis.source && node != vis.destination {

		isWalled := vis.grid.Walls[node]
		vis.grid.Walls[node] = !isWalled

		if !isWalled {
			vis.getGridNode(node).setColor(walledNodeCol)
		} else {
			vis.getGridNode(node).setColor(simpleNodeCol)
		}
	}
}

func (vis *visualization) onMouseDownCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == choosingSource {
		if vis.sourceSet() {
			vis.getGridNode(vis.source).setColor(simpleNodeCol)
		}

		vis.source = node
		vis.getGridNode(node).setColor(sourceNodeCol)
	} else if vis.status == choosingDestination {
		if vis.destinationSet() {
			vis.getGridNode(vis.destination).setColor(simpleNodeCol)
		}

		vis.destination = node
		vis.getGridNode(node).setColor(destinationNodeCol)
	}
	vis.status = walling
}

//AStarVisualization begins the visualization GUI for the A* Pathfinding
func AStarVisualization() {
	vis := visualization{}
	vis.init()
}
