package visualization

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"github.com/parkerqueen/a-star-go/astar"
)

const gridRows = 30
const gridCols = 30

type painter struct {
	app fyne.App
	win fyne.Window

	grid         *astar.Grid
	walls        map[astar.Node]bool
	source, dest astar.Node
}

func (artist *painter) Paint(node astar.Node, opened bool, closed bool) {
	index := int(node.R*artist.grid.Cols + node.C)

	container := artist.win.Content().(*fyne.Container)
	gridNode := container.Objects[index].(*canvas.Rectangle)

	var color color.RGBA
	if closed {
		color = closedCellCol
	} else if opened {
		color = openedCellCol
	}

	gridNode.FillColor = color
	gridNode.Refresh()
	time.Sleep(10 * time.Millisecond)
}

func (artist *painter) init() {
	artist.grid = &astar.Grid{Rows: gridRows, Cols: gridCols, Artist: artist}

	artist.app = app.New()
	artist.win = artist.app.NewWindow("A* Path Visualization")
	artist.win.Resize(fyne.Size{Width: 1200, Height: 800})
	artist.win.SetFixedSize(true)
	artist.win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeySpace {
			artist.grid.AStarSearch(artist.source, artist.dest)
		}
	})

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(gridCols))
	for r := uint(0); r < gridRows; r++ {
		for c := uint(0); c < gridCols; c++ {
			container.AddObject(newGridNode("", simpleCellCol))
		}
	}

	artist.win.SetContent(container)
	artist.win.Show()
	artist.app.Run()
}

//AStarVisualization performs the A Star path finding animation
//on a Grid provided by the astar package
func AStarVisualization() {
	artist := painter{}
	artist.init()
}
