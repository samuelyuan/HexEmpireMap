package hexempire

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	fileDir = "./images"
)

type HexMap struct {
	RandomSeed             int
	Background             *ebiten.Image
	BackgroundDirt         *ebiten.Image
	BackgroundGrass        *ebiten.Image
	BackgroundSea          *ebiten.Image
	BackgroundTiles        *ebiten.Image
	UI                     *ebiten.Image
	Board                  *Board
	Pathfinder             *Pathfinder
	TextFont               font.Face
	BackgroundImageOptions map[string]*LandImageOptions
	WaterImageOptions      map[string]*WaterImageOptions
	TownImageOptions       map[string]*LandImageOptions
}

type Board struct {
	MapNumber       int
	XMax            int
	YMax            int
	HexWidth        int
	HexHeight       int
	Fields          map[string]*Field
	LandCount       int
	LandGroups      [][]*Field
	Towns           []*Field
	PartiesCapitals []*Field
	TownNames       []string
}

type Field struct {
	FX        int
	FY        int
	X         int
	Y         int
	LandId    int
	Type      string
	Capital   int
	Neighbors [6]*Point2D
	IsLand    bool
	Estate    string
	TownName  string
}

type Point2D struct {
	X int
	Y int
}

type WaterImageOptions struct {
	WaterBgImgIndex int
	FlipH           int
	FlipV           int
	RotateDegrees   int
}

type LandImageOptions struct {
	BgDirtImgIndex  int
	BgGrassImgIndex int
	FlipH           int
	FlipV           int
	RotateDegrees   int
}

func NewHexMap(mapNumber int, textFont font.Face) *HexMap {
	hexMap := &HexMap{}
	hexMap.RandomSeed = mapNumber
	hexMap.Background = ebiten.NewImage(800, 600)
	hexMap.BackgroundDirt = ebiten.NewImage(800, 600)
	hexMap.BackgroundGrass = ebiten.NewImage(800, 600)
	hexMap.BackgroundSea = ebiten.NewImage(800, 600)
	hexMap.BackgroundTiles = ebiten.NewImage(800, 600)
	hexMap.UI = ebiten.NewImage(800, 600)
	hexMap.Board = NewBoard()
	hexMap.Board.MapNumber = mapNumber
	hexMap.Board.TownNames = generateAllTowns()
	hexMap.Pathfinder = NewPathfinder()
	hexMap.TextFont = textFont
	return hexMap
}

func NewBoard() *Board {
	return &Board{
		XMax:            20,
		YMax:            11,
		HexWidth:        50,
		HexHeight:       40,
		Fields:          make(map[string]*Field),
		LandCount:       0,
		LandGroups:      make([][]*Field, 0),
		Towns:           make([]*Field, 0),
		PartiesCapitals: make([]*Field, 4),
	}
}

