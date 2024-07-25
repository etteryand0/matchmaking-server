package match

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func LogMatch(c *gin.Context) {
	testName, foundTestNameParam := c.GetQuery("test_name")
	epoch, foundEpochParam := c.GetQuery("epoch")

	if !foundTestNameParam || !foundEpochParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	testFile, err := os.Open(path.Join("tests", testName, "test.json"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	defer testFile.Close()

	byteValue, err := io.ReadAll(testFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while reading epoch file"})
		return
	}

	var testCollection map[string]*json.RawMessage
	if err := json.Unmarshal(byteValue, &testCollection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while unmarshalling test json file"})
		return
	}

	nextEpoch, ok := testCollection[epoch]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such epoch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"epoch":     nextEpoch,
		"test_name": testName,
	})
}
