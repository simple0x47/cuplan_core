package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// assure clean after time cleans expired items

type CacheTestSuite struct {
	suite.Suite
	cache *Cache
}

func (c *CacheTestSuite) SetupTest() {
	c.cache = NewCache(time.Hour)
}

func (c *CacheTestSuite) TestCache_CleanExpiredItemsAfterTime() {
	cache := NewCache(time.Second)
	key := "expire"
	const value = true

	cache.Set(key, value, time.Second)
	time.Sleep(time.Second * 2)

	assert.False(c.T(), cache.HasKey(key), "Cached value has not been removed after cleaning time has elapsed.")
}

func (c *CacheTestSuite) TestCache_Get_ReturnsPreviouslySetValue() {
	key := "example"
	const value = "abcd"
	c.cache.Set(key, value, time.Hour)

	cachedValue := c.cache.Get(key)

	assert.Equal(c.T(), value, cachedValue.Unwrap())
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
