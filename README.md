# HexEmpireMap

Procedurally generated maps for Hex Empire. This program will generate the same map that can be found in game as a static image.

When the program starts up, it will display a randomnly generated map. If you want to change the map, you can click on the button "Random Map" below, which will generate a different map and display it on the screen.

A Hex Empire map contains 4 capital cities, one at each corner, which are displayed with a unique color (red, violet, blue, green). There are a list of cities and ports spread out across the map that the players have to control in order to increase morale.

Screenshot from map generator:

<div style="display:inline-block;">
<img src="https://raw.githubusercontent.com/samuelyuan/HexEmpireMap/master/screenshots/generated.png" alt="generated map" width="400" height="300" />
</div>

Original game for reference to show the map in the context of the game:

<div style="display:inline-block;">
<img src="https://raw.githubusercontent.com/samuelyuan/HexEmpireMap/master/screenshots/original.png" alt="original game" width="400" height="300" />
</div>

## Unit Tests

```
go test -v ./...
```

## Algorithm

When the map generators runs, it first assigns most of the tiles to be water tiles and only a small portion are assigned to be land tiles. The algorithm reduces the number of water tiles by checking each water tile and setting it to land if there is at least one out of the six neighbors that have land tiles.

Once the land and water tiles are fixed, the algorithm partitions the map into separate landmasses and for each landmass, a random tile is assigned to be a city if it doesn't border any water tiles or another city tile. After the cities are generated and added to a list, the list is shuffled. For each adjacent pair of cities in the list, a path is calculated avoiding water tiles and if it's impossible to generate a path between two cities without crossing water or the path is too long, another path will be generated that allows using water tiles and any tiles on that path which border the water will be assigned as ports.
