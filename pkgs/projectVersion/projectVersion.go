package projectVersion

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

func ManuallySetVersion(setVersion string, defaultGameIniPath string) error {
	// loading ini file
	ini.PrettyFormat = false
	cfg, err := ini.ShadowLoad(defaultGameIniPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("error loading ini file: %s", err), 1)
	}

	// updating game version in ini
	cfg.Section("/Script/EngineSettings.GeneralProjectSettings").Key("ProjectVersion").SetValue(setVersion)
	err = cfg.SaveTo(defaultGameIniPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("error saving ini file: %s", err), 1)
	}

	return nil
}

func SetVersionWithPerforce(devVersion string, relVersion string, branch string, defaultGameIniPath string) (string, error) {
	currentVersion := devVersion

	// if the new rel version has a different major/minor version than the current rel version, copy number from dev
	if strings.Contains(branch, "rel") {
		relVerSplit := strings.Split(strings.TrimPrefix(branch, "rel"), ".")
		curDevVerSplit := strings.Split(devVersion, ".")

		majorDev, _ := strconv.Atoi(curDevVerSplit[0])
		minorDev, _ := strconv.Atoi(curDevVerSplit[1])
		majorRel, _ := strconv.Atoi(relVerSplit[0])
		minorRel, _ := strconv.Atoi(relVerSplit[1])

		if majorDev != majorRel || minorRel != minorDev {
			currentVersion = devVersion
		} else {
			currentVersion = relVersion
		}
	}

	semVerSplit := strings.Split(currentVersion, ".")

	major, _ := strconv.Atoi(semVerSplit[0])
	minor, _ := strconv.Atoi(semVerSplit[1])
	patch, _ := strconv.Atoi(semVerSplit[2])
	build, _ := strconv.Atoi(semVerSplit[3])

	// hard setting the version tag if the branch is dev
	if strings.Contains(branch, "dev") {
		devVerSplit := strings.Split(strings.TrimPrefix(branch, "dev"), ".")
		newMajor, _ := strconv.Atoi(devVerSplit[0])
		newMinor, _ := strconv.Atoi(devVerSplit[1])

		// Overriding the major/minor version if it's not currently set
		if newMajor != major || minor != newMinor {
			major = newMajor
			minor = newMinor
			patch = 0
			build = 1
		} else {
			build += 1
		}
	}

	// we should never be doing any dev on main branch.  So we just return the current version at the end

	// patch branches are for hotfixes applied directly to rel branch
	if strings.Contains(branch, "patch") {
		patchVerSplit := strings.Split(strings.TrimPrefix(branch, "patch"), ".")

		newMajor, _ := strconv.Atoi(patchVerSplit[0])
		newMinor, _ := strconv.Atoi(patchVerSplit[1])
		newPatch, _ := strconv.Atoi(patchVerSplit[2])

		// Overriding the major/minor/patch version if it's not currently set to this branch
		if newMajor != major || minor != newMinor || newPatch != patch {
			major = newMajor
			minor = newMinor
			patch = newPatch
			build = 1
		} else {
			build += 1
		}
	}

	// for task branches we want to increment the build number
	if strings.Contains(branch, "task") {
		build += 1
	}

	newVersionNumber := fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build)

	// loading ini file
	ini.PrettyFormat = false
	cfg, err := ini.ShadowLoad(defaultGameIniPath)
	if err != nil {
		return "", cli.Exit(fmt.Errorf("error loading ini: %s", err), 1)
	}

	// updating game version in ini
	cfg.Section("/Script/EngineSettings.GeneralProjectSettings").Key("ProjectVersion").SetValue(newVersionNumber)
	err = cfg.SaveTo(defaultGameIniPath)
	if err != nil {
		return "", cli.Exit(fmt.Errorf("error saving ini: %s", err), 1)
	}

	// printing the string to standard output
	w := os.Stdout
	if currentVersion == relVersion {
		_, err = fmt.Fprintf(w, "##teamcity[setParameter name='REL_VERSION' value='%s']", newVersionNumber)
	} else {
		_, err = fmt.Fprintf(w, "##teamcity[setParameter name='DEV_VERSION' value='%s']", newVersionNumber)
	}

	return newVersionNumber, nil
}

func GetVersion(defaultGameIniPath string) (string, error) {
	ini.PrettyFormat = false
	cfg, err := ini.ShadowLoad(defaultGameIniPath)
	if err != nil {
		return "", cli.Exit(fmt.Errorf("error loading ini: %s", err), 1)
	}

	version := cfg.Section("/Script/EngineSettings.GeneralProjectSettings").Key("ProjectVersion").Value()

	fmt.Printf(version)
	return version, nil
}