func generateAllTowns() []string {
	return []string{
		"Abu Dhabi", "Abuja", "Accra", "Addis Ababa", "Algiers", "Amman", "Amsterdam", "Ankara", "Antananarivo", "Apia", "Ashgabat", "Asmara", "Astana", "Asunción", "Athens",
		"Baghdad", "Baku", "Bamako", "Bangkok", "Bangui", "Banjul", "Basseterre", "Beijing", "Beirut", "Belgrade", "Belmopan", "Berlin", "Bern", "Bishkek", "Bissau", "Bogotá",
		"Brasília", "Bratislava", "Brazzaville", "Bridgetown", "Brussels", "Bucharest", "Budapest", "Buenos Aires", "Bujumbura", "Cairo", "Canberra",
		"Cape Town", "Caracas", "Castries", "Chisinau", "Conakry", "Copenhagen", "Cotonou",
		"Dakar", "Damascus", "Dhaka", "Dili", "Djibouti", "Dodoma", "Doha", "Dublin", "Dushanbe", "Delhi",
		"Freetown", "Funafuti", "Gabarone", "Georgetown", "Guatemala City", "Hague", "Hanoi", "Harare", "Havana", "Helsinki", "Honiara", "Hong Kong",
		"Islamabad", "Jakarta", "Jerusalem", "Kabul", "Kampala", "Kathmandu", "Khartoum", "Kyiv", "Kigali", "Kingston", "Kingstown", "Kinshasa", "Kuala Lumpur", "Kuwait City",
		"La Paz", "Liberville", "Lilongwe", "Lima", "Lisbon", "Ljubljana", "Lobamba", "Lomé", "London", "Luanda", "Lusaka", "Luxembourg",
		"Madrid", "Majuro", "Malé", "Managua", "Manama", "Manila", "Maputo", "Maseru", "Mbabane", "Melekeok", "Mexico City", "Minsk", "Mogadishu", "Monaco", "Monrovia", "Montevideo", "Moroni", "Moscow", "Muscat",
		"Nairobi", "Nassau", "Naypyidaw", "N'Djamena", "New Delhi", "Niamey", "Nicosia", "Nouakchott", "Nuku'alofa", "Nuuk",
		"Oslo", "Ottawa", "Ouagadougou", "Palikir", "Panama City", "Paramaribo", "Paris", "Phnom Penh", "Podgorica", "Prague", "Praia", "Pretoria", "Pyongyang",
		"Quito", "Rabat", "Ramallah", "Reykjavík", "Riga", "Riyadh", "Rome", "Roseau",
		"San José", "San Marino", "San Salvador", "Sanaá", "Santiago", "Santo Domingo", "Sao Tomé", "Sarajevo", "Seoul", "Singapore", "Skopje", "Sofia", "South Tarawa", "St. George's", "St. John's", "Stockholm", "Sucre", "Suva",
		"Taipei", "Tallinn", "Tashkent", "Tbilisi", "Tegucigalpa", "Teheran", "Thimphu", "Tirana", "Tokyo", "Tripoli", "Tunis", "Ulaanbaatar",
		"Vaduz", "Valletta", "Victoria", "Vienna", "Vientiane", "Vilnius", "Warsaw", "Washington", "Wellington", "Windhoek", "Yamoussoukro", "Yaoundé", "Yerevan", "Zagreb", "Zielona Góra",
		"Poznań", "Wrocław", "Gdańsk", "Szczecin", "Łódź", "Białystok", "Toruń", "St. Petersburg", "Turku", "Örebro", "Chengdu",
		"Wuppertal", "Frankfurt", "Düsseldorf", "Essen", "Duisburg", "Magdeburg", "Bonn", "Brno", "Tours", "Bordeaux", "Nice", "Lyon", "Stara Zagora", "Milan", "Bologna", "Sydney", "Venice", "New York",
		"Barcelona", "Zaragoza", "Valencia", "Seville", "Graz", "Munich", "Birmingham", "Naples", "Cologne", "Turin", "Marseille", "Leeds", "Kraków", "Palermo", "Genoa",
		"Stuttgart", "Dortmund", "Rotterdam", "Glasgow", "Málaga", "Bremen", "Sheffield", "Antwerp", "Plovdiv", "Thessaloniki", "Kaunas", "Lublin", "Varna", "Ostrava", "Iaşi", "Katowice",
		"Cluj-Napoca", "Timişoara", "Constanţa", "Pskov", "Vitebsk", "Arkhangelsk", "Novosibirsk", "Samara", "Omsk", "Chelyabinsk", "Ufa", "Volgograd", "Perm", "Kharkiv", "Odessa", "Donetsk", "Dnipropetrovsk",
		"Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "Dallas", "Detroit", "Indianapolis", "San Francisco", "Atlanta", "Austin", "Vermont", "Toronto", "Montreal", "Vancouver", "Gdynia", "Edmonton",
	}
}

func (hexMap *HexMap) rand(n int) int {
	hexMap.RandomSeed = (hexMap.RandomSeed*9301 + 49297) % 233280
	return int(math.Floor(float64(hexMap.RandomSeed) / 233280 * float64(n)))
}

func (hexMap *HexMap) flipImageMatrix(options *ebiten.DrawImageOptions, img *ebiten.Image, h int, v int) {
	imgWidth, imgHeight := img.Size()
	width := float64(imgWidth)
	height := float64(imgHeight)

	if h > 0 && v > 0 {
		options.GeoM.Scale(-1, -1)
		options.GeoM.Translate(width, height)
	} else if h > 0 {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(width, 0)
	} else if v > 0 {
		options.GeoM.Scale(1, -1)
		options.GeoM.Translate(0, height)
	}
}

