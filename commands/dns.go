package commands

import (
	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/dns"
)

var (

	// DnsCmd ...
	DnsCmd = &cobra.Command{
		Use:   "dns",
		Short: "Manage dns aliases for local applications.",
		Long: `
Manages dns aliases for local applications. This modifies
your local hosts file, requiring administrative privileges.
		`,
	}
)

//
func init() {
	DnsCmd.AddCommand(dns.AddCmd)
	DnsCmd.AddCommand(dns.RemoveCmd)
	DnsCmd.AddCommand(dns.RemoveAllCmd)
	DnsCmd.AddCommand(dns.ListCmd)
}
