package jsonToXml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dimchansky/utfbom"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

// XML for JUnit
type testsuite struct {
	Errors    int           `xml:"errors,attr"`
	Failures  int           `xml:"failures,attr"`
	Tests     int           `xml:"tests,attr"`
	Skipped   int           `xml:"skipped,attr"`
	Time      float64       `xml:"time,attr"`
	Name      string        `xml:"name,attr"`
	TimeStamp string        `xml:"timestamp,attr"`
	TestCase  []interface{} `xml:"testcase"`
}

type TestSuccess struct {
	ClassName string  `xml:"classname,attr"`
	Name      string  `xml:"name,attr"`
	Time      float64 `xml:"time,attr"`
}

type TestFailed struct {
	ClassName string     `xml:"classname,attr"`
	Name      string     `xml:"name,attr"`
	Time      float64    `xml:"time,attr"`
	Failure   FailureTag `xml:"failure,omitempty"`
}

type FailureTag struct {
	Message       string `xml:"message,attr"`
	Type          string `xml:"type,attr"`
	OutputMessage string `xml:",innerxml"`
}

// UE4 JSON Test Unit Structure
type AutomationTestJson struct {
	Succeeded     int     `json:"succeeded"`
	Failed        int     `json:"failed"`
	TotalDuration float64 `json:"totalDuration"`
	TimeStamp     string  `json:"reportCreatedOn"`
	Tests         []Test  `json:"tests"`
}

type Test struct {
	TestDisplayName string      `json:"testDisplayName"`
	FullTestPath    string      `json:"fullTestPath"`
	Duration        float64     `json:"Duration"`
	Errors          int         `json:"errors"`
	Entries         []TestEntry `json:"entries,omitempty"`
}

type TestEntry struct {
	Event TestEvent `json:"event"`
}

type TestEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ParseTestOutput(jsonFilePath string, outPath string, testSuiteName string) error {
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return cli.Exit(fmt.Errorf("error reading json file: %s", err), 1)
	}

	jsonDataWithoutBom, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(jsonData)))
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to remove bom from json file: %s", err), 1)
	}

	testResults := AutomationTestJson{}

	err = json.Unmarshal(jsonDataWithoutBom, &testResults)
	if err != nil {
		return cli.Exit(fmt.Errorf("error unmarshalling json file: %s", err), 1)
	}

	var testCases []interface{}

	for _, test := range testResults.Tests {
		className := test.FullTestPath
		name := test.TestDisplayName
		time := test.Duration

		// If the test fails, add in new failed test xml else add in standard success
		if test.Errors > 0 {
			testCases = append(testCases, &TestFailed{
				ClassName: className,
				Name:      name,
				Time:      time,
				Failure: FailureTag{
					Message:       test.Entries[0].Event.Message,
					Type:          test.Entries[0].Event.Type,
					OutputMessage: test.Entries[0].Event.Message,
				}})
		} else {
			testCases = append(testCases, &TestSuccess{
				ClassName: className,
				Name:      name,
				Time:      time,
			})
		}
	}

	uatXml := testsuite{
		Errors:    0,
		Failures:  testResults.Failed,
		Name:      testSuiteName,
		Time:      testResults.TotalDuration,
		Skipped:   0,
		Tests:     len(testResults.Tests),
		TestCase:  testCases,
		TimeStamp: testResults.TimeStamp,
	}

	file, _ := xml.MarshalIndent(uatXml, "", " ")
	jUnitXML := []byte(xml.Header + string(file))

	_ = ioutil.WriteFile(outPath, jUnitXML, 0644)

	return nil
}
