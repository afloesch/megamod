package internal

import (
	"fmt"
	"regexp"
	"strconv"
)

const versionGTE versionOperator = ">="
const versionLTE versionOperator = "<="

var re *regexp.Regexp = regexp.MustCompile(`^([>|<]?=?)(?:v)([\d]+).([\d]+).([\d]+)(?:-(.*))?$`)

type versionOperator string

type SemVer string

// Version is a semantic version augmented with a VersionOperator
// for finer grained versioning rules.
type Version struct {
	Operator versionOperator `json:"operator,omitempty" yaml:"operator,omitempty"`
	Major    int             `json:"major" yaml:"major"`
	Minor    int             `json:"minor" yaml:"minor"`
	Patch    int             `json:"patch" yaml:"patch"`
	Build    string          `json:"build,omitempty" yaml:"build,omitempty"`
}

func (v *Version) String() string {
	return fmt.Sprintf("v%v.%v.%v", v.Major, v.Minor, v.Patch)
}

/*
OpCompare tests the Version Operator against the passed version (d) and
returns true if d passes the operator rule.
*/
func (v *Version) OpCompare(d *Version) bool {
	i := v.Compare(d)

	switch v.Operator {
	case versionGTE:
		return i <= 0
	case versionLTE:
		return i >= 0
	default:
		return i == 0
	}
}

/*
Compare checks the two versions and returns 1 if v is greater than d, -1 if
v is less than d, and 0 if they are equal.
*/
func (v *Version) Compare(d *Version) int {
	if v.Major > d.Major {
		return 1
	}

	if v.Major < d.Major {
		return -1
	}

	if v.Minor > d.Minor {
		return 1
	}

	if v.Minor < d.Minor {
		return -1
	}

	if v.Patch > d.Patch {
		return 1
	}

	if v.Patch < d.Patch {
		return -1
	}

	if v.Build != d.Build {
		return 1
	}

	return 0
}

func (s SemVer) String() string {
	return string(s)
}

// Version parses a semantic version string to a Version object.
func (v SemVer) Get() *Version {
	parts := re.FindStringSubmatch(v.String())
	if len(parts) != 6 {
		return nil
	}

	maj, _ := strconv.Atoi(parts[2])
	min, _ := strconv.Atoi(parts[3])
	patch, _ := strconv.Atoi(parts[4])
	return &Version{
		Operator: versionOperator(parts[1]),
		Major:    maj,
		Minor:    min,
		Patch:    patch,
		Build:    parts[5],
	}
}
