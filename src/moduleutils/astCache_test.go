package moduleutils

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestBasicCache(t *testing.T) {
	ASTCache.Clear()

	cache1, err := ASTCache.GetTree("./testdata/cache1.ni")
	if err != nil {
		t.Fatal(err)
	}

	cache1_1, err := ASTCache.GetTree("./testdata/cache1.ni")
	if err != nil {
		t.Fatal(err)
	}

	if cache1 != cache1_1 {
		t.Fatal("Returned trees are not the same object")
	}
}

func TestMultipleCache(t *testing.T) {
	ASTCache.Clear()

	cache1, err := ASTCache.GetTree("./testdata/cache1.ni")
	if err != nil {
		t.Fatal(err)
	}

	// Same tree but different script name
	cache2, err := ASTCache.GetTree("./testdata/cache2.ni")
	if err != nil {
		t.Fatal(err)
	}

	cache1_1, err := ASTCache.GetTree("./testdata/cache1.ni")
	if err != nil {
		t.Fatal(err)
	}

	cache2_1, err := ASTCache.GetTree("./testdata/cache2.ni")
	if err != nil {
		t.Fatal(err)
	}

	if cache1 != cache1_1 {
		t.Fatal("Returned cache1 trees are not the same object")
	}

	if cache2 != cache2_1 {
		t.Fatal("Returned cache2 trees are not the same object")
	}

	if cache1_1 == cache2_1 {
		t.Fatal("Returned cache1 and cache2 trees are the same object")
	}
}

const cacheMissTestScript = `let str = "Hello, world!"
println(str)
`

func TestCacheMiss(t *testing.T) {
	ASTCache.Clear()
	// Copy a test data script tp play with
	copyFileContents("./testdata/cache1.ni", "./testdata/cache1-1.ni")
	defer os.Remove("./testdata/cache1-1.ni")

	// Get first cache
	cache1, err := ASTCache.GetTree("./testdata/cache1-1.ni")
	if err != nil {
		t.Fatal(err)
	}

	// Modify the file
	if err := ioutil.WriteFile("./testdata/cache1-1.ni", []byte(cacheMissTestScript), 0644); err != nil {
		t.Fatal(err)
	}

	// Get second cache, hopefully a miss and a new tree
	cache1_1, err := ASTCache.GetTree("./testdata/cache1-1.ni")
	if err != nil {
		t.Fatal(err)
	}

	if cache1 == cache1_1 {
		t.Fatal("Cache didn't miss when it should have")
	}
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
