package queue

import (
	"encoding/gob"
	"log"
	"os"
)

const diskDebug = false

// runBackup loads queue data from file if file exists once and saves backups
// whenever its asked to.
func runBackup() {
	filenameLocal := "localQueueBackup"
	filenameRemote := "remoteQueueBackup"

	local.loadFromDisk(filenameLocal)
	// remote.loadFromDisk(filenameRemote)

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
	fi, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	if err := gob.NewEncoder(fi).Encode(q); err != nil {
		return err
	}

	if diskDebug {
		log.Printf("Successful save of file %s\n", filename)
	}
	return nil
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to the queue it's invoked on if the file is present.
func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Printf("Backup file %s exists, processing...\n", filename)
		fi, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fi.Close()

		if err := gob.NewDecoder(fi).Decode(&q); err != nil {
			return err
		}
	}
	return nil
}