func (hexMap *HexMap) degreesToRadians(degrees int) float64 {
	return (math.Pi / 180) * float64(degrees)
}

func (hexMap *HexMap) rotateImageMatrix(options *ebiten.DrawImageOptions, img *ebiten.Image, rotateDegrees int) {
	imgWidth, imgHeight := img.Size()
	width := float64(imgWidth)
	height := float64(imgHeight)

	options.GeoM.Translate(-width/2, -height/2)
	options.GeoM.Rotate(hexMap.degreesToRadians(rotateDegrees))
	options.GeoM.Translate(width/2, height/2)
}

func (hexMap *HexMap) generateBackgroundImageOptions() {
	hexMap.BackgroundImageOptions = make(map[string]*LandImageOptions)
	for x := 0; x < 6; x++ {
		for y := 0; y < 4; y++ {
			bgDirtImgIndex := hexMap.rand(6)
			bgGrassImgIndex := hexMap.rand(6)
			flipH := hexMap.rand(2)
			flipV := hexMap.rand(2)
			rotateDegrees := hexMap.rand(4) * 90

			hexMap.BackgroundImageOptions[getFieldKey(x, y)] = &LandImageOptions{
				BgDirtImgIndex:  bgDirtImgIndex,
				BgGrassImgIndex: bgGrassImgIndex,
				FlipH:           flipH,
				FlipV:           flipV,
				RotateDegrees:   rotateDegrees,
			}
		}
	}
}

func (hexMap *HexMap) drawMapBackground() {
	dirtBgImg := make([]*ebiten.Image, 6)
	grassBgImg := make([]*ebiten.Image, 6)

	// Load images once
	for i := 0; i < 6; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/ld_%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		dirtBgImg[i] = bgImage
	}
	for i := 0; i < 6; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/l_%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		grassBgImg[i] = bgImage
	}

	for x := 0; x < 6; x++ {
		for y := 0; y < 4; y++ {
			destX := x*125 - 15
			destY := y*125 - 15
			options := &ebiten.DrawImageOptions{}

			backgroundImageOptions := hexMap.BackgroundImageOptions[getFieldKey(x, y)]
			dirtImg := dirtBgImg[backgroundImageOptions.BgDirtImgIndex]
			grassImg := grassBgImg[backgroundImageOptions.BgGrassImgIndex]
			flipH := backgroundImageOptions.FlipH
			flipV := backgroundImageOptions.FlipV
			rotateDegrees := backgroundImageOptions.RotateDegrees
			hexMap.flipImageMatrix(options, grassImg, flipH, flipV)
			hexMap.rotateImageMatrix(options, grassImg, rotateDegrees)

			options.GeoM.Translate(float64(destX), float64(destY))
			hexMap.BackgroundDirt.DrawImage(dirtImg, options)
			hexMap.BackgroundGrass.DrawImage(grassImg, options)
		}
	}
}

func getFieldKey(x int, y int) string {
	return "f" + fmt.Sprint(x) + "x" + fmt.Sprint(y)
}

func (hexMap *HexMap) getField(x int, y int, board *Board) *Field {
	key := getFieldKey(x, y)
	return board.Fields[key]
}

func (board *Board) getNeighborField(field *Field, neighborIndex int) *Field {
	neighborLocation := field.Neighbors[neighborIndex]
	if neighborLocation == nil {
		return nil
	}
	key := getFieldKey(neighborLocation.X, neighborLocation.Y)
	return board.Fields[key]
}

func (hexMap *HexMap) addField(x int, y int, board *Board) {
	key := getFieldKey(x, y)
	board.Fields[key] = &Field{}
	field := board.Fields[key]
	field.FX = x
	field.FY = y
	px := x*(board.HexWidth/4*3) + board.HexWidth/2
	var py int
	if x%2 == 0 {
		py = y*board.HexHeight + board.HexHeight/2
	} else {
		py = y*board.HexHeight + board.HexHeight
	}
	field.X = px
	field.Y = py
	field.LandId = -1
	if (x == 1 && y == 1) ||
		(x == board.XMax-2 && y == 1) ||
		(x == board.XMax-2 && y == board.YMax-2) ||
		(x == 1 && y == board.YMax-2) {
		field.Type = "land"
	} else {
		if hexMap.rand(10) <= 1 {
			field.Type = "land"
		} else {
			field.Type = "water"
		}
	}
	field.Capital = -1
	field.Estate = ""
	field.TownName = ""
}

