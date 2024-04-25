package exportfile

import (
	"fmt"
	"os"
	"time"
	Authtimer "xApp/pkg/service/timer"
)

func CreateTimeToFile() error {
	// Create a new text file
	file, err := os.Create("time.txt")
	if err != nil {
		return fmt.Errorf("error creating to file: %v", err)
	}
	defer file.Close()
	return nil
}

func WriteTimeToFile() error {
	// Get the current time
	currentTime := time.Now()

	// Format the time value
	timeString := currentTime.Format("2006-01-02 15:04:05")

	// Create a new text file
	file, err := os.Create("time.txt")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the time value to the text file
	_, err = file.WriteString(timeString)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Println("Time value exported to time.txt")
	return nil
}

func ReadTimeFromFile(Type int, endTime time.Time) error {
	serviceTime := endTime.Sub(Authtimer.GetStartTime(1))
	if Type == 1 {
		fmt.Println("First Authentication transmission time: %v", serviceTime)
	} else {
		fmt.Println("NORA-AKA transmission time: %v", serviceTime)
	}
	//err := ioutil.WriteFile("time.txt", []byte(fmt.Sprintf("%d", serviceTime)), 0644)
	//if err != nil {
	//	return fmt.Errorf("error writing time to file: %v", err)
	//}

	return nil
}
