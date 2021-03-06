package flister

import (
	"bufio"
	"fmt"
	"github.com/matryer/filedb"
	"github.com/nteissler/stringfixer"
	"log"
	"os"
)

const Dbpath = "./database"

type Client struct {
	Matches  chan []byte
	Progress chan int
	Done     chan struct{}
}

func ParseFileToDB(filename string) {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()
	// check file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fatalErr = err
		return
	}

	// create the database if it doesn't exist
	db, err := checkDB()
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()
	collectionName, err := stringfixer.DeleteExtension(filename)
	if err != nil {
		fatalErr = err
		return
	}
	col, err := db.C(collectionName)
	if err != nil {
		fatalErr = err
		return
	}

	// read lines in from the file and make one db entry per line
	file, err := os.Open(filename)
	if err != nil {
		fatalErr = err
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		col.Insert(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		fatalErr = err
		return
	}
}

// checkDB is a wrapper around filedb.Dial that will also create the database
// if it doesn't already exist
func checkDB() (*filedb.DB, error) {
	// check that the database exists and create if it doesn't
	if _, err := os.Stat(Dbpath); os.IsNotExist(err) {
		err := os.Mkdir(Dbpath, 0777)
		if err != nil {
			return nil, err
		}
	}

	db, err := filedb.Dial(Dbpath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Find searches the database and adds matches to the match
// channel as they are found so they can be reported to the
// user asap and not as a big dump of data
func (c *Client) Find(query string, r Retriever) {
	db, err := checkDB()
	if err != nil {
		log.Fatalln(err)
	}
	done := false
	defer db.Close()
	go func() {
		<-c.Done
		done = true
	}()
	collections, _ := db.ColNames()
	for _, colString := range collections {
		col, err := db.C(colString)
		if err != nil {
			log.Fatalln(err)
		}
		col.ForEach(func(_ int, data []byte) bool {
			if r.Match(query, string(data)) {
				if done {
					return true
				}
				c.Matches <- []byte(fmt.Sprintf("%v/%v", colString, string(data)))
			}
			return false
		})
		if done {
			break
		}

	}
	close(c.Matches)
}

// The same as Find, but with a progess channel that will output ints 0-100 until it is done
func (c *Client) FindProgress(query string, r Retriever) {
	db, err := checkDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	done := false
	collections, _ := db.ColNames()
	total := float64(len(collections))
	go func() {
		<-c.Done
		done = true
		close(c.Done)
	}()
	for i, colString := range collections {
		col, err := db.C(colString)
		if err != nil {
			log.Fatalln(err)
		}
		col.ForEach(func(_ int, data []byte) bool {
			if done {
				return true
			}
			if r.Match(query, string(data)) {
				c.Matches <- []byte(fmt.Sprintf("%v/%v", colString, string(data)))
			}
			return false
		})
		if done {
			break
		}
		c.Progress <- int(float64(i+1) / total * 100)

	}
	c.Progress <- 100
	close(c.Progress)
	close(c.Matches)
}
