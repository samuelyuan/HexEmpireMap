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

### Overview
The map generation algorithm begins by creating a mostly water-dominated map and incrementally carves out landmasses, cities, and ports through neighbor-based transformations and pathfinding. The process mimics organic land formation and logical placement of settlements.

### Step-by-Step Process

#### 1. Initial Tile Assignment
- Most tiles are randomly assigned as water.
- A small portion are randomly assigned as land, usually near the four corners for balance.

#### 2. Water Reduction by Neighbor Influence
- For each tile marked as water:
  - If at least one of its 6 neighbors is land, convert it into land.
- This causes land to "spread" organically from seed points.

#### 3. Landmass Partitioning
- After the land tiles are finalized, they are grouped into distinct landmasses.
  - A landmass is a connected group of adjacent land tiles.
  - Each landmass gets a unique ID.

#### 4. City Generation
- For each landmass:
  - A number of cities is determined based on its size.
  - Cities are placed only on tiles that:
    - Don’t border any water tiles.
    - Don’t border another city tile.

#### 5. City Shuffle & Pathfinding
- The city list is shuffled to ensure randomness in route connections.
- For each adjacent city pair:
  - A path is attempted avoiding water.
  - If the path is too long or impossible:
    - A secondary path allows water traversal.
    - Any land tile on this water-based path that borders water is converted into a port.

### Notes
- This results in naturally shaped continents with functional city-to-city connections and coastlines.
- The ports provide logical naval access between cities across water barriers.
