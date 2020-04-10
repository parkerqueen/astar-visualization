package visualization

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"github.com/parkerqueen/a-star-go/astar"
)

type painter struct {
	app          fyne.App
	win          fyne.Window
	grid         *astar.Grid
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
}

//AStarVisualization performs the A Star path finding animation
//on a Grid provided by the astar package
func AStarVisualization(grid *astar.Grid, source, dest astar.Node) {
	artist := painter{}
	artist.grid = grid
	artist.source = source
	artist.dest = dest

	artist.app = app.New()
	artist.win = artist.app.NewWindow("A* Path Visualization")
	artist.win.Resize(fyne.Size{Width: 1200, Height: 800})
	artist.win.SetFixedSize(true)
	artist.win.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeySpace {
			grid.AStarSearch(source, dest)
		}
	})

	grid.Artist = &artist

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(int(grid.Cols)))
	for r := uint(0); r < grid.Rows; r++ {
		for c := uint(0); c < grid.Rows; c++ {
			node := astar.Node{R: r, C: c}
			cellCol := simpleCellCol

			if source == node {
				cellCol = sourceCellCol
			} else if dest == node {
				cellCol = destCellCol
			} else if grid.Walls[node] {
				cellCol = blockedCellCol
			}

			container.AddObject(canvas.NewRectangle(cellCol))
		}
	}

	artist.win.SetContent(container)
	artist.win.Show()
	artist.app.Run()
}
