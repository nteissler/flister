package flister

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// the retriver interface represents ways of seeing
// if a search query matches an entry. Think of
// all the ways to search for files:
// exact name path, within folder, filename contains, etc...
type Retriever interface {
	Match(query, entry string) bool
}

type FilenameRegex struct {
	re   *regexp.Regexp
	once sync.Once
}
type AncestorOf struct{}
type FilenameExact struct{}
type FilenameContains struct{} // may remove this and just rely on Filename Regex

var UseFilenameRegex FilenameRegex
var UseAncestorOf AncestorOf
var UseFilenameExact FilenameExact
var UseFilenameContains FilenameContains

// a case sensitive regex comparison of a query to a filename
// more flexible than FilenameContains
// according to regex here: https://github.com/google/re2/wiki/Syntax
// expects a well formed regular expression, so you must do your
// own error checking for that. if regex is malformed, simply
// returns false
func (f *FilenameRegex) Match(query, entry string) bool {
	var err error
	f.once.Do(func() {
		f.re, err = regexp.Compile(query)
		fmt.Printf("regex match: %T: &p=%p i=%v\n", f, &f, f)
		fmt.Println("regex compiled")
	})
	if err != nil {
		log.Println("bad regex. err:", err)
		return false
	}
	loc := f.re.FindStringIndex(entry)
	if loc == nil {
		return false
	}
	return true
}

// a case insensitive seach for a match in all of a file's
// parent directories
func (_ AncestorOf) Match(directory, entry string) bool {
	dirs, _ := filepath.Split(entry)
	// strictly speaking, we could check the entire dirs string
	// for containing directory, but the for loop opens
	// up a little extra fleixibility  if needed in the future
	for _, dir := range strings.Split(dirs, "/") {
		if strings.Contains(strings.ToLower(dir), strings.ToLower(directory)) {
			return true
		}
	}
	return false
}

// a case insensitive search type for exact filename match
func (_ FilenameExact) Match(query, entry string) bool {
	return strings.ToLower(filepath.Base(query)) == strings.ToLower(filepath.Base(entry))
}

// a case insensitive search type for a part of the filename
func (_ FilenameContains) Match(query, entry string) bool {
	basename := strings.ToLower(filepath.Base(entry))
	return strings.Contains(basename, strings.ToLower(query))
}
