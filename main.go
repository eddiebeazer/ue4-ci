package main

import (
	"github.com/eddiebeazer/ue4-ci/pkgs/clean"
	"github.com/eddiebeazer/ue4-ci/pkgs/jsonToXml"
	"github.com/eddiebeazer/ue4-ci/pkgs/projectVersion"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "ue4-ci helper",
		Usage:     "Tools for CI pipelines",
		Version:   "0.1.0",
		UsageText: "Tools for CI pipelines.  This CLI will assume the command is being ran in a directory with a .uproject file in it",
		Suggest:   true,
		Authors: []*cli.Author{
			{
				Name:  "Edward Beazer",
				Email: "ebeazer@thedigitalsages.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "clean",
				Usage: "Deletes build files from the project",
				Subcommands: []*cli.Command{
					{
						Name:  "dist",
						Usage: "Deletes the given archive/dist folder from the project",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "directory",
								Aliases:  []string{"d"},
								Value:    "",
								Usage:    "Location of the archive/dist folder to delete. Relative or absolute",
								Required: true,
							},
						},
						Action: func(cCtx *cli.Context) error {
							directory := cCtx.String("directory")
							return clean.Dist(directory)
						},
					},
				},
			},
			{
				Name:  "jsonToXml",
				Usage: "Parses the UAT tools output json into JUnit XML for ci's",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "jsonFile",
						Aliases:  []string{"j"},
						Value:    "",
						Usage:    "UAT test output json file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "outPath",
						Aliases:  []string{"o"},
						Value:    "",
						Usage:    "Output path of the new JUnit Xml file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "testSuiteName",
						Aliases: []string{"t"},
						Value:   "Unreal Automation Testing JUnit Test Report",
						Usage:   "Test suite name.  Optional",
					},
				},
				Action: func(cCtx *cli.Context) error {
					jsonFilePath := cCtx.String("jsonFile")
					outPath := cCtx.String("outPath")
					testSuiteName := cCtx.String("testSuiteName")
					return jsonToXml.ParseTestOutput(jsonFilePath, outPath, testSuiteName)
				},
			},
			{
				Name:  "projectVersion",
				Usage: "Sets and gets the current project version from DefaultGame.ini (must be writable)",
				Subcommands: []*cli.Command{
					{
						Name:  "set",
						Usage: "Sets the project version to the inputted string",
						Subcommands: []*cli.Command{
							{
								Name: "perforce",
								Usage: "Sets the version using perforce style branching. ex - rel0.2, dev0.1, patch4.3 and task9.4 \n" +
									"The ideal usage of this command is to save the current dev and rel versions as separate vars \n" +
									"and then provide them as parameters.  The new sem ver version will be spit out at the end.",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "relVersion",
										Aliases:  []string{"s"},
										Value:    "",
										Usage:    "Current release version of the project",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "devVersion",
										Aliases:  []string{"s"},
										Value:    "",
										Usage:    "Current dev version of the project",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "branch",
										Aliases:  []string{"b"},
										Value:    "",
										Usage:    "Branch of the build",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "iniPath",
										Aliases:  []string{"i"},
										Value:    "",
										Usage:    "file path to the DefaultGame ini file",
										Required: true,
									},
								},
								Action: func(cCtx *cli.Context) error {
									devVersion := cCtx.String("devVersion")
									relVersion := cCtx.String("relVersion")
									branch := cCtx.String("branch")
									defaultGameIniPath := cCtx.String("iniPath")
									_, err := projectVersion.SetVersionWithPerforce(devVersion, relVersion, branch, defaultGameIniPath)
									return err
								},
							},
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "setVersion",
								Aliases:  []string{"s"},
								Value:    "",
								Usage:    "Version to set in Project",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "iniPath",
								Aliases:  []string{"i"},
								Value:    "",
								Usage:    "file path to the DefaultGame ini file",
								Required: true,
							},
						},
						Action: func(cCtx *cli.Context) error {
							setVersion := cCtx.String("setVersion")
							defaultGameIniPath := cCtx.String("iniPath")
							return projectVersion.ManuallySetVersion(setVersion, defaultGameIniPath)
						},
					},
					{
						Name:  "get",
						Usage: "Prints the current project version",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "iniPath",
								Aliases:  []string{"i"},
								Value:    "",
								Usage:    "file path to the DefaultGame ini file",
								Required: true,
							},
						},
						Action: func(cCtx *cli.Context) error {
							defaultGameIniPath := cCtx.String("iniPath")
							_, err := projectVersion.GetVersion(defaultGameIniPath)
							return err
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
