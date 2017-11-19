package moduleutils

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

// ASTCache is a global cache of Program AST nodes keyed to a script filename
var ASTCache = newASTCache()

type astCache struct {
	m     sync.Mutex
	cache map[string]*cacheItem
}

type cacheItem struct {
	tree    *ast.Program
	modTime time.Time
}

func newASTCache() *astCache {
	return &astCache{
		cache: make(map[string]*cacheItem),
	}
}

func (c *astCache) GetTree(file string) (*ast.Program, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	c.m.Lock()
	defer c.m.Unlock()

	cachedItem, cached := c.cache[file]
	if cached && cachedItem.modTime.Equal(fileinfo.ModTime()) { // hit
		return cachedItem.tree, nil
	}

	// miss
	l, err := lexer.NewFile(file)
	if err != nil {
		return nil, err
	}

	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, errors.New(p.Errors()[0])
	}

	if cachedItem == nil {
		cachedItem = &cacheItem{}
	}
	cachedItem.tree = program
	cachedItem.modTime = fileinfo.ModTime()
	c.cache[file] = cachedItem

	return program, nil
}

func (c *astCache) ClearAll() {
	c.m.Lock()
	c.cache = make(map[string]*cacheItem)
	c.m.Unlock()
}

func (c *astCache) ClearOld(d time.Duration) {
	now := time.Now()
	c.m.Lock()
	for k, v := range c.cache {
		if v.modTime.Add(d).Before(now) {
			delete(c.cache, k)
		}
	}
	c.m.Unlock()
}

func (c *astCache) Remove(file string) {
	c.m.Lock()
	delete(c.cache, file)
	c.m.Unlock()
}
