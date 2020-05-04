<p align="center"><img width="100px" src="https://i.imgur.com/aKsxqTq.png" alt="A* Visualization"/></p>

---

`astar-visualization` is just a simple tool to visualize the A* Pathfinding Algorithm and I built it in an attempt to learn the ropes of Golang. Below is a gif of the tool in action: (**2x Speed**)

<p align="center"><img src="https://imgur.com/DQSkFvs.gif" alt="A* Visualization"/></p>

## Instructions
* Click & drag around to create walls/obstacles.
* Click the `CHOOSE SOURCE` or `CHOOSE DESTINATION` button and click on any cell to mark them as the source/destination.
* Click `RUN` to run the visualization. If you run without choosing a source and destination, the visualization will not run.
* Click `RESET` at the end of the visualization to reset the board.

**The visualization considers a diagonal step as a valid step. Keep that in mind when drawing obstacles as it may appear that the algorithm is bypassing the walls/obstacles.**

## To Run (Ubuntu)
`astar-visualization` is built with `fyne`, a cross-platform GUI library and hence can run on `Linux`, `Windows` & `MacOS`. However, I have only tested it on Ubuntu 19.10 and will provide instructions for the same.

Make sure you have `gcc` & `golang`(Version 1.11+) installed. Then run the following.
```
sudo apt install libgl1-mesa-dev xorg-dev
git clone git@github.com:parkerqueen/astar-visualization.git
cd astar-visualization
go run main.go
```
If you wish to run on `Windows` or `MacOS` or any other flavor of `Linux`, you can visit [fyne's getting started](https://fyne.io/develop/index) to install all the prerequisites for it. Once done, you can clone the repository and use `go run main.go` to run the tool.

## Issues
* Final help text of resetting the board is sometimes not displayed correctly.
