package parseLintReport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLintReport(t *testing.T) {
	err := ParseReport("./Lintreport.json")
	assert.Nil(t, err)
}
