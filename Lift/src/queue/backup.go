package queue

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"../defs"
	//"../hw" //NB! ikke gjør dette!
)

// runBackup loads queue data from file if file exists once and saves backups
// whenever its asked to.
func runBackup() {
	filenameLocal := "queueBackupFile1"
	filenameRemote := "queueBackupFile2"

	local.loadFromDisk(filenameLocal)
	// remote.loadFromDisk(filenameRemote)
	//remote = local //stygfiks er ikke beste fiks, funket ikke

	if !local.isEmpty() {
		log.Println(local)
		//go syncLightsAfterBackup()  
		defs.SyncLightsChan <- true
//problemet er at vi setter eksterne lys bare hvis de er i remote queue, hva er den beste måten å fikse dette på?
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
//problemet med dette er at syncLights kjøres etter dette, og siden det ikke ligger i remote skrives lysene over. 
/*func syncLightsAfterBackup() { // NB! duplikering av kode!
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if (b == defs.ButtonUp && f == defs.NumFloors-1) ||
				(b == defs.ButtonDown && f == 0) {
				continue
			} else {
				hw.SetButtonLamp(f, b, IsOrder(f, b))
			}
		}
	}
}*/


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
