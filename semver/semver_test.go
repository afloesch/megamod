package SemVer

import (
	"testing"

	. "github.com/franela/goblin"
)

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
