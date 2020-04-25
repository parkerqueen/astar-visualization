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
	finished
)

const (
	wallingLabel             = "Click & drag to draw obstacles/walls."
	choosingSourceLabel      = "Click on any node to mark it as the source."
	choosingDestinationLabel = "Click on any node to mark it as the destination."
	runningLabel             = "Sit back & enjoy the animation."
	finishedLabel            = "Reset the board to run the visualization again."
)

type visualization struct {
	app fyne.App
	win fyne.Window

	nodesContainer   *fyne.Container
	actionsContainer *fyne.Container

	runAction         *widget.Button
	sourceAction      *widget.Button
	destinationAction *widget.Button
	resetAction       *widget.Button
	label             *widget.Label

	grid        *astar.Grid
	source      astar.Node
	destination astar.Node

	status vstatus
}

func (vis *visualization) init() {
	vis.app = app.New()
	vis.win = vis.app.NewWindow("A* Path Visualization")
	vis.win.Resize(fyne.Size{Width: 1200, Height: 800})
	vis.win.SetFixedSize(true)
	vis.win.SetContent(vis.setup())
	vis.win.Show()
	vis.app.Run()
}

func (vis *visualization) setup() *fyne.Container {
	vis.source = astar.Node{R: gridRows, C: gridCols}
	vis.destination = astar.Node{R: gridRows, C: gridCols}
	vis.grid = &astar.Grid{Rows: gridRows, Cols: gridCols, Walls: make(map[astar.Node]bool),
		Artist: vis}

	vis.nodesContainer = fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for r := uint(0); r < gridRows; r++ {
		for c := uint(0); c < gridCols; c++ {
			vis.nodesContainer.AddObject(vis.newGridNode(astar.Node{R: r, C: c}, simpleNodeCol))
		}
	}

	vis.runAction = widget.NewButton("RUN", func() {
		if !vis.sourceSet() || !vis.destinationSet() ||
			vis.status == running || vis.status == finished {
			return
		}

		vis.status = running
		vis.label.SetText(runningLabel)
		path := vis.grid.AStarSearch(vis.source, vis.destination)
		for _, node := range path {
			vis.getGridNode(node).setColor(pathedNodeCol)
		}
		vis.status = finished
		vis.label.SetText(finishedLabel)
		vis.showResetAction()
	})

	vis.sourceAction = widget.NewButton("CHOOSE SOURCE", func() {
		if vis.status != running && vis.status != finished {
			vis.status = choosingSource
			vis.label.SetText(choosingSourceLabel)
		}
	})

	vis.destinationAction = widget.NewButton("CHOOSE DESTINATION", func() {
		if vis.status != running && vis.status != finished {
			vis.status = choosingDestination
			vis.label.SetText(choosingDestinationLabel)
		}
	})

	vis.resetAction = widget.NewButton("RESET", func() {
		vis.win.SetContent(vis.setup())
	})
	vis.resetAction.Hide()

	vis.actionsContainer = fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(),
		vis.sourceAction, vis.destinationAction, vis.runAction, vis.resetAction, layout.NewSpacer())

	vis.status = walling
	vis.label = widget.NewLabel(wallingLabel)

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), vis.nodesContainer, vis.actionsContainer,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), vis.label, layout.NewSpacer()))
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

func (vis *visualization) onMouseInCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == walling && ev.Button == desktop.LeftMouseButton &&
		node != vis.source && node != vis.destination {
		vis.toggleWalled(node)
	}
}

func (vis *visualization) onMouseDownCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == walling && node != vis.source && node != vis.destination {
		vis.toggleWalled(node)
	} else if vis.status == choosingSource {
		vis.setSource(node)
		vis.status = walling
		vis.label.SetText(wallingLabel)
	} else if vis.status == choosingDestination {
		vis.setDestination(node)
		vis.status = walling
		vis.label.SetText(wallingLabel)
	}
}

func (vis *visualization) showResetAction() {
	vis.runAction.Hide()
	vis.sourceAction.Hide()
	vis.destinationAction.Hide()
	vis.resetAction.Show()
}

func (vis *visualization) toggleWalled(node astar.Node) {
	isWalled := vis.grid.Walls[node]
	vis.grid.Walls[node] = !isWalled

	if !isWalled {
		vis.getGridNode(node).setColor(walledNodeCol)
	} else {
		vis.getGridNode(node).setColor(simpleNodeCol)
	}
}

func (vis *visualization) setSource(node astar.Node) {
	if vis.sourceSet() {
		vis.getGridNode(vis.source).setColor(simpleNodeCol)
	}

	vis.source = node
	vis.getGridNode(node).setColor(sourceNodeCol)
}

func (vis *visualization) setDestination(node astar.Node) {
	if vis.destinationSet() {
		vis.getGridNode(vis.destination).setColor(simpleNodeCol)
	}

	vis.destination = node
	vis.getGridNode(node).setColor(destinationNodeCol)
}

func (vis *visualization) newGridNode(node astar.Node, color color.RGBA) *gridNode {
	return newGridNode(node, "", color, vis.onMouseDownCB, vis.onMouseInCB)
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

//AStarVisualization begins the visualization GUI for the A* Pathfinding
func AStarVisualization() {
	vis := visualization{}
	vis.init()
}
