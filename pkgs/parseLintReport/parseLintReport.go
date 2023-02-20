package parseLintReport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dimchansky/utfbom"
	"github.com/urfave/cli/v2"
	"io"
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

func EscapeTeamCityString(str string) string {
	newString := strings.Replace(str, "|", "||", -1)
	newString = strings.Replace(newString, "'", "|'", -1)
	newString = strings.Replace(newString, "\r", "|r", -1)
	newString = strings.Replace(newString, "\n", "|n", -1)
	newString = strings.Replace(newString, "]", "|]", -1)
	newString = strings.Replace(newString, "[", "|[", -1)
	return newString
}

func WriteTeamCityMsg(w io.Writer, str string) error {
	_, err := fmt.Fprintf(w, str)
	if err != nil {
		return cli.Exit(fmt.Errorf("error writing message: %s", err), 1)
	}
	return nil
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
		err = WriteTeamCityMsg(w, fmt.Sprintf("##teamcity[testStarted name='%s: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath))
		if err != nil {
			return err
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
			err = WriteTeamCityMsg(w, EscapeTeamCityString(fmt.Sprintf("##teamcity[testFailed name='%s: %s' message='%s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath, strings.Join(errors, "\n"))))
			if err != nil {
				return err
			}
		}
		if len(warnings) > 0 {
			err = WriteTeamCityMsg(w, EscapeTeamCityString(fmt.Sprintf("##teamcity[testStdOut name='%s: %s' out='warning: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath, strings.Join(warnings, "\n"))))
			if err != nil {
				return err
			}
		}

		err = WriteTeamCityMsg(w, fmt.Sprintf("##teamcity[testFinished name='%s: %s']\n", violator.ViolatorAssetName, violator.ViolatorAssetPath))
		if err != nil {
			return err
		}
	}

	err = WriteTeamCityMsg(w, fmt.Sprintf("##teamcity[testSuiteFinished name='%s']\n", "Linter"))
	if err != nil {
		return err
	}
	err = WriteTeamCityMsg(w, fmt.Sprintf("##teamcity[buildStatisticValue key='%s' value='%d']\n", "Lint Errors", errorCount))
	if err != nil {
		return err
	}
	err = WriteTeamCityMsg(w, fmt.Sprintf("##teamcity[buildStatisticValue key='%s' value='%d']\n", "Lint Warnings", warningCount))
	if err != nil {
		return err
	}

	return nil
}
