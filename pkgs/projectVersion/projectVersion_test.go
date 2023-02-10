package projectVersion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	defaultGameIniPath = "./DefaultGame.ini"
)

func TestSetVersion(t *testing.T) {
	expectedVersion := "2.1.1"
	err := ManuallySetVersion(expectedVersion, defaultGameIniPath)
	assert.Nil(t, err)
}

func TestGetVersion(t *testing.T) {
	expectedVersion := "5.3.2"
	err := ManuallySetVersion(expectedVersion, defaultGameIniPath)
	assert.Nil(t, err)
	version, err := GetVersion(defaultGameIniPath)
	assert.Nil(t, err)
	assert.Equal(t, expectedVersion, version)
}

func TestPerforceRelVersion(t *testing.T) {
	version, err := SetVersionWithPerforce("1.0.0.0", "0.1.0.0", "rel0.1", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "1.0.0.0", version)

	version, err = SetVersionWithPerforce("1.0.0.0", "1.0.0.0", "rel0.1", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "1.0.0.0", version)
}

func TestPerforceDevVersion(t *testing.T) {
	version, err := SetVersionWithPerforce("1.2.0.5", "0.1.0.0", "dev1.2", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "1.2.0.6", version)

	version, err = SetVersionWithPerforce("1.2.0.5", "0.1.0.0", "dev2.3", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "2.3.0.1", version)
}

func TestPerforcePatchVersion(t *testing.T) {
	version, err := SetVersionWithPerforce("3.1.5.6", "3.1.5.6", "patch3.1.5", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "3.1.5.7", version)
}

func TestPerforceTaskVersion(t *testing.T) {
	version, err := SetVersionWithPerforce("1.1.0.33", "0.1.0.0", "task1.1", defaultGameIniPath)
	if err != nil {
		print(err.Error())
	}
	assert.Nil(t, err)
	assert.Equal(t, "1.1.0.34", version)
}
