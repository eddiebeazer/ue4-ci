package jsonToXml

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSuccessJUnitTest(t *testing.T) {
	err := ParseTestOutput("successTest.json", "successTest.xml", "")
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)

	// Open our xmlFile
	xmlFile, err := os.Open("successTest.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer func(xmlFile *os.File) {
		err := xmlFile.Close()
		if err != nil {

		}
	}(xmlFile)

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var report testsuite
	err = xml.Unmarshal(byteValue, &report)
	if err != nil {
		print(err)
	}
	assert.Nil(t, err)

	assert.Equal(t, 309, report.Tests)
	assert.Equal(t, 0, report.Failures)
}

func TestFailedJUnitTest(t *testing.T) {
	err := ParseTestOutput("failedTest.json", "failedXml.xml", "")
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)

	// Open our xmlFile
	xmlFile, err := os.Open("failedXml.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer func(xmlFile *os.File) {
		err := xmlFile.Close()
		if err != nil {

		}
	}(xmlFile)

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var report testsuite
	err = xml.Unmarshal(byteValue, &report)
	if err != nil {
		print(err)
	}
	assert.Nil(t, err)

	assert.Equal(t, 309, report.Tests)
	assert.Equal(t, 1, report.Failures)
}
