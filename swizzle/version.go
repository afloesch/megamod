package swizzle

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const versionGTE versionOperator = ">="
const versionLTE versionOperator = "<="

var re *regexp.Regexp = regexp.MustCompile(`^([>|<]?=?)(?:v)([\d]+).([\d]+).([\d]+)(?:-(.*))?$`)

type versionOperator string

/*
SemVer is a semantic version string with additional support for
greater than or equal to ">=", or less than or equal to "<="
comparison operations. For example:

>=v1.3.1

<=v3.0.0

v1.0.2

v0.0.1-alpha

A SemVer string can be parsed to a Version for value parsing or
comparisons. Example:

	version := swizzle.SemVer("v1.0.0").Get()
*/
type SemVer string

/*
Version is a semantic version augmented with a VersionOperator
for finer grained versioning rules.
*/
type Version struct {
	Operator versionOperator `json:"operator,omitempty" yaml:"operator,omitempty"`
	Major    uint16          `json:"major" yaml:"major"`
	Minor    uint16          `json:"minor" yaml:"minor"`
	Patch    uint16          `json:"patch" yaml:"patch"`
	Build    string          `json:"build,omitempty" yaml:"build,omitempty"`
}

// SemVer returns the SemVer value for the Version.
func (v *Version) SemVer() SemVer {
	var s strings.Builder
	s.WriteString(string(v.Operator))
	s.WriteString(v.String())
	return SemVer(s.String())
}

// String returns the Version in semantic version string format.
func (v *Version) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("v%v.%v.%v", v.Major, v.Minor, v.Patch))
	if v.Build != "" {
		s.WriteString("-")
		s.WriteString(v.Build)
	}
	return s.String()
}

/*
OpCompare tests any current Version operator against the passed version and
returns false if the passed version violates the operator rule.
*/
func (v *Version) OpCompare(version *Version) bool {
	i := v.Compare(version)

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
Compare checks the two versions and returns 1 if current Version is greater than
the passed version, -1 if the current Version is less than the passed version,
and 0 if they are equal.
*/
func (v *Version) Compare(version *Version) int {
	if v.Major > version.Major {
		return 1
	}

	if v.Major < version.Major {
		return -1
	}

	if v.Minor > version.Minor {
		return 1
	}

	if v.Minor < version.Minor {
		return -1
	}

	if v.Patch > version.Patch {
		return 1
	}

	if v.Patch < version.Patch {
		return -1
	}

	if v.Build != version.Build {
		return 1
	}

	return 0
}

// String returns the SemVer string value.
func (s SemVer) String() string {
	return string(s)
}

// Get returns a Version from the SemVer string.
func (v SemVer) Get() *Version {
	parts := re.FindStringSubmatch(v.String())
	if len(parts) != 6 {
		return &Version{}
	}

	maj, _ := strconv.ParseInt(parts[2], 10, 16)
	min, _ := strconv.ParseInt(parts[3], 10, 16)
	patch, _ := strconv.ParseInt(parts[4], 10, 16)

	return &Version{
		Operator: versionOperator(parts[1]),
		Major:    uint16(maj),
		Minor:    uint16(min),
		Patch:    uint16(patch),
		Build:    parts[5],
	}
}
