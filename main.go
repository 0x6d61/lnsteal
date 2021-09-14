package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var serverIp string

func lnstealMain(c *cobra.Command, args []string) {
	build, err := c.PersistentFlags().GetBool("build")
	if err != nil {
		fmt.Println(err)
	}

	server, err := c.PersistentFlags().GetBool("server")
	if err != nil {
		fmt.Println(err)
	}

	if build && serverIp != "" {
		fmt.Println("powershell")
		os.Exit(0)
	} else if server {

	}

}

func main() {
	rootCmd := &cobra.Command{
		Use: "lnsteal",
		Run: lnstealMain,
	}

	rootCmd.PersistentFlags().BoolP("build", "b", true, "build the client powershell code")
	rootCmd.PersistentFlags().StringVarP(&serverIp, "ip", "i", "", "IP address for client build")
	rootCmd.PersistentFlags().BoolP("server", "s", true, "starting the server")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