func validateLocation(x int, y int, board *Board) *Point2D {
	if _, ok := board.Fields[getFieldKey(x, y)]; ok {
		return &Point2D{
			X: x,
			Y: y,
		}
	} else {
		return nil
	}
}

func (hexMap *HexMap) findNeighbors(field *Field, board *Board) {
	field.Neighbors = [6]*Point2D{}
	fx := field.FX
	fy := field.FY
	if fx%2 == 0 {
		field.Neighbors[0] = validateLocation(fx+1, fy, board)
		field.Neighbors[1] = validateLocation(fx, fy+1, board)
		field.Neighbors[2] = validateLocation(fx-1, fy, board)
		field.Neighbors[3] = validateLocation(fx-1, fy-1, board)
		field.Neighbors[4] = validateLocation(fx, fy-1, board)
		field.Neighbors[5] = validateLocation(fx+1, fy-1, board)
	} else {
		field.Neighbors[0] = validateLocation(fx+1, fy+1, board)
		field.Neighbors[1] = validateLocation(fx, fy+1, board)
		field.Neighbors[2] = validateLocation(fx-1, fy+1, board)
		field.Neighbors[3] = validateLocation(fx-1, fy, board)
		field.Neighbors[4] = validateLocation(fx, fy-1, board)
		field.Neighbors[5] = validateLocation(fx+1, fy, board)
	}
}

func (hexMap *HexMap) setLandFields(board *Board) {
	for x := 0; x < hexMap.Board.XMax; x++ {
		for y := 0; y < hexMap.Board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			if field.Type == "water" {
				landFields := 0
				for n := 0; n < 6; n++ {
					neighbor := board.getNeighborField(field, n)
					if neighbor == nil {
						continue
					}
					if neighbor.Type == "land" {
						landFields++
					}
				}
				if landFields >= 1 {
					hexMap.getField(x, y, board).IsLand = true
				}
			}
		}
	}

	for x := 0; x < hexMap.Board.XMax; x++ {
		for y := 0; y < hexMap.Board.YMax; y++ {
			if hexMap.getField(x, y, board).IsLand {
				hexMap.getField(x, y, board).Type = "land"
			}
		}
	}

	for x := 0; x < hexMap.Board.XMax; x++ {
		for y := 0; y < hexMap.Board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			if field.Type == "water" {
				waterFields := 0
				for n := 0; n < 6; n++ {
					neighbor := board.getNeighborField(field, n)
					if neighbor == nil {
						continue
					}
					if neighbor.Type == "water" {
						waterFields++
					}
				}

				if waterFields == 0 {
					hexMap.getField(x, y, board).Type = "land"
				}
			}
		}
	}
}

func (hexMap *HexMap) addNeighborsToLandGroup(field *Field, board *Board, landId int) int {
	newFields := 0
	for n := 0; n < 6; n++ {
		neighbor := board.getNeighborField(field, n)
		if neighbor != nil && neighbor.Type == "land" && neighbor.LandId < 0 {
			board.LandGroups[landId] = append(board.LandGroups[landId], neighbor)
			neighbor.LandId = landId
			newFields++
		}
	}
	return newFields
}

func (hexMap *HexMap) generateLandGroups(board *Board) {
	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			if hexMap.getField(x, y, board).Type == "land" {
				board.LandCount = board.LandCount + 1
			}
		}
	}

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			if hexMap.getField(x, y, board).Type == "land" && hexMap.getField(x, y, board).LandId < 0 {
				var countLandId = len(board.LandGroups)
				board.LandGroups = append(board.LandGroups, make([]*Field, 0))
				board.LandGroups[countLandId] = append(board.LandGroups[countLandId], hexMap.getField(x, y, board))
				hexMap.getField(x, y, board).LandId = countLandId
				groupSize := 0
				fieldCount := groupSize
				for groupSize >= fieldCount {
					groupSize = groupSize + hexMap.addNeighborsToLandGroup(board.LandGroups[countLandId][fieldCount], board, countLandId)
					fieldCount++
				}
			}
		}
	}
}

