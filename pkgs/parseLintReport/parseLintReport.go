package parseLintReport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dimchansky/utfbom"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strings"
)

type LinterJson struct {
	Violators []Violator `json:"Violators"`
}

type Violator struct {
	ViolatorAssetName string      `json:"ViolatorAssetName"`
	ViolatorAssetPath string      `json:"ViolatorAssetPath"`
	ViolatorFullName  string      `json:"ViolatorFullName"`
	Violations        []Violation `json:"Violations"`
}

type Violation struct {
	RuleGroup             string `json:"RuleGroup"`
	RuleTitle             string `json:"RuleTitle"`
	RuleDesc              string `json:"RuleDesc"`
	RuleSeverity          int    `json:"RuleSeverity"`
	RuleRecommendedAction string `json:"RuleRecommendedAction"`
}

func ParseReport(jsonFilePath string) error {
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return cli.Exit(fmt.Errorf("error reading json file: %s", err), 1)
	}

	jsonDataWithoutBom, err := ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(jsonData)))
	if err != nil {
		return cli.Exit(fmt.Errorf("failed to remove bom from json file: %s", err), 1)
	}

	testResults := LinterJson{}

	err = json.Unmarshal(jsonDataWithoutBom, &testResults)

	errorCount := 0
	warningCount := 0

	w := os.Stdout
	_, err = fmt.Fprintf(w, "##teamcity[testSuiteStarted name='%s']\n", "Linter")
	if err != nil {
		return cli.Exit(fmt.Errorf("error writing message: %s", err), 1)
	}

	for _, violator := range testResults.Violators {
		_, err = fmt.Fprintf(w, "##teamcity[testStarted name='%s: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath)
		if err != nil {
			return cli.Exit(fmt.Errorf("error writing message: %s", err), 1)
		}

		var warnings []string
		var errors []string

		for _, violations := range violator.Violations {
			formattedMessage := ""
			if violations.RuleRecommendedAction == "" {
				formattedMessage = fmt.Sprintf("%s - %s", violations.RuleGroup, violations.RuleDesc)
			} else {
				formattedMessage = fmt.Sprintf("%s - %s. Fix: %s", violations.RuleGroup, violations.RuleDesc, violations.RuleRecommendedAction)
			}
			// 0 = error, 1 = warn
			if violations.RuleSeverity == 0 {
				errors = append(errors, formattedMessage)
				errorCount += 1
			} else {
				warnings = append(warnings, formattedMessage)
				warningCount += 1
			}
		}
		if len(errors) > 0 {
			_, err = fmt.Fprintf(w, "##teamcity[testFailed name='%s: %s' message='%s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath, strings.Join(errors, "\n"))
		}
		if len(warnings) > 0 {
			_, err = fmt.Fprintf(w, "##teamcity[testFailed name='%s: %s' out='warning: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath, strings.Join(warnings, "\n"))
		}

		_, err = fmt.Fprintf(w, "##teamcity[testFinished name='%s: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath)
	}

	_, err = fmt.Fprintf(w, "##teamcity[testSuiteFinished name='%s']\n", "Linter")
	_, err = fmt.Fprintf(w, "##teamcity[buildStatisticValue key='%s' value='%d']\n", "Lint Errors", errorCount)
	_, err = fmt.Fprintf(w, "##teamcity[buildStatisticValue key='%s' value='%d']\n", "Lint Warnings", warningCount)

	return nil
}
