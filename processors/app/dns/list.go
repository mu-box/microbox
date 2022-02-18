package dns

import (
	"fmt"

	"github.com/mu-box/microbox/models"
	"github.com/mu-box/microbox/util/dns"
)

// List lists all dns entries for an app
func List(a *models.App) error {
	if a.ID == "" {
		fmt.Println("No DNS aliases registered")
		return nil
	}

	// print the header
	fmt.Printf("\nDNS Aliases\n")

	// iterate
	for _, domain := range dns.List(a.ID) {
		fmt.Printf("  %s\n", domain.Domain)
	}

	fmt.Println()

	return nil
}