func (hexMap *HexMap) generatePartyCapitals(board *Board) {
	capital := 0
	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			if (x == 1 && y == 1) ||
				(x == board.XMax-2 && y == 1) ||
				(x == board.XMax-2 && y == board.YMax-2) ||
				(x == 1 && y == board.YMax-2) {
				hexMap.getField(x, y, board).Estate = "town"
				board.Towns = append(board.Towns, hexMap.getField(x, y, board))
				hexMap.getField(x, y, board).Capital = capital
				board.PartiesCapitals[capital] = hexMap.getField(x, y, board)
				capital++
			}
		}
	}
}

func (hexMap *HexMap) generateTowns(board *Board) {
	for landNum := 0; landNum < len(board.LandGroups); landNum++ {
		townCount := int(math.Floor((float64(len(board.LandGroups[landNum])) / 10) + 1))
		for townNum := 0; townNum < townCount; townNum++ {
			created := false
			attempts := 0
			for !created {
				attempts++
				if attempts > 10 {
					created = true
				}
				townIndex := hexMap.rand(len(board.LandGroups[landNum]))
				if board.LandGroups[landNum][townIndex].Estate == "" {
					ok := true
					for n := 0; n < 6; n++ {
						field := board.LandGroups[landNum][townIndex]
						neighbor := board.getNeighborField(field, n)
						if neighbor == nil {
							continue
						}
						if neighbor.Type == "water" || neighbor.Estate != "" {
							ok = false
						}
					}
					if ok {
						board.LandGroups[landNum][townIndex].Estate = "town"
						board.Towns = append(board.Towns, board.LandGroups[landNum][townIndex])
						created = true
					}
				}
			}
		}
	}
}

func (hexMap *HexMap) shuffle(arr []*Field) {
	for index := 0; index < len(arr); index++ {
		tmp := arr[index]
		randIndex := hexMap.rand(len(arr))
		// Swap with random index
		arr[index] = arr[randIndex]
		arr[randIndex] = tmp
	}
}

func (hexMap *HexMap) generatePorts(board *Board) {
	portNum := 0
	pathNum := 0
	for town := 0; town < len(board.Towns)-1; town++ {
		path := hexMap.Pathfinder.findPath(board, board.Towns[town], board.Towns[town+1], []string{"town"}, true)
		if path == nil || len(path) > portNum {
			path = hexMap.Pathfinder.findPath(board, board.Towns[town], board.Towns[town+1], []string{"town"}, false)
			pathNum++
		}
		for pathIndex := 1; pathIndex < len(path)-1; pathIndex++ {
			if path[pathIndex].Type == "land" && path[pathIndex+1].Type == "water" {
				path[pathIndex].Estate = "port"
				portNum++
			}
			if path[pathIndex].Type == "land" && path[pathIndex-1].Type == "water" {
				path[pathIndex].Estate = "port"
				portNum++
			}
		}
	}
}

func (hexMap *HexMap) assignWaterImageOptions(board *Board) {
	hexMap.WaterImageOptions = make(map[string]*WaterImageOptions)

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			if field.Type == "water" {
				waterBgImgIndex := hexMap.rand(6)
				flipH := hexMap.rand(2)
				flipV := hexMap.rand(2)
				rotateDegrees := hexMap.rand(2) * 180

				hexMap.WaterImageOptions[getFieldKey(x, y)] = &WaterImageOptions{
					WaterBgImgIndex: waterBgImgIndex,
					FlipH:           flipH,
					FlipV:           flipV,
					RotateDegrees:   rotateDegrees,
				}
			}
		}
	}
}

