package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OptionTestSuite struct {
	suite.Suite
}

func (o *OptionTestSuite) TestOption_IsSome_ReturnsFalseForNone() {
	assert.False(o.T(), None[any]().IsSome())
}

func (o *OptionTestSuite) TestOption_IsNone_ReturnsTrueForNone() {
	assert.True(o.T(), None[any]().IsNone())
}

func (o *OptionTestSuite) TestOption_Unwrap_PanicsIfNone() {
	assert.Panics(o.T(), func() { None[any]().Unwrap() })
}

func (o *OptionTestSuite) TestOption_Unwrap_ReturnsValueIfSome() {
	const value = "theValue:)"
	option := Some(value)

	unwrappedValue := option.Unwrap()

	assert.Equal(o.T(), value, unwrappedValue)
}

func TestOptionTestSuite(t *testing.T) {
	suite.Run(t, new(OptionTestSuite))
}
