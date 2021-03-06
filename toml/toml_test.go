package toml_test

import (
	//	"reflect"

	"github.com/megamsys/megamd/toml"
	"gopkg.in/check.v1"
)

type S struct{}

var _ = check.Suite(&S{})

// Ensure that megabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_MB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("200m"))
	c.Assert(err, check.NotNil)
	c.Assert(sb, check.Not(check.Equals), 200*(1<<20))
}

// Ensure that gigabyte sizes can be parsed.
func (s *S) TestSize_UnmarshalText_GB(c *check.C) {
	var sb toml.Size
	err := sb.UnmarshalText([]byte("10g"))
	c.Assert(err, check.NotNil)
	c.Assert(sb, check.Not(check.Equals), (10 * (1 << 30)))

}
