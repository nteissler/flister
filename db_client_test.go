package flister

import (
	"fmt"
	"github.com/matryer/filedb"
	"log"
	"os"
	"testing"
)

func TestCheckDB(t *testing.T) {
	// case where the file doesn't exist already
	db, err := checkDB()
	if err != nil {
		log.Fatalln(err)
	}
	c, err := db.C("testcollection")
	input := "here is some data"
	err = c.Insert([]byte(input))
	if err != nil {
		t.Error(err)
	}
	c.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
		if string(data) != input {
			t.Error("Didn't select right: %v, queried value != %v, input value", string(data), input)
		}
		return true, data, false
	})
}

func TestParseFileToDB(t *testing.T) {
	ParseFileToDB("/Users/nteissler/Documents/Sourcetree repos/SVN Finder/datatext/ABC150.svndb")
	if _, err := os.Stat("./database/ABC150.filedb"); os.IsNotExist(err) {
		t.Errorf("file not created correctly")
	}
	// if successful, delete the file and move cleanup test
	if err := os.Remove("./database/ABC150.filedb"); err != nil {
		t.Errorf("DB collection file not correctly deleted")
	}
}

func TestFindProgress(t *testing.T) {
	// Where the file exists
	db, _ := checkDB()
	makeTestCollection(db)
	query := "MetroM2_CAN1.h"
	received := make(chan string)
	progress := make(chan int)
	go FindProgress(query, UseFilenameExact, received, progress)
	go func() {
		for found := range received {
			fmt.Println(found)
		}
	}()

	for prog := range progress {
		fmt.Println(prog)
	}

	// Where the file doesn't exist
	// todo
	deleteTestCollection(db)
}

func makeTestCollection(db *filedb.DB) {
	ParseFileToDB("/Users/nteissler/Documents/Sourcetree repos/SVN Finder/datatext/ABC150.svndb")
}

func deleteTestCollection(db *filedb.DB) {
	os.Remove("./database/ABC150.filedb")
}
