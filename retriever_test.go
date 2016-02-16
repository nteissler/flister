package flister

import (
	"testing"
)

func TestRegex(t *testing.T) {
	match := UseFilenameRegex.Match(".*hello.*", "/path/to/hello.txt")
	if !match {
		t.Error("First regex test failed")
	}

	match = UseFilenameRegex.Match("not here", "/path/to/nothing")
	if match {
		t.Error("Second regex failed")
	}
}

func TestAncestor(t *testing.T) {
	entry := "/Users/nteissler/Dropbox/Camera Uploads/2015-03-02 05.28.39.png"
	directory := "dropboX"
	match := UseAncestorOf.Match(directory, entry)
	if !match {
		t.Error("TestAncestor should have found a match on test 1")
	}
	directory = "svn"
	match = UseAncestorOf.Match(directory, entry)
	if match {
		t.Error("TestAncestor should have not found a match")
	}
}

func TestExact(t *testing.T) {
	query := "NICK.pdf"
	entry := "/Users/nteissler/nick.pdf"
	match := UseFilenameExact.Match(query, entry)
	if !match {
		t.Error("TestExact first test should have matched")
	}
	query = "poo.pdf"
	match = UseFilenameExact.Match(query, entry)
	if match {
		t.Error("TestExact second test should not have matched")
	}
}

func TestContains(t *testing.T) {
	query := "ovc"
	entry := "/Users/nteissler/OVC battery.pdf"
	match := UseFilenameContains.Match(query, entry)
	if !match {
		t.Error("TestExact first test should have matched")
	}
	query = "poo.pdf"
	match = UseFilenameContains.Match(query, entry)
	if match {
		t.Error("TestExact second test should not have matched")
	}
}
