package queue

import (
	def "../config"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// suggestBackup tells the backup routine to make a new backup if the channel
// backupChan does not currently block.
func suggestBackup() {
	select {
	case backupChan <- true:
	default:
	}
}

// runBackup loads queue data from file if file exists once, and saves
// backups whenever its asked to.
func runBackup() {
	filenameLocal := "queueBackupFile1"
	filenameRemote := "queueBackupFile2"

	var backup queue
	backup.loadFromDisk(filenameLocal)

	// Resend all orders found on loaded backup file:
	if !backup.isEmpty() {
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if backup.isActiveOrder(f, b) {
					if b == def.ButtonCommand {
						AddLocalOrder(f, b)
					} else {
						def.MessageChan <- def.Message{
							Kind:   def.NewOrder,
							Floor:  f,
							Button: b}
					}
				}
			}
		}
	}

	go func() {
		for {
			<-backupChan
			if err := local.saveToDisk(filenameLocal); err != nil {
				log.Println(err)
			}
			if err := remote.saveToDisk(filenameRemote); err != nil {
				log.Println(err)
			}
		}
	}()
}

// saveToDisk saves a queue to disk.
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
// saves its contents to a queue if the file is present.
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
