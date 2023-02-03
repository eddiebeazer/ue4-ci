package clean

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCleanDirectory(t *testing.T) {
	dirToTest := "Test"
	err := os.Mkdir(dirToTest, 777)
	assert.Nil(t, err)

	err = Dist(dirToTest)
	assert.Nil(t, err)

	err = Dist("")
	assert.NotNil(t, err)
}
