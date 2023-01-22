package hexempire

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratedBoard(t *testing.T) {
	hexMap := NewHexMap(148660, nil)
	hexMap.generateBoard(hexMap.Board)

	generatedBoardJson, err := json.MarshalIndent(hexMap.Board, "", " ")
	if err != nil {
		log.Fatal("Failed to marshal map data: ", err)
	}

	referenceFilename := "../generated/map148660.json"
	jsonFile, err := os.Open(referenceFilename)
	if err != nil {
		log.Fatal("Failed to open json file", err)
	}
	defer jsonFile.Close()

	referenceBoardJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	assert.JSONEq(t, string(generatedBoardJson[:]), string(referenceBoardJson[:]))
}
