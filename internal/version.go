package internal

import (
	"regexp"
	"strconv"
)

const VersionGTE VersionOperator = ">="
const VersionLTE VersionOperator = "<="

var re *regexp.Regexp = regexp.MustCompile(`^([>|<]?=?)(?:v)([\d]+).([\d]+).([\d]+)(?:-(.*))?$`)

type VersionOperator string

// Version is a semantic version augmented with a VersionOperator
// for finer grained versioning rules.
type Version struct {
	Operator VersionOperator
	Major    int
	Minor    int
	Patch    int
	Build    string
}

/*
OpCompare tests the Version Operator against the passed version (d) and
returns true if d passes the operator rule.
*/
func (v *Version) OpCompare(d *Version) bool {
	i := v.Compare(d)

	switch v.Operator {
	case VersionGTE:
		return i <= 0
	case VersionLTE:
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

// SemVersion parses a semantic version string to a Version object.
func SemVersion(v string) *Version {
	parts := re.FindStringSubmatch(v)
	if len(parts) != 6 {
		return nil
	}

	maj, _ := strconv.Atoi(parts[2])
	min, _ := strconv.Atoi(parts[3])
	patch, _ := strconv.Atoi(parts[4])
	return &Version{
		Operator: VersionOperator(parts[1]),
		Major:    maj,
		Minor:    min,
		Patch:    patch,
		Build:    parts[5],
	}
}
