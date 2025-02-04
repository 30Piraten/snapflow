package handlers

import (
	"log"

	"github.com/30Piraten/snapflow/config"
)

// Handler to initiate print job
func InitiatePrintJob(customerEmail, photoID, processedS3Location string) error {
	err := config.SendPrintRequest(customerEmail, photoID, processedS3Location)
	if err != nil {
		log.Printf("failed to initiate print job: %v", err)
		return err
	}

	log.Println("✅ Print job initiated successfully.")
	return nil
}