func (hexMap *HexMap) drawWaterAndPorts(board *Board) {
	waterBgImg := make([]*ebiten.Image, 6)
	portBgImg := make([]*ebiten.Image, 2)

	// Load images once
	for i := 0; i < 6; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/m_%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		waterBgImg[i] = bgImage
	}
	for i := 0; i < 2; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/m_p%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		portBgImg[i] = bgImage
	}

	portImageNum := [6]int{2, 1, 2, 2, 1, 2}
	portFlipH := [6]int{1, 0, 0, 0, 0, 1}
	portFlipV := [6]int{1, 1, 1, 0, 0, 0}
	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			if field.Type == "water" {
				options := &ebiten.DrawImageOptions{}

				waterImageOptions := hexMap.WaterImageOptions[getFieldKey(x, y)]
				waterImg := waterBgImg[waterImageOptions.WaterBgImgIndex]
				flipH := waterImageOptions.FlipH
				flipV := waterImageOptions.FlipV
				rotateDegrees := waterImageOptions.RotateDegrees
				hexMap.flipImageMatrix(options, waterImg, flipH, flipV)
				hexMap.rotateImageMatrix(options, waterImg, rotateDegrees)

				width, height := waterImg.Size()
				destX := field.X - width/2
				destY := field.Y - height/2
				options.GeoM.Translate(float64(destX), float64(destY))
				hexMap.BackgroundSea.DrawImage(waterImg, options)

				for n := 0; n < 6; n++ {
					field := hexMap.getField(x, y, board)
					neighbor := board.getNeighborField(field, n)
					if neighbor != nil && neighbor.Estate == "port" {
						portOptions := &ebiten.DrawImageOptions{}
						portImg := portBgImg[portImageNum[n]-1]
						hexMap.flipImageMatrix(portOptions, portImg, portFlipH[n], portFlipV[n])
						width, height := portImg.Size()
						portOptions.GeoM.Translate(float64(field.X)-float64(width/2), float64(field.Y)-float64(height/2))
						hexMap.BackgroundSea.DrawImage(portImg, portOptions)
					}
				}
			}
		}
	}
}

func (hexMap *HexMap) generateTownImageOptions() *LandImageOptions {
	bgDirtImgIndex := hexMap.rand(6)
	bgGrassImgIndex := hexMap.rand(6)
	flipH := hexMap.rand(2)
	flipV := hexMap.rand(2)
	rotateDegrees := hexMap.rand(360)

	return &LandImageOptions{
		BgDirtImgIndex:  bgDirtImgIndex,
		BgGrassImgIndex: bgGrassImgIndex,
		FlipH:           flipH,
		FlipV:           flipV,
		RotateDegrees:   rotateDegrees,
	}
}

func (hexMap *HexMap) addTown(
	townImageOptions *LandImageOptions,
	x int,
	y int,
	board *Board,
	townBgDirtImg []*ebiten.Image,
	townBgGrassImg []*ebiten.Image,
) {
	options := &ebiten.DrawImageOptions{}

	dirtImg := townBgDirtImg[townImageOptions.BgDirtImgIndex]
	grassImg := townBgGrassImg[townImageOptions.BgGrassImgIndex]
	flipH := townImageOptions.FlipH
	flipV := townImageOptions.FlipV
	rotateDegrees := townImageOptions.RotateDegrees
	hexMap.flipImageMatrix(options, grassImg, flipH, flipV)
	hexMap.rotateImageMatrix(options, grassImg, rotateDegrees)

	width, height := grassImg.Size()
	destX := hexMap.getField(x, y, board).X - (width / 2)
	destY := hexMap.getField(x, y, board).Y - (height / 2)
	options.GeoM.Translate(float64(destX), float64(destY))
	hexMap.BackgroundDirt.DrawImage(dirtImg, options)
	hexMap.BackgroundGrass.DrawImage(grassImg, options)
}

func (hexMap *HexMap) randTown() string {
	townNames := hexMap.Board.TownNames
	randIndex := hexMap.rand(len(townNames))

	// Swap values between index 0 and randIndex
	townName := townNames[randIndex]
	townNames[randIndex] = townNames[0]
	townNames[0] = townName

	// Return the town at index 0 and remove from the list
	hexMap.Board.TownNames = townNames[1:]
	return townName
}

