package hexempire

import "math"

// MapGenerator deterministically builds a Board (and the visual-variation
// choices needed to render it) from a seed. It owns all randomness used
// during generation; MapRenderer consumes its output but never calls
// random itself, so a given seed always reproduces the same map.
type MapGenerator struct {
	mapNumber  int
	random     *SeededRandom
	pathfinder *Pathfinder
}

func NewMapGenerator(mapNumber int) *MapGenerator {
	return &MapGenerator{
		mapNumber:  mapNumber,
		random:     NewSeededRandom(mapNumber),
		pathfinder: NewPathfinder(),
	}
}

// GeneratedMap is the deterministic output of a MapGenerator run: the board
// itself plus the sprite-variation choices to use when rendering it.
type GeneratedMap struct {
	Board                  *Board
	BackgroundImageOptions map[string]*LandImageOptions
	WaterImageOptions      map[string]*WaterImageOptions
	TownImageOptions       map[string]*LandImageOptions
}

func (gen *MapGenerator) Generate() *GeneratedMap {
	board := NewBoard()
	board.MapNumber = gen.mapNumber
	board.TownNames = generateAllTowns()

	backgroundImageOptions := gen.generateBackgroundImageOptions()

	board.forEachField(func(x, y int) { gen.addField(x, y, board) })
	board.forEachField(func(x, y int) { board.findNeighbors(board.fieldAt(x, y)) })

	gen.setLandFields(board)
	gen.generateLandGroups(board)
	gen.generatePartyCapitals(board)
	gen.generateTowns(board)
	gen.shuffle(board.Towns)
	gen.generatePorts(board)
	waterImageOptions := gen.assignWaterImageOptions(board)
	townImageOptions := gen.assignTownNames(board)

	return &GeneratedMap{
		Board:                  board,
		BackgroundImageOptions: backgroundImageOptions,
		WaterImageOptions:      waterImageOptions,
		TownImageOptions:       townImageOptions,
	}
}

func (gen *MapGenerator) addField(x int, y int, board *Board) {
	field := &Field{}
	board.Fields[fieldKey(x, y)] = field
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
	if board.isCapitalCorner(x, y) {
		field.Type = FieldTypeLand
	} else if gen.random.Intn(10) <= 1 {
		field.Type = FieldTypeLand
	} else {
		field.Type = FieldTypeWater
	}
	field.Capital = -1
	field.Estate = EstateNone
	field.TownName = ""
}

// setLandFields promotes water fields to land in two passes: first any
// water field touching land becomes land (coastline fill-in), then any
// remaining water field with no water neighbor (a landlocked puddle) is
// also converted to land, since isolated single-tile lakes aren't playable.
func (gen *MapGenerator) setLandFields(board *Board) {
	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Type != FieldTypeWater {
			return
		}
		if gen.countNeighborsOfType(board, field, FieldTypeLand) >= 1 {
			field.IsLand = true
		}
	})

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.IsLand {
			field.Type = FieldTypeLand
		}
	})

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Type != FieldTypeWater {
			return
		}
		if gen.countNeighborsOfType(board, field, FieldTypeWater) == 0 {
			field.Type = FieldTypeLand
		}
	})
}

func (gen *MapGenerator) countNeighborsOfType(board *Board, field *Field, fieldType FieldType) int {
	count := 0
	for n := 0; n < 6; n++ {
		neighbor := board.neighborOf(field, n)
		if neighbor != nil && neighbor.Type == fieldType {
			count++
		}
	}
	return count
}

func (gen *MapGenerator) generateLandGroups(board *Board) {
	board.forEachField(func(x, y int) {
		if board.fieldAt(x, y).Type == FieldTypeLand {
			board.LandCount++
		}
	})

	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Type != FieldTypeLand || field.LandId >= 0 {
			return
		}

		landId := len(board.LandGroups)
		board.LandGroups = append(board.LandGroups, []*Field{field})
		field.LandId = landId

		fieldCount := 0
		groupSize := 0
		for groupSize >= fieldCount {
			groupSize += gen.addNeighborsToLandGroup(board.LandGroups[landId][fieldCount], board, landId)
			fieldCount++
		}
	})
}

func (gen *MapGenerator) addNeighborsToLandGroup(field *Field, board *Board, landId int) int {
	newFields := 0
	for n := 0; n < 6; n++ {
		neighbor := board.neighborOf(field, n)
		if neighbor != nil && neighbor.Type == FieldTypeLand && neighbor.LandId < 0 {
			board.LandGroups[landId] = append(board.LandGroups[landId], neighbor)
			neighbor.LandId = landId
			newFields++
		}
	}
	return newFields
}

func (gen *MapGenerator) generatePartyCapitals(board *Board) {
	capital := 0
	board.forEachField(func(x, y int) {
		if !board.isCapitalCorner(x, y) {
			return
		}
		field := board.fieldAt(x, y)
		field.Estate = EstateTown
		board.Towns = append(board.Towns, field)
		field.Capital = capital
		board.PartiesCapitals[capital] = field
		capital++
	})
}

func (gen *MapGenerator) generateTowns(board *Board) {
	for landNum := 0; landNum < len(board.LandGroups); landNum++ {
		group := board.LandGroups[landNum]
		townCount := int(math.Floor((float64(len(group)) / 10) + 1))
		for townNum := 0; townNum < townCount; townNum++ {
			gen.placeTownInGroup(board, group)
		}
	}
}

