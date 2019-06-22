package cmd

import (
	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

// DdevShareCommand contains the "ddev share" command
var DdevShareCommand = &cobra.Command{
	Use:   "share",
	Short: "Share project on the internet via ngrok.",
	Long:  `Use "ddev share" or add on extra ngrok commands, like "ddev share --subdomain some-subdomain". Although a few ngrok commands are supported directly, any ngrok flag can be added in the ngrok_args section of .ddev/config.yaml. You will want to create an account on ngrok.com and use the "ngrok authtoken" command to set up ngrok.`,
	Example: `ddev share
ddev share --subdomain some-subdomain
ddev share --use-http`,
	//Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := ddevapp.GetActiveApp("")
		if err != nil {
			util.Failed("Failed to get requested project: %v", err)
		}
		if app.SiteStatus() != ddevapp.SiteRunning {
			util.Failed("Project is not yet running. Use 'ddev start' first.")
		}
		// If they provided the --use-http flag, we'll not try both https and http
		useHTTP, err := cmd.Flags().GetBool("use-http")
		if err != nil {
			util.Failed("failed to get use-http flag: %v", err)
		}
		var flags []string
		// Pass along --subdomain argument
		if cmd.Flags().Changed("subdomain") {
			s, err := cmd.Flags().GetString("subdomain")
			if err != nil {
				util.Failed("Unable to get subdomain argument: %v", err)
			}
			flags = append(flags, "-subdomain", s)
		}

		err = app.Share(useHTTP, flags, 0)
		if err != nil {
			util.Failed("Unable to share: %v", err)
		}

		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(DdevShareCommand)
	DdevShareCommand.Flags().String("subdomain", "", `ngrok --subdomain argument, as in "ngrok --subdomain my-subdomain:, requires paid ngrok.com account"`)
	DdevShareCommand.Flags().Bool("use-http", false, `Set to true to use unencrypted http local tunnel (required if you do not have an ngrok.com account)"`)
}
