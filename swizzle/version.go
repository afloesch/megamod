package swizzle

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// See https://regex101.com/r/CkWF3o/1 for regex testing.
var opRe string = `[>|<]?=?`
var svRe string = `(?:v)?([\d]+).([\d]+).([\d]+)(?:-((?:[.|-]?[\d\w]+)+))?(?:\+)?((?:[.|-]?[\d\w]+)+)?`
var re *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("^(%s)?%s$", opRe, svRe))

var defaultOpSettings OpSettings = OpSettings{
	GT:  Operator(">"),
	GTE: Operator(">="),
	LT:  Operator("<"),
	LTE: Operator("<="),
	re:  re,
}

// OpSettings defines a custom set of operators for
// semantic version parsing comparisons.
type OpSettings struct {
	// GT is a greater than Operator.
	GT Operator
	// GTE is a greater than or equal to Operator.
	GTE Operator
	// LT is a less than Operator.
	LT Operator
	// LTE is a less than or equal to Operator.
	LTE Operator
	// Regex is the Go comparible regex string which matches
	// the defined Operators. Should exclude any capture group.
	RegEx string

	re *regexp.Regexp
}

/*
Operator is a comparison operator to be applied to a Version.
*/
type Operator string

/*
SemVer is a semantic version string with additional support for
an optional comparison Operator. For example:

- >=v1.3.1

- <=v3.0.0

- >1.0.2

- 0.0.1-alpha

A SemVer string can be parsed to a Version for value parsing or
comparisons. Example:

	version := swizzle.SemVer("v1.0.0").Get(nil)

The "v" string character at the beginning of the version technically
does not conform to the https://semver.org specification, but is a
common convention when representing a semantic version in string format.
For this reason SemVer treats the "v" in a version string as optional.
*/
type SemVer string

/*
Version is a semantic version augmented with an Operator for fine grained
versioning rules and simple comparisons.

See https://semver.org/ for more info on semantic versioning and version
comparisons.
*/
type Version struct {
	Operator      Operator
	Major         uint16
	Minor         uint16
	Patch         uint16
	PreRelease    string
	BuildMetadata string

	settings *OpSettings
}

// SemVer returns the SemVer value for the Version.
func (v *Version) SemVer() SemVer {
	var s strings.Builder
	s.WriteString(string(v.Operator))
	s.WriteString(v.String())
	return SemVer(s.String())
}

// String returns the Version in semantic version string format.
//
// v{Major}.{Minor}.{Patch}-{PreRelease}+{BuildMetadata}
func (v *Version) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("v%v.%v.%v", v.Major, v.Minor, v.Patch))
	if v.PreRelease != "" {
		s.WriteString("-")
		s.WriteString(v.PreRelease)
	}
	if v.BuildMetadata != "" {
		s.WriteString("+")
		s.WriteString(v.BuildMetadata)
	}
	return s.String()
}

/*
OpCompare tests any current Version Operator against the version param and
returns false if the passed version violates the Operator rule.

Version Operators on the version param are ignored.

Example:

	ver := swizzle.SemVer(">=v1.2.3").Get(nil)
	fail := ver.OpCompare(swizzle.SemVer("1.0.0").Get(nil))
	fmt.Println(fail)
*/
func (v *Version) OpCompare(version *Version) bool {
	i := v.Compare(version)

	switch v.Operator {
	case v.settings.GTE:
		return i <= 0
	case v.settings.GT:
		return i < 0
	case v.settings.LTE:
		return i >= 0
	case v.settings.LT:
		return i > 0
	default:
		return i == 0
	}
}

/*
Compare checks the two versions and returns 1 if the current Version is greater than
the version param, -1 if the current Version is less than the version param, and
0 if they are equal.

Comparison logic is implemented to the https://semver.org specification.
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

	return v.comparePreRelease(version.PreRelease)
}

/*
comparePreRelease is an internal method that evalutes only the current Version
PreRelease value against the preRelease param. Similar to Compare, it returns
1 if the current Version.PreRelease is greater than the preRelease param, -1 if
the current Version.PreRelease is less than the preRelease param, and 0 if they
are equal.
*/
func (v *Version) comparePreRelease(preRelease string) int {
	if v.PreRelease == "" && preRelease == "" {
		return 0
	}

	if v.PreRelease == "" && preRelease != "" {
		return 1
	}

	if v.PreRelease != "" && preRelease == "" {
		return -1
	}

	// split pre release string parts
	vp := strings.FieldsFunc(v.PreRelease, splitRelease)
	versionp := strings.FieldsFunc(preRelease, splitRelease)

	// fill missing values
	if len(vp) < len(versionp) {
		for i := len(vp); i < len(versionp); i++ {
			vp = append(vp, "")
		}
	}

	if len(vp) > len(versionp) {
		for i := len(versionp); i < len(vp); i++ {
			versionp = append(versionp, "")
		}
	}

	// compare all pre release parts
	for i, v := range vp {
		if v == versionp[i] {
			continue
		} else if v > versionp[i] {
			return 1
		} else {
			return -1
		}
	}

	return 0
}

// splitRelease is a helper method to split a string on '-' or '.' characters.
func splitRelease(r rune) bool {
	return r == '-' || r == '.'
}

// String returns the SemVer string value.
func (s SemVer) String() string {
	return string(s)
}

/*
Get returns a Version from the SemVer string. SemVer strings which are not
valid semantic version strings will evaluate to v0.0.0.

Example:

	version := swizzle.SemVer("3.14.15").Get(nil)
	fmt.Println("Major version:", version.Major)
*/
func (v SemVer) Get(settings *OpSettings) *Version {
	set := &defaultOpSettings
	if settings != nil && settings.RegEx != "" {
		set = settings
		set.re = regexp.MustCompile(fmt.Sprintf("^(%s)?%s$", settings.RegEx, svRe))
	}

	parts := set.re.FindStringSubmatch(v.String())
	if len(parts) != 7 {
		return &Version{}
	}

	maj, _ := strconv.ParseInt(parts[2], 10, 16)
	min, _ := strconv.ParseInt(parts[3], 10, 16)
	patch, _ := strconv.ParseInt(parts[4], 10, 16)

	return &Version{
		Operator:      Operator(parts[1]),
		Major:         uint16(maj),
		Minor:         uint16(min),
		Patch:         uint16(patch),
		PreRelease:    parts[5],
		BuildMetadata: parts[6],

		settings: set,
	}
}
