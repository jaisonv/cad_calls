package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jaisonv/telegram-cad-bot/internal/cad"
)

func main() {
	baseURL := flag.String("base-url", "https://southmiamipdfl.policetocitizen.com", "CAD API base URL")
	agencyID := flag.Int("agency-id", 386, "Agency ID")
	verifySSL := flag.Bool("verify-ssl", false, "Verify SSL certificates")
	flag.Parse()

	fmt.Printf("Testing CAD API Connection\n")
	fmt.Printf("===========================\n")
	fmt.Printf("Base URL:  %s\n", *baseURL)
	fmt.Printf("Agency ID: %d\n", *agencyID)
	fmt.Printf("Verify SSL: %v\n\n", *verifySSL)

	// Create client
	client := cad.NewClient(*baseURL, *agencyID, *verifySSL, 30*time.Second)

	fmt.Println("Fetching active CAD calls...")
	resp, err := client.GetActiveCalls(10)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}

	fmt.Printf("\n✅ Success! Retrieved %d calls (Total: %d)\n\n", len(resp.CADCalls), resp.Total)

	if len(resp.CADCalls) > 0 {
		fmt.Println("Sample call:")
		call := resp.CADCalls[0]
		fmt.Printf("  Incident ID: %s\n", call.IncidentID)
		fmt.Printf("  Type: %s\n", call.CallType)
		fmt.Printf("  Nature: %s\n", call.Nature)
		fmt.Printf("  Address: %s\n", call.Address)
		fmt.Printf("  Time: %s\n", call.StartTime.Format("2006-01-02 15:04"))
	} else {
		fmt.Println("No active calls at this time.")
	}

	fmt.Println("\n✅ Configuration is working! You can use these settings with the bot.")
}
