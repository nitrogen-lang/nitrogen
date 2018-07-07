package moduleutils

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/compiler"
	"github.com/nitrogen-lang/nitrogen/src/compiler/marshal"
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

func (c *codeBlockCache) GetBlock(file string, name string) (*compiler.CodeBlock, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	c.m.Lock()

	cachedItem, cached := c.cache[file]
	if cached && cachedItem.modTime.Equal(fileinfo.ModTime()) { // hit
		c.m.Unlock()
		return cachedItem.block, nil
	}

	// miss
	if cachedItem == nil {
		cachedItem = &cbCacheItem{}
	}

	if filepath.Ext(file) == ".nib" {
		srcfile := file[:len(file)-1]
		srcinfo, err := os.Stat(srcfile)
		if err != nil {
			c.m.Unlock()
			if os.IsNotExist(err) {
				return nil, errors.New("source file for compiled nib not found")
			}
			return nil, err
		}

		code, modinfo, err := marshal.ReadFile(file)
		if err != nil {
			c.m.Unlock()
			if marshal.IsErrVersion(err) {
				return c.GetBlock(srcfile, name)
			}
			return nil, err
		}

		if !modinfo.ModTime.Equal(srcinfo.ModTime().Round(time.Second)) {
			c.m.Unlock()
			return c.GetBlock(srcfile, name)
		}

		cachedItem.block = code
	} else {
		program, err := ASTCache.GetTree(file)
		if err != nil {
			c.m.Unlock()
			return nil, err
		}
		cachedItem.block = compiler.Compile(program, name)

		ext := path.Ext(file)
		outfile := file[:len(file)-len(ext)] + ".nib"
		marshal.WriteFile(outfile, cachedItem.block, fileinfo.ModTime())
	}
	cachedItem.modTime = fileinfo.ModTime()
	c.cache[file] = cachedItem

	c.m.Unlock()
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