func (hexMap *HexMap) assignTownNames(board *Board) {
	hexMap.TownImageOptions = make(map[string]*LandImageOptions)

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			switch hexMap.getField(x, y, board).Estate {
			case "town":
				hexMap.TownImageOptions[getFieldKey(x, y)] = hexMap.generateTownImageOptions()
				hexMap.getField(x, y, board).TownName = hexMap.randTown()
				break
			case "port":
				hexMap.TownImageOptions[getFieldKey(x, y)] = hexMap.generateTownImageOptions()
				hexMap.getField(x, y, board).TownName = hexMap.randTown()
				break
			}
		}
	}
}

func (hexMap *HexMap) drawTownBackground(board *Board) {
	townBgDirtImg := make([]*ebiten.Image, 6)
	townBgGrassImg := make([]*ebiten.Image, 6)

	// Load images once
	for i := 0; i < 6; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/cd_%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		townBgDirtImg[i] = bgImage
	}
	for i := 0; i < 6; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf(fileDir+"/c_%d.png", i+1))
		if err != nil {
			log.Fatal(err)
		}
		townBgGrassImg[i] = bgImage
	}

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			switch hexMap.getField(x, y, board).Estate {
			case "town":
				townImageOptions := hexMap.TownImageOptions[getFieldKey(x, y)]
				hexMap.addTown(townImageOptions, x, y, board, townBgDirtImg, townBgGrassImg)
				break
			case "port":
				townImageOptions := hexMap.TownImageOptions[getFieldKey(x, y)]
				hexMap.addTown(townImageOptions, x, y, board, townBgDirtImg, townBgGrassImg)
				break
			}
		}
	}
}

func (hexMap *HexMap) drawHex(background *ebiten.Image, xCenter float64, yCenter float64) {
	lineColor := color.RGBA{uint8(255), uint8(255), uint8(102), 50}
	hexMap.drawLine(background, xCenter-12.5, yCenter-20, xCenter-23, yCenter-0, lineColor)
	hexMap.drawLine(background, xCenter-23, yCenter-0, xCenter-12.5, yCenter+20, lineColor)
	hexMap.drawLine(background, xCenter-12.5, yCenter+20, xCenter+12.5, yCenter+20, lineColor)
	hexMap.drawLine(background, xCenter+12.5, yCenter+20, xCenter+23, yCenter+0, lineColor)
	hexMap.drawLine(background, xCenter+23, yCenter+0, xCenter+12.5, yCenter-20, lineColor)
	hexMap.drawLine(background, xCenter+12.5, yCenter-20, xCenter-12.5, yCenter-20, lineColor)
}

func (hexMap *HexMap) drawLine(background *ebiten.Image, x1 float64, y1 float64, x2 float64, y2 float64, lineColor color.RGBA) {
	vector.StrokeLine(background, float32(x1), float32(y1), float32(x2), float32(y2), 0.5, lineColor, true)
}

func (hexMap *HexMap) drawTowns(board *Board) {
	capitalsImg := make([]*ebiten.Image, 4)
	partyImgNames := [4]string{"capital_red.png", "capital_violet.png", "capital_blue.png", "capital_green.png"}

	// Load images once
	for i := 0; i < 4; i++ {
		bgImage, _, err := ebitenutil.NewImageFromFile(fileDir + "/" + partyImgNames[i])
		if err != nil {
			log.Fatal(err)
		}
		capitalsImg[i] = bgImage
	}

	cityImage, _, err := ebitenutil.NewImageFromFile(fileDir + "/city.png")
	if err != nil {
		log.Fatal(err)
	}
	portImage, _, err := ebitenutil.NewImageFromFile(fileDir + "/port.png")
	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < hexMap.Board.XMax; x++ {
		for y := 0; y < hexMap.Board.YMax; y++ {
			options := &ebiten.DrawImageOptions{}
			field := hexMap.getField(x, y, board)

			// Draw outline
			hexMap.drawHex(hexMap.BackgroundTiles, float64(field.X), float64(field.Y))

			if field.Estate == "town" {
				width, height := cityImage.Size()
				options.GeoM.Translate(float64(field.X)-float64(width/2), float64(field.Y)-float64(height/2))

				if field.Capital >= 0 {
					// Draw capital icon to mark each player's capital
					hexMap.BackgroundTiles.DrawImage(capitalsImg[field.Capital], options)
				} else {
					hexMap.BackgroundTiles.DrawImage(cityImage, options)
				}
			} else if field.Estate == "port" {
				width, height := portImage.Size()
				options.GeoM.Translate(float64(field.X)-float64(width/2), float64(field.Y)-float64(height/2))
				hexMap.BackgroundTiles.DrawImage(portImage, options)
			}
		}
	}
}

