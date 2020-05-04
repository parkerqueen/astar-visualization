package visualization

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/parkerqueen/astar-visualization/astar"
)

//The number of rows & cols to show in the visualization
const gridRows = 30
const gridCols = 30

//A type to store the current status of the visualization
type vstatus int

//Go-style enums for vstatus
const (
	walling vstatus = iota
	choosingSource
	choosingDestination
	running
	finished
)

//Help text to display for the current visualization state
const (
	wallingLabel             = "Click & drag to draw obstacles/walls."
	choosingSourceLabel      = "Click on any node to mark it as the source."
	choosingDestinationLabel = "Click on any node to mark it as the destination."
	runningLabel             = "Sit back & enjoy the animation."
	finishedLabel            = "Reset the board to run the visualization again."
)

//A struct to run & manage the visualization GUI
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

//The only function to be called on a visualization object from outside
//This function intializes and launches the GUI
func (vis *visualization) init() {
	vis.app = app.New()
	vis.win = vis.app.NewWindow("A* Path Visualization")
	vis.win.Resize(fyne.Size{Width: 1200, Height: 800})
	vis.win.SetFixedSize(true)
	vis.win.SetContent(vis.setup())
	vis.win.Show()
	vis.app.Run()
}

//A function used to set up an entirely new window and the required
//data structures for the visualization
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
			time.Sleep(10 * time.Millisecond)
		}

		vis.runAction.Hide()
		vis.sourceAction.Hide()
		vis.destinationAction.Hide()
		vis.resetAction.Show()
		vis.status = finished
		vis.label.SetText(finishedLabel)
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

//A function which is called by the astar algorithm everytime a astar.Node's state is changed
//This function essentially updates the appearance (color & label) of the corresponding gridNode
//according to the new state of astar.Node
func (vis *visualization) Paint(node astar.Node, cost uint, opened bool, closed bool) {
	if node == vis.source || node == vis.destination {
		return
	}

	gridNode := vis.getGridNode(node)
	gridNode.setLabel(strconv.FormatUint(uint64(cost), 10))
	var color color.RGBA
	if closed {
		color = closedNodeCol
	} else if opened {
		color = openedNodeCol
	}
	gridNode.setColor(color)
	time.Sleep(10 * time.Millisecond)
}

//A callback which is fired whenever the mouse pointer enters a gridNode
//This callback handles the toggling of 'walled' status for each node
func (vis *visualization) onMouseInCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == walling && ev.Button == desktop.LeftMouseButton &&
		node != vis.source && node != vis.destination {
		vis.toggleWalled(node)
	}
}

//A function which handles a click event for any gridNode and marks it either as the source,
//the destination or changes it's 'walled' status
func (vis *visualization) onMouseDownCB(node astar.Node, ev *desktop.MouseEvent) {
	if vis.status == walling && node != vis.source && node != vis.destination {
		vis.toggleWalled(node)
	} else if vis.status == choosingSource && node != vis.destination {
		vis.setSource(node)
		vis.status = walling
		vis.label.SetText(wallingLabel)
	} else if vis.status == choosingDestination && node != vis.source {
		vis.setDestination(node)
		vis.status = walling
		vis.label.SetText(wallingLabel)
	}
}

//A helper function used by the onMouseInCB to toggle the walled status
func (vis *visualization) toggleWalled(node astar.Node) {
	isWalled := vis.grid.Walls[node]
	vis.grid.Walls[node] = !isWalled

	if !isWalled {
		vis.getGridNode(node).setColor(walledNodeCol)
	} else {
		vis.getGridNode(node).setColor(simpleNodeCol)
	}
}

//A helper function to set any node as the source node
func (vis *visualization) setSource(node astar.Node) {
	if vis.sourceSet() {
		vis.getGridNode(vis.source).setColor(simpleNodeCol)
	}

	vis.source = node
	vis.getGridNode(node).setColor(sourceNodeCol)
}

//A helper function to fset any node as the destination node
func (vis *visualization) setDestination(node astar.Node) {
	if vis.destinationSet() {
		vis.getGridNode(vis.destination).setColor(simpleNodeCol)
	}

	vis.destination = node
	vis.getGridNode(node).setColor(destinationNodeCol)
}

//A helper function to create a new gridNode with the provided color
func (vis *visualization) newGridNode(node astar.Node, color color.RGBA) *gridNode {
	return newGridNode(node, "", color, vis.onMouseDownCB, vis.onMouseInCB)
}

//A helper function to grab the gridNode of any astar.Node
func (vis *visualization) getGridNode(node astar.Node) *gridNode {
	index := int(node.R*vis.grid.Cols + node.C)
	return vis.nodesContainer.Objects[index].(*gridNode)
}

//A helper function to check if a source has been set
func (vis *visualization) sourceSet() bool {
	return vis.source.R != gridRows
}

//A helper function to check if a destination has been set
func (vis *visualization) destinationSet() bool {
	return vis.destination.R != gridRows
}

//AStarVisualization create a visualization object and calls init on it
func AStarVisualization() {
	vis := visualization{}
	vis.init()
}
