/*
Package SemVer provides methods for parsing and comparing semantic versions
following the https://semver.org specification, as well as support for version
comparison operations with custom defined operator syntax.
*/
package SemVer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// See https://regex101.com/r/CkWF3o/1 for regex testing.
var opRe string = `[>|<]+=?`
var semverRe string = `(?:v)?([\d]+).([\d]+).([\d]+)(?:-((?:[.|-]?[\d\w]+)+))?(?:\+)?((?:[.|-]?[\d\w]+)+)?`
var re *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("^(%s)?%s$", opRe, semverRe))

var defaultConf config = config{
	ops: Operators{
		GT:  Operator(">"),
		GTE: Operator(">="),
		LT:  Operator("<"),
		LTE: Operator("<="),
	},
	re: re,
}

// Operators defines a set of operator syntax for semantic version comparisons.
type Operators struct {
	// GT is a greater than Operator.
	GT Operator
	// GTE is a greater than or equal to Operator.
	GTE Operator
	// LT is a less than Operator.
	LT Operator
	// LTE is a less than or equal to Operator.
	LTE Operator
}

type config struct {
	ops Operators
	re  *regexp.Regexp
}

/*
Config returns an intialized config object which can be passed to the SemVer.Get
method and define custom operator syntax and regex.

The regex string to parse the operators is combined with the SemVer regex. An invalid
regex string will result in a panic.

	conf := SemVer.Config(SemVer.Operators{
		GT: SemVer.Operator("+"),
		GTE: SemVer.Operator("+="),
		LT: SemVer.Operator("-"),
		LTE: SemVer.Operator("-="),
	}, `[\+|-]+=?`)

	v := SemVer("+=v1.0.0").Get(conf)
	fmt.Println(v.Compare(SemVer("v1.1.0").Get()))
*/
func Config(ops Operators, regex string) *config {
	regex = strings.TrimPrefix(regex, "^")
	regex = strings.TrimSuffix(regex, "$")
	return &config{
		ops: ops,
		re:  regexp.MustCompile(fmt.Sprintf("^(%s)?%s$", regex, semverRe)),
	}
}

/*
Operator is a comparison operator to be applied to a version.
*/
type Operator string

/*
String is a semantic version string with additional support for
an optional comparison Operator. For example:

>=v1.3.1

<=v3.0.0

>1.0.2

0.0.1-alpha

A String can be parsed to a Version for value parsing or
comparisons.

	ver := SemVer("v1.0.0").Get()
	i := ver.Compare(SemVer("v2.0.0").Get())
	fmt.Println(i == -1)

The "v" string character at the beginning of the version technically
does not conform to the https://semver.org specification, but is a
common convention when representing a semantic version in string format.
For this reason String treats the "v" in a version string as optional.
*/
type SemVer string

/*
version is a semantic version augmented with an Operator for fine grained
versioning rules and simple comparisons.

See https://semver.org/ for more info on semantic versioning and version
comparisons.
*/
type version struct {
	// Operator is an optional value for the set of comparison operators for the
	// version. See the SemVer.Operators for more info.
	Operator Operator
	// Major is the semantic major release version number. Must be a positive
	// integer.
	Major uint16
	// Minor is the semantic minor release version number. Must be a positive
	// integer.
	Minor uint16
	// Patch is the semantic patch release version number. Must be a positive
	// integer.
	Patch uint16
	// PreRelease is the string data contained after the '-' in a semantic
	// version string, but before the '+' denoting BuildMetadata. Can contain
	// only alphanumeric characters separated by '-' or '.'.
	PreRelease string
	// BuildMetadata is the string data after the '+' character in a semantic
	// version string. It can contain only alphanumeric characters separated by
	// a '-' or '.', and is not factored into version comparisons.
	BuildMetadata string

	conf *config
}

// ToSemVer returns the SemVer for the version.
func (v *version) ToSemVer() SemVer {
	var s strings.Builder
	s.WriteString(string(v.Operator))
	s.WriteString(v.String())
	return SemVer(s.String())
}

// String returns the version in semantic version string format.
//
// v{Major}.{Minor}.{Patch}-{PreRelease}+{BuildMetadata}
func (v *version) String() string {
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
OpCompare tests any current version Operator against the version param and
returns false if the passed version violates the Operator rule.

Version Operators on the version param are ignored.

Example:

	ver := SemVer(">=v1.2.3").Get()
	fail := ver.OpCompare(SemVer("1.0.0").Get())
	fmt.Println(fail)
*/
func (v *version) OpCompare(version *version) bool {
	i := v.Compare(version)

	switch v.Operator {
	case v.conf.ops.GTE:
		return i <= 0
	case v.conf.ops.GT:
		return i < 0
	case v.conf.ops.LTE:
		return i >= 0
	case v.conf.ops.LT:
		return i > 0
	default:
		return i == 0
	}
}

/*
Compare checks the two versions and returns 1 if the current version is greater than
the version param, -1 if the current version is less than the version param, and
0 if they are equal.

Comparison logic is implemented to the https://semver.org specification.
*/
func (v *version) Compare(version *version) int {
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
comparePreRelease is an internal method that evalutes only the current version
PreRelease value against the preRelease param. Similar to Compare, it returns
1 if the current version.PreRelease is greater than the preRelease param, -1 if
the current version.PreRelease is less than the preRelease param, and 0 if they
are equal.
*/
func (v *version) comparePreRelease(preRelease string) int {
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

// String returns the SemVer.String type as a string.
//
// {Operator}v{Major}.{Minor}.{Patch}-{PreRelease}+{BuildMetadata}
func (v SemVer) String() string {
	return string(v)
}

/*
Get returns a version from the String. Strings which are not
valid semantic versions will evaluate to v0.0.0.

Example:

	ver := SemVer("3.14.15").Get()
	fmt.Println("Major version:", ver.Major)
*/
func (v SemVer) Get(conf ...*config) *version {
	set := &defaultConf
	if conf != nil && conf[0] != nil {
		set = conf[0]
	}

	parts := set.re.FindStringSubmatch(v.String())
	if len(parts) != 7 {
		return &version{}
	}

	maj, _ := strconv.ParseInt(parts[2], 10, 16)
	min, _ := strconv.ParseInt(parts[3], 10, 16)
	patch, _ := strconv.ParseInt(parts[4], 10, 16)

	return &version{
		Operator:      Operator(parts[1]),
		Major:         uint16(maj),
		Minor:         uint16(min),
		Patch:         uint16(patch),
		PreRelease:    parts[5],
		BuildMetadata: parts[6],

		conf: set,
	}
}