func (hexMap *HexMap) drawTownNames(board *Board) {
	for x := 0; x < hexMap.Board.XMax; x++ {
		for y := 0; y < hexMap.Board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			if field.Estate == "town" || field.Estate == "port" {
				// Draw text twice to make the font bold and easier to read
				textX := field.X - (len(field.TownName) * 3)
				textY := field.Y - int(float64(hexMap.Board.HexHeight)/3)
				text.Draw(hexMap.BackgroundTiles, field.TownName, hexMap.TextFont, textX, textY, color.White)
				text.Draw(hexMap.BackgroundTiles, field.TownName, hexMap.TextFont, textX, textY, color.White)
			}
		}
	}
}

func (hexMap *HexMap) generateBoard(board *Board) {
	hexMap.generateBackgroundImageOptions()

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			hexMap.addField(x, y, board)
		}
	}

	for x := 0; x < board.XMax; x++ {
		for y := 0; y < board.YMax; y++ {
			field := hexMap.getField(x, y, board)
			hexMap.findNeighbors(field, board)
		}
	}

	hexMap.setLandFields(board)
	hexMap.generateLandGroups(board)
	hexMap.generatePartyCapitals(board)
	hexMap.generateTowns(board)
	hexMap.shuffle(board.Towns)
	hexMap.generatePorts(board)
	hexMap.assignWaterImageOptions(board)
	hexMap.assignTownNames(board)
}

func (hexMap *HexMap) generateMap() {
	hexMap.generateBoard(hexMap.Board)

	hexMap.drawMapBackground()
	hexMap.drawWaterAndPorts(hexMap.Board)
	hexMap.drawTownBackground(hexMap.Board)
	hexMap.drawTowns(hexMap.Board)
	hexMap.drawTownNames(hexMap.Board)

	// Crop image to only include hexes
	mapView := image.Rect(0, 0,
		int(float64(hexMap.Board.HexWidth)*(float64(hexMap.Board.XMax)/1.35)),
		int(float64(hexMap.Board.HexHeight)*(float64(hexMap.Board.YMax)+0.5)))

	hexMap.Background.DrawImage(hexMap.BackgroundGrass.SubImage(mapView).(*ebiten.Image), nil)
	hexMap.Background.DrawImage(hexMap.BackgroundSea.SubImage(mapView).(*ebiten.Image), nil)
	hexMap.Background.DrawImage(hexMap.BackgroundTiles, nil)

	// Draw UI
	redButtonImage, _, err := ebitenutil.NewImageFromFile(fileDir + "/red_button.png")
	if err != nil {
		log.Fatal(err)
	}
	randomMapOptions := &ebiten.DrawImageOptions{}
	randomMapOptions.GeoM.Translate(250, 500)
	hexMap.UI.DrawImage(redButtonImage, randomMapOptions)
	text.Draw(hexMap.UI, "Random Map", hexMap.TextFont, 260, 517, color.White)

	grayButtonImage, _, err := ebitenutil.NewImageFromFile(fileDir + "/gray_button.png")
	if err != nil {
		log.Fatal(err)
	}
	mapNumberOptions := &ebiten.DrawImageOptions{}
	mapNumberOptions.GeoM.Translate(150, 500)
	hexMap.UI.DrawImage(grayButtonImage, mapNumberOptions)
	text.Draw(hexMap.UI, fmt.Sprintf("%v", hexMap.Board.MapNumber), hexMap.TextFont, 175, 517, color.Black)

	text.Draw(hexMap.UI, "Map Number", hexMap.TextFont, 75, 517, color.White)
}

func (hexMap *HexMap) isMouseCursorOnRandomMapButton(x int, y int) bool {
	return x >= 255 && x <= 345 && y >= 500 && y <= 530
}

func (hexMap *HexMap) drawBackground(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(25, 25)
	screen.DrawImage(hexMap.Background, options)
	screen.DrawImage(hexMap.UI, nil)
}
