package cmd

import (
	"dmtool/action"
	"dmtool/unix"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	toolCmd = &cobra.Command{
		Use: "tool",
		Run: toolFunc,
	}

	defaultSockerFile = "./dm.Sock"
)

func init() {
	toolCmd.Flags().StringP("action", "a", "", "tool action: ["+RESET_ACCOUNT+"]")
	if err := toolCmd.MarkFlagRequired("action"); err != nil {
		panic(err)
	}

	toolCmd.Flags().StringP("socket", "s", "./dm.sock", "target unix socket file")
}

func toolFunc(cmd *cobra.Command, args []string) {
	actionName, err := cmd.Flags().GetString("action")
	if err != nil {
		fmt.Printf("failed to get action: %v\n", err)
		return
	}

	defaultSockerFile, err := cmd.Flags().GetString("socket")
	if err != nil {
		fmt.Printf("failed to get socket file: %v\n", err)
		return
	}

	unixClient := unix.NewUnixClient(defaultSockerFile)
	if unixClient == nil {
		return
	}
	defer unixClient.Close()

	switch actionName {
	case RESET_ACCOUNT:
		action.HandleResetAccount(unixClient)
	default:
		panic("unknown action, please use -h to get valid actions")
	}
}

func Execute() {
	if err := toolCmd.Execute(); err != nil {
		panic(err)
	}
}
