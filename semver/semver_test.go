package SemVer

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestSemVerParse(t *testing.T) {
	g := Goblin(t)
	g.Describe("SemVer string parsing to version", func() {
		g.It("Should parse all semantic version parts", func() {
			v := SemVer(">=v1.2.3-pre+meta").Get()
			g.Assert(string(v.Operator)).Equal(">=")
			g.Assert(int(v.Major)).Equal(1)
			g.Assert(int(v.Minor)).Equal(2)
			g.Assert(int(v.Patch)).Equal(3)
			g.Assert(v.PreRelease).Equal("pre")
			g.Assert(v.BuildMetadata).Equal("meta")
		})

		g.It("Should return string value for SemVer", func() {
			v := SemVer(">=v1.2.3-pre+meta")
			g.Assert(v.String()).Equal(">=v1.2.3-pre+meta")
		})

		g.It("Should return semantic string for version", func() {
			v := SemVer(">=v1.2.3-pre+meta").Get()
			g.Assert(v.String()).Equal("v1.2.3-pre+meta")
		})

		g.It("Should return SemVer string for version", func() {
			v := SemVer(">=v1.2.3-pre+meta").Get()
			g.Assert(v.ToSemVer().String()).Equal(">=v1.2.3-pre+meta")
		})

		g.It("Should parse invalid semantic version to v0.0.0", func() {
			v := SemVer("nosemver").Get()
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

			v := SemVer("+=v1.0.0").Get(conf)
			g.Assert(v.OpCompare(SemVer("v1.1.0").Get())).IsTrue()
			g.Assert(v.OpCompare(SemVer("v0.9.0").Get())).IsFalse()
			v = SemVer("+v1.0.0").Get(conf)
			g.Assert(v.OpCompare(SemVer("v1.1.0").Get())).IsTrue()
			g.Assert(v.OpCompare(SemVer("v1.0.0").Get())).IsFalse()
			g.Assert(v.OpCompare(SemVer("v0.9.0").Get())).IsFalse()
			v = SemVer("-v1.0.0").Get(conf)
			g.Assert(v.OpCompare(SemVer("v1.1.0").Get())).IsFalse()
			g.Assert(v.OpCompare(SemVer("v1.0.0").Get())).IsFalse()
			g.Assert(v.OpCompare(SemVer("v0.9.0").Get())).IsTrue()
			v = SemVer("-=v1.0.0").Get(conf)
			g.Assert(v.OpCompare(SemVer("v1.1.0").Get())).IsFalse()
			g.Assert(v.OpCompare(SemVer("v1.0.0").Get())).IsTrue()
		})
	})
}

func TestOpCompare(t *testing.T) {
	g := Goblin(t)
	g.Describe("Version operator compare", func() {
		g.It("Evaluate greater than operator", func() {
			v := SemVer(">v1.0.0").Get()
			v2 := SemVer("v1.1.0").Get()
			v3 := SemVer("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsFalse()
		})
		g.It("Evaluate greater than or equal to operator", func() {
			v := SemVer(">=v1.0.0").Get()
			v2 := SemVer("v1.1.0").Get()
			v3 := SemVer("v1.0.0").Get()
			v4 := SemVer("v0.9.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsTrue()
			g.Assert(v.OpCompare(v4)).IsFalse()
		})
		g.It("Evaluate less than operator", func() {
			v := SemVer("<v1.0.0").Get()
			v2 := SemVer("v0.9.0").Get()
			v3 := SemVer("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsTrue()
			g.Assert(v.OpCompare(v3)).IsFalse()
		})
		g.It("Evaluate less than or equal to operator", func() {
			v := SemVer("<=v1.0.0").Get()
			v2 := SemVer("v1.1.0").Get()
			v3 := SemVer("v1.0.0").Get()
			v4 := SemVer("v0.9.0").Get()
			g.Assert(v.OpCompare(v2)).IsFalse()
			g.Assert(v.OpCompare(v3)).IsTrue()
			g.Assert(v.OpCompare(v4)).IsTrue()
		})
		g.It("Evaluate equality", func() {
			v := SemVer("v1.0.0").Get()
			v2 := SemVer("v1.1.0").Get()
			v3 := SemVer("v1.0.0").Get()
			g.Assert(v.OpCompare(v2)).IsFalse()
			g.Assert(v.OpCompare(v3)).IsTrue()
		})
	})
}

func TestCompare(t *testing.T) {
	g := Goblin(t)

	g.Describe("Version compare", func() {
		g.It("Major version", func() {
			v1 := SemVer("v0.1.0").Get()
			v2 := SemVer("v1.0.0").Get()
			v3 := SemVer("v1.0.0").Get()
			v4 := SemVer("v2.0.0").Get()

			g.Assert(v2.Compare(v3)).Equal(0)
			g.Assert(v2.Compare(v1)).Equal(1)
			g.Assert(v2.Compare(v4)).Equal(-1)
		})
		g.It("Minor version", func() {
			v1 := SemVer("v0.0.0").Get()
			v2 := SemVer("v0.1.0").Get()
			v3 := SemVer("v0.1.0").Get()
			v4 := SemVer("v0.2.0").Get()

			g.Assert(v2.Compare(v3)).Equal(0)
			g.Assert(v2.Compare(v1)).Equal(1)
			g.Assert(v2.Compare(v4)).Equal(-1)
		})
		g.It("Patch version", func() {
			v1 := SemVer("v0.0.0").Get()
			v2 := SemVer("v0.0.1").Get()
			v3 := SemVer("v0.0.1").Get()
			v4 := SemVer("v0.0.2").Get()

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
			v := version{PreRelease: ""}
			g.Assert(v.comparePreRelease("")).Equal(0)
		})
		g.It("Should return 1 for clean vs dirty build", func() {
			v := version{PreRelease: ""}
			g.Assert(v.comparePreRelease("1")).Equal(1)
		})
		g.It("Should return -1 for dirty vs clean build", func() {
			v := version{PreRelease: "alpha"}
			g.Assert(v.comparePreRelease("")).Equal(-1)
		})
		g.It("should handle alphabetical compare", func() {
			v := version{PreRelease: "b"}
			g.Assert(v.comparePreRelease("a")).Equal(1)
			v = version{PreRelease: "a"}
			g.Assert(v.comparePreRelease("b")).Equal(-1)
			v = version{PreRelease: "b"}
			g.Assert(v.comparePreRelease("b")).Equal(0)
		})
		g.It("should handle numerical compare", func() {
			v := version{PreRelease: "2"}
			g.Assert(v.comparePreRelease("1")).Equal(1)
			v = version{PreRelease: "1"}
			g.Assert(v.comparePreRelease("2")).Equal(-1)
			v = version{PreRelease: "1"}
			g.Assert(v.comparePreRelease("1")).Equal(0)
		})
		g.It("should handle '.' and '-' delimited data", func() {
			v := version{PreRelease: "alpha.2"}
			g.Assert(v.comparePreRelease("alpha-1")).Equal(1)
			v = version{PreRelease: "alpha.1"}
			g.Assert(v.comparePreRelease("alpha-2")).Equal(-1)
			v = version{PreRelease: "alpha.2"}
			g.Assert(v.comparePreRelease("alpha-2")).Equal(0)
		})
		g.It("should handle mismatched sizes of delimited data", func() {
			v := version{PreRelease: "alpha.2.1"}
			g.Assert(v.comparePreRelease("alpha-1")).Equal(1)
			v = version{PreRelease: "alpha.1"}
			g.Assert(v.comparePreRelease("alpha-1.1")).Equal(-1)
		})
	})
}
