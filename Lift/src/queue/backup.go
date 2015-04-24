package queue

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// runBackup loads queue data from file if file exists once and saves backups
// whenever its asked to.
func runBackup() {
	filenameLocal := "queueBackupFile1"
	filenameRemote := "queueBackupFile2"

	local.loadFromDisk(filenameLocal)
	// remote.loadFromDisk(filenameRemote)

	if !local.isEmpty() {
		newOrder <- true
	}

	for {
		<-backup
		if err := local.saveToDisk(filenameLocal); err != nil {
			log.Println(err)
		}
		if err := remote.saveToDisk(filenameRemote); err != nil {
			log.Println(err)
		}
	}
}

func (q *queue) saveToDisk(filename string) error {

	data, err := json.Marshal(&q)
	if err != nil {
		log.Println("json.Marshal() error: Failed to backup.")
		return err
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println("ioutil.WriteFile() error: Failed to backup.")
		return err
	}

	return nil
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to the queue it's invoked on if the file is present.
func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Printf("Backup file %s exists, processing...\n", filename)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("loadFromDisk() error: Failed to read file.")
		}

		if err := json.Unmarshal(data, q); err != nil {
			log.Println("loadFromDisk() error: Failed to Unmarshal.")
		}
	}

	return nil
}
