package SemVer

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestSemVerParse(t *testing.T) {
	g := Goblin(t)
	g.Describe("SemVer string parsing to version", func() {
		g.It("Should parse all semantic version parts", func() {
			v := String(">=v1.2.3-pre+meta").Get()
			g.Assert(string(v.Operator)).Equal(">=")
			g.Assert(int(v.Major)).Equal(1)
			g.Assert(int(v.Minor)).Equal(2)
			g.Assert(int(v.Patch)).Equal(3)
			g.Assert(v.PreRelease).Equal("pre")
			g.Assert(v.BuildMetadata).Equal("meta")
		})

		g.It("Should return string value for SemVer", func() {
			v := String(">=v1.2.3-pre+meta")
			g.Assert(v.String()).Equal(">=v1.2.3-pre+meta")
		})

		g.It("Should return semantic string for version", func() {
			v := String(">=v1.2.3-pre+meta").Get()
			g.Assert(v.String()).Equal("v1.2.3-pre+meta")
		})

		g.It("Should return SemVer string for version", func() {
			v := String(">=v1.2.3-pre+meta").Get()
			g.Assert(v.ToString().String()).Equal(">=v1.2.3-pre+meta")
		})

		g.It("Should parse invalid semantic version to v0.0.0", func() {
			v := String("nosemver").Get()
			g.Assert(string(v.Operator)).Equal("")
			g.Assert(int(v.Major)).Equal(0)
			g.Assert(int(v.Minor)).Equal(0)
			g.Assert(int(v.Patch)).Equal(0)
			g.Assert(v.PreRelease).Equal("")
			g.Assert(v.BuildMetadata).Equal("")
		})
	})
}

func TestConfig(t *testing.T) {
	g := Goblin(t)
	g.Describe("Parse with custom config", func() {
		g.It("Should support custom Operator syntax", func() {
			conf := Config(Operators{
				GT:  Operator("+"),
				GTE: Operator("+="),
				LT:  Operator("-"),
				LTE: Operator("-="),
			}, `[\+|-]+=?`)

			v := String("+=v1.0.0").Get(conf)
			g.Assert(v.OpCompare(String("v1.1.0").Get())).IsTrue()
			g.Assert(v.OpCompare(String("v0.9.0").Get())).IsFalse()
			v = String("+v1.0.0").Get(conf)
			g.Assert(v.OpCompare(String("v1.1.0").Get())).IsTrue()
			g.Assert(v.OpCompare(String("v1.0.0").Get())).IsFalse()
			g.Assert(v.OpCompare(String("v0.9.0").Get())).IsFalse()
			v = String("-v1.0.0").Get(conf)
			g.Assert(v.OpCompare(String("v1.1.0").Get())).IsFalse()
			g.Assert(v.OpCompare(String("v1.0.0").Get())).IsFalse()
			g.Assert(v.OpCompare(String("v0.9.0").Get())).IsTrue()
			v = String("-=v1.0.0").Get(conf)
			g.Assert(v.OpCompare(String("v1.1.0").Get())).IsFalse()
			g.Assert(v.OpCompare(String("v1.0.0").Get())).IsTrue()
		})
	})
}

