package match

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/perimeterx/marshmallow"
)

func LogMatch(c *gin.Context) {
	testName, foundTestNameParam := c.GetQuery("test_name")
	epoch, foundEpochParam := c.GetQuery("epoch")

	if !foundTestNameParam || !foundEpochParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}
	if epoch == "last" {
		c.JSON(http.StatusBadRequest, gin.H{"Nostradamus": "No... no... no..."})
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

	testCollection := TestCollection{}
	result, err := marshmallow.Unmarshal(byteValue, &testCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while unmarshalling test json file"})
		return
	}

	nextEpoch, ok := result[epoch]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such epoch"})
		return
	}

	if testCollection.Last == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No last epoch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"epoch":         nextEpoch,
		"is_last_epoch": testCollection.Last == nextEpoch,
	})
}