// placeTownInGroup repeatedly picks a random field in the group until it
// finds one with no estate and no water or occupied neighbor, or gives up
// after 10 failed attempts.
func (gen *MapGenerator) placeTownInGroup(board *Board, group []*Field) {
	attempts := 0
	for {
		attempts++
		candidate := group[gen.random.Intn(len(group))]
		if candidate.Estate == EstateNone && gen.hasNoWaterOrOccupiedNeighbor(board, candidate) {
			candidate.Estate = EstateTown
			board.Towns = append(board.Towns, candidate)
			return
		}
		if attempts > 10 {
			return
		}
	}
}

func (gen *MapGenerator) hasNoWaterOrOccupiedNeighbor(board *Board, field *Field) bool {
	for n := 0; n < 6; n++ {
		neighbor := board.neighborOf(field, n)
		if neighbor == nil {
			continue
		}
		if neighbor.Type == FieldTypeWater || neighbor.Estate != EstateNone {
			return false
		}
	}
	return true
}

func (gen *MapGenerator) shuffle(arr []*Field) {
	for index := 0; index < len(arr); index++ {
		randIndex := gen.random.Intn(len(arr))
		arr[index], arr[randIndex] = arr[randIndex], arr[index]
	}
}

// generatePorts connects each consecutive pair of towns with a path,
// preferring one that avoids water; any land field on that path adjacent
// to water becomes a port.
func (gen *MapGenerator) generatePorts(board *Board) {
	portNum := 0
	for town := 0; town < len(board.Towns)-1; town++ {
		path := gen.pathfinder.findPath(board, board.Towns[town], board.Towns[town+1], []Estate{EstateTown}, true)
		if path == nil || len(path) > portNum {
			path = gen.pathfinder.findPath(board, board.Towns[town], board.Towns[town+1], []Estate{EstateTown}, false)
		}
		portNum += markPortsAlongPath(path)
	}
}

func markPortsAlongPath(path []*Field) int {
	portsAdded := 0
	for i := 1; i < len(path)-1; i++ {
		if path[i].Type == FieldTypeLand && path[i+1].Type == FieldTypeWater {
			path[i].Estate = EstatePort
			portsAdded++
		}
		if path[i].Type == FieldTypeLand && path[i-1].Type == FieldTypeWater {
			path[i].Estate = EstatePort
			portsAdded++
		}
	}
	return portsAdded
}

func (gen *MapGenerator) assignWaterImageOptions(board *Board) map[string]*WaterImageOptions {
	options := make(map[string]*WaterImageOptions)
	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Type != FieldTypeWater {
			return
		}

		waterBgImgIndex := gen.random.Intn(6)
		flipH := gen.random.Intn(2)
		flipV := gen.random.Intn(2)
		rotateDegrees := gen.random.Intn(2) * 180

		options[fieldKey(x, y)] = &WaterImageOptions{
			WaterBgImgIndex: waterBgImgIndex,
			FlipH:           flipH,
			FlipV:           flipV,
			RotateDegrees:   rotateDegrees,
		}
	})
	return options
}

// generateBackgroundImageOptions picks a sprite variant for each tile of
// the repeating 6x4 dirt/grass texture used to paint the whole background.
func (gen *MapGenerator) generateBackgroundImageOptions() map[string]*LandImageOptions {
	options := make(map[string]*LandImageOptions)
	for x := 0; x < 6; x++ {
		for y := 0; y < 4; y++ {
			bgDirtImgIndex := gen.random.Intn(6)
			bgGrassImgIndex := gen.random.Intn(6)
			flipH := gen.random.Intn(2)
			flipV := gen.random.Intn(2)
			rotateDegrees := gen.random.Intn(4) * 90

			options[fieldKey(x, y)] = &LandImageOptions{
				BgDirtImgIndex:  bgDirtImgIndex,
				BgGrassImgIndex: bgGrassImgIndex,
				FlipH:           flipH,
				FlipV:           flipV,
				RotateDegrees:   rotateDegrees,
			}
		}
	}
	return options
}

func (gen *MapGenerator) generateTownImageOptions() *LandImageOptions {
	bgDirtImgIndex := gen.random.Intn(6)
	bgGrassImgIndex := gen.random.Intn(6)
	flipH := gen.random.Intn(2)
	flipV := gen.random.Intn(2)
	rotateDegrees := gen.random.Intn(360)

	return &LandImageOptions{
		BgDirtImgIndex:  bgDirtImgIndex,
		BgGrassImgIndex: bgGrassImgIndex,
		FlipH:           flipH,
		FlipV:           flipV,
		RotateDegrees:   rotateDegrees,
	}
}

// randTown pops a random, not-yet-used name off board.TownNames.
func (gen *MapGenerator) randTown(board *Board) string {
	townNames := board.TownNames
	randIndex := gen.random.Intn(len(townNames))

	townName := townNames[randIndex]
	townNames[randIndex] = townNames[0]
	townNames[0] = townName

	board.TownNames = townNames[1:]
	return townName
}

func (gen *MapGenerator) assignTownNames(board *Board) map[string]*LandImageOptions {
	options := make(map[string]*LandImageOptions)
	board.forEachField(func(x, y int) {
		field := board.fieldAt(x, y)
		if field.Estate != EstateTown && field.Estate != EstatePort {
			return
		}
		options[fieldKey(x, y)] = gen.generateTownImageOptions()
		field.TownName = gen.randTown(board)
	})
	return options
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