func TestOpCompare(t *testing.T) {
	g := Goblin(t)
	g.Describe("Version operator compare", func() {
		g.It("Evaluate greater than operator", func() {
			v := String(">v1.0.0").Get()
			v2 := String("v1.1.0").Get()
			v3 := String("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsFalse()
		})
		g.It("Evaluate greater than or equal to operator", func() {
			v := String(">=v1.0.0").Get()
			v2 := String("v1.1.0").Get()
			v3 := String("v1.0.0").Get()
			v4 := String("v0.9.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsTrue()
			g.Assert(v.OpCompare(v4)).IsFalse()
		})
		g.It("Evaluate less than operator", func() {
			v := String("<v1.0.0").Get()
			v2 := String("v0.9.0").Get()
			v3 := String("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsFalse()
		})
		g.It("Evaluate less than or equal to operator", func() {
			v := String("<=v1.0.0").Get()
			v2 := String("v1.1.0").Get()
			v3 := String("v1.0.0").Get()
			v4 := String("v0.9.0").Get()
			g.Assert(v.OpCompare(v2)).IsFalse()
			g.Assert(v.OpCompare(v3)).IsTrue()
			g.Assert(v.OpCompare(v4)).IsTrue()
		})
		g.It("Evaluate equality", func() {
			v := String("v1.0.0").Get()
			v2 := String("v1.1.0").Get()
			v3 := String("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsFalse()
			g.Assert(v.OpCompare(v3)).IsTrue()
		})
		g.It("Should handle invalid comparison operator", func() {
			v := String("~~v1.0.0").Get()
			v2 := String("v1.1.0").Get()
			g.Assert(v.OpCompare(v2)).IsFalse()
		})
	})
}

func TestCompare(t *testing.T) {
	g := Goblin(t)

	g.Describe("Version compare", func() {
		g.It("Major version", func() {
			v1 := String("v0.1.0").Get()
			v2 := String("v1.0.0").Get()
			v3 := String("v1.0.0").Get()
			v4 := String("v2.0.0").Get()

			g.Assert(v2.Compare(v3)).Equal(0)
			g.Assert(v2.Compare(v1)).Equal(1)
			g.Assert(v2.Compare(v4)).Equal(-1)
		})
		g.It("Minor version", func() {
			v1 := String("v0.0.0").Get()
			v2 := String("v0.1.0").Get()
			v3 := String("v0.1.0").Get()
			v4 := String("v0.2.0").Get()

			g.Assert(v2.Compare(v3)).Equal(0)
			g.Assert(v2.Compare(v1)).Equal(1)
			g.Assert(v2.Compare(v4)).Equal(-1)
		})
		g.It("Patch version", func() {
			v1 := String("v0.0.0").Get()
			v2 := String("v0.0.1").Get()
			v3 := String("v0.0.1").Get()
			v4 := String("v0.0.2").Get()

			g.Assert(v2.Compare(v3)).Equal(0)
			g.Assert(v2.Compare(v1)).Equal(1)
			g.Assert(v2.Compare(v4)).Equal(-1)
		})
	})
}

func TestComparePreRelease(t *testing.T) {
	g := Goblin(t)

	g.Describe("Compare pre release version", func() {
		g.It("Should return 0 for empty pre release info", func() {
			v := Version{PreRelease: ""}
			g.Assert(v.comparePreRelease("")).Equal(0)
		})
		g.It("Should return 1 for clean vs dirty build", func() {
			v := Version{PreRelease: ""}
			g.Assert(v.comparePreRelease("1")).Equal(1)
		})
		g.It("Should return -1 for dirty vs clean build", func() {
			v := Version{PreRelease: "alpha"}
			g.Assert(v.comparePreRelease("")).Equal(-1)
		})
		g.It("should handle alphabetical compare", func() {
			v := Version{PreRelease: "b"}
			g.Assert(v.comparePreRelease("a")).Equal(1)
			v = Version{PreRelease: "a"}
			g.Assert(v.comparePreRelease("b")).Equal(-1)
			v = Version{PreRelease: "b"}
			g.Assert(v.comparePreRelease("b")).Equal(0)
		})
		g.It("should handle numerical compare", func() {
			v := Version{PreRelease: "2"}
			g.Assert(v.comparePreRelease("1")).Equal(1)
			v = Version{PreRelease: "1"}
			g.Assert(v.comparePreRelease("2")).Equal(-1)
			v = Version{PreRelease: "1"}
			g.Assert(v.comparePreRelease("1")).Equal(0)
		})
		g.It("should handle '.' and '-' delimited data", func() {
			v := Version{PreRelease: "alpha.2"}
			g.Assert(v.comparePreRelease("alpha-1")).Equal(1)
			v = Version{PreRelease: "alpha.1"}
			g.Assert(v.comparePreRelease("alpha-2")).Equal(-1)
			v = Version{PreRelease: "alpha.2"}
			g.Assert(v.comparePreRelease("alpha-2")).Equal(0)
		})
		g.It("should handle mismatched sizes of delimited data", func() {
			v := Version{PreRelease: "alpha.2.1"}
			g.Assert(v.comparePreRelease("alpha-1")).Equal(1)
			v = Version{PreRelease: "alpha.1"}
			g.Assert(v.comparePreRelease("alpha-1.1")).Equal(-1)
		})
	})
}

func Example() {
	v := String("v3.14.15").Get()
	v2 := String("3.14.15").Get()
	fmt.Println(v.Compare(v2))
	// Output: 0
}

func Example_alt() {
	v := String("v3.14.15").Get()
	v2 := String("1.0.0").Get()
	fmt.Println(v.Compare(v2))
	// Output: 1
}

func Example_marshal() {
	// Because a SemVer.String is a just a string it can be
	// marshaled and unmarshaled to other data formats
	type Data struct {
		Version String `json:"version"`
	}

	jsonData := []byte(`{
		"version": ">=v3.14.15"
	}`)

	var data Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Version.Get().Minor)
	// Output: 14
}

func Example_opcompare() {
	v := String(">=v3.14.15").Get()
	v2 := String("3.14.16").Get()
	fmt.Println(v.OpCompare(v2))
	// Output: true
}

func ExampleString_String() {
	v := String(">=v3.14.15")
	fmt.Println(v.String())
	// Output: >=v3.14.15
}

func ExampleVersion_String() {
	v := String(">=v3.14.15").Get()
	fmt.Println(v.String())
	// Output: v3.14.15
}

func ExampleVersion() {
	s := String("v1.2.3")
	v := s.Get()
	fmt.Println(v.Patch)
	// Output: 3
}

func ExampleVersion_Compare_gt() {
	ver := String("v2.0.0").Get()
	i := ver.Compare(String("v1.0.0").Get())
	fmt.Println(i)
	// Output: 1
}

func ExampleVersion_Compare_lt() {
	ver := String("v1.0.0").Get()
	i := ver.Compare(String("v2.0.0").Get())
	fmt.Println(i)
	// Output: -1
}

func ExampleVersion_Compare_equal() {
	ver := String("v1.0.0").Get()
	i := ver.Compare(String("v1.0.0").Get())
	fmt.Println(i)
	// Output: 0
}

func ExampleVersion_OpCompare() {
	// By dropping any operator in the version OpCompare
	// will produce an equality check.
	ver := String(">=v1.0.0").Get()
	ok := ver.OpCompare(String("v1.0.0").Get())
	fmt.Println(ok)
	// Output: true
}

func ExampleVersion_OpCompare_equal() {
	ver := String("v1.0.0").Get()
	ok := ver.OpCompare(String("v1.0.0").Get())
	fmt.Println(ok)
	// Output: true
}

func ExampleConfig() {
	// Create a custom syntax for version comparison operators.
	conf := Config(Operators{
		GT:  Operator("+"),
		GTE: Operator("+="),
		LT:  Operator("-"),
		LTE: Operator("-="),
	}, `[\+|-]+=?`)

	v := String("+=v1.0.0").Get(conf)
	fmt.Println(v.OpCompare(String("v1.0.0").Get()))
	// Output: true
}

func ExampleConfig_gteorlte() {
	// Support only GTE or LTE comparisons.
	conf := Config(Operators{
		GT:  Operator(">="),
		GTE: Operator(">="),
		LT:  Operator("<="),
		LTE: Operator("<="),
	}, `[>|<]+=`)

	v := String(">=v1.0.0").Get(conf)
	fmt.Println(v.OpCompare(String("v1.0.0").Get()))
	// Output: true
}

func ExampleConfig_gte() {
	// Support only GTE comparisons with the ~ as the
	// identifying character.
	conf := Config(Operators{
		GT:  Operator("~"),
		GTE: Operator("~"),
		LT:  Operator("~"),
		LTE: Operator("~"),
	}, `~`)

	v := String("~v1.0.0").Get(conf)
	fmt.Println(v.OpCompare(String("v1.0.0").Get()))
	// Output: true
}
