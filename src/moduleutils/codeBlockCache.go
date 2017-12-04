package moduleutils

import (
	"os"
	"sync"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
)

// CodeBlockCache is a global cache of Code Blocks keyed to a script filename
var (
	CodeBlockCache = newCodeBlockCache()
)

type codeBlockCache struct {
	m     sync.Mutex
	cache map[string]*cbCacheItem
}

type cbCacheItem struct {
	block   *compiler.CodeBlock
	modTime time.Time
}

func newCodeBlockCache() *codeBlockCache {
	return &codeBlockCache{
		cache: make(map[string]*cbCacheItem),
	}
}

func (c *codeBlockCache) GetBlock(file string) (*compiler.CodeBlock, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	c.m.Lock()
	defer c.m.Unlock()

	cachedItem, cached := c.cache[file]
	if cached && cachedItem.modTime.Equal(fileinfo.ModTime()) { // hit
		return cachedItem.block, nil
	}

	// miss
	program, err := ASTCache.GetTree(file)
	if err != nil {
		return nil, err
	}

	if cachedItem == nil {
		cachedItem = &cbCacheItem{}
	}
	cachedItem.block = compiler.Compile(program, "__main")
	cachedItem.modTime = fileinfo.ModTime()
	c.cache[file] = cachedItem

	return cachedItem.block, nil
}

func (c *codeBlockCache) ClearAll() {
	c.m.Lock()
	c.cache = make(map[string]*cbCacheItem)
	c.m.Unlock()
}

func (c *codeBlockCache) ClearOld(d time.Duration) {
	now := time.Now()
	c.m.Lock()
	for k, v := range c.cache {
		if v.modTime.Add(d).Before(now) {
			delete(c.cache, k)
		}
	}
	c.m.Unlock()
}

func (c *codeBlockCache) Remove(file string) {
	c.m.Lock()
	delete(c.cache, file)
	c.m.Unlock()
}
