package gofontweb

import "testing"

func TestFS(t *testing.T) {
	fsys := FS()
	f, err := fsys.Open("LICENSE.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
}
