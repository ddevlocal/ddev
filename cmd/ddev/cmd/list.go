package cmd

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"time"

	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/output"
	"github.com/drud/ddev/pkg/util"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// continuous, if set, makes list continuously output
var continuous bool

// activeOnly, if set, shows only running projects
var activeOnly bool

// continuousSleepTime is time to sleep between reads with --continuous
var continuousSleepTime = 1

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  `List projects. Shows all projects by default, shows active projects only with --active-only`,
	Example: `ddev list
ddev list -A`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			apps, err := ddevapp.GetProjects(activeOnly)
			if err != nil {
				util.Failed("failed getting GetProjects: %v", err)
			}
			appDescs := make([]map[string]interface{}, 0)

			if len(apps) < 1 {
				output.UserOut.WithField("raw", appDescs).Println("No ddev projects were found.")
				return
			}

			if cmd.Flag("gui").Changed {
				tguiPages(apps)
			} else {
				table := ddevapp.CreateAppTable()
				for _, app := range apps {
					desc, err := app.Describe()
					if err != nil {
						util.Error("Failed to describe project %s: %v", app.GetName(), err)
					}
					appDescs = append(appDescs, desc)
					ddevapp.RenderAppRow(table, desc)
				}
				output.UserOut.WithField("raw", appDescs).Print(table.String() + "\n" + ddevapp.RenderRouterStatus())
			}

			if !continuous {
				break
			}

			time.Sleep(time.Duration(continuousSleepTime) * time.Second)
		}
	},
}

func init() {
	ListCmd.Flags().BoolVarP(&activeOnly, "active-only", "A", false, "If set, only currently active projects will be displayed.")
	ListCmd.Flags().BoolVarP(&continuous, "continuous", "", false, "If set, project information will be emitted until the command is stopped.")
	ListCmd.Flags().IntVarP(&continuousSleepTime, "continuous-sleep-interval", "I", 1, "Time in seconds between ddev list --continous output lists.")
	ListCmd.Flags().Bool("gui", true, "Use gui mode.")

	RootCmd.AddCommand(ListCmd)
}

func tguiPages(apps []*ddevapp.DdevApp) {
	l, err := os.OpenFile("/tmp/listTable.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	tguiApp := tview.NewApplication()

	pages := tview.NewPages()
	table := listTable(apps)

	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			tguiApp.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		pages.SwitchToPage(apps[row-1].Name)
		tguiApp.SetFocus(pages)
	})

	pages.AddPage("listTable", table, true, true)

	for _, app := range apps {
		pages.AddPage(app.Name, descFlex(app), true, false)
	}

	if err := tguiApp.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func descFlex(app *ddevapp.DdevApp) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	d := tview.NewDropDown().SetOptions([]string{"stop", "start"}, nil).SetBorder(true).SetTitle(fmt.Sprintf("Title: %s", app.Name))
	flex.AddItem(d, 1, 1, true)

	desc, err := app.Describe()
	if err != nil {
		util.Failed("Failed to describe project %s: %v", app.Name, err)
	}
	renderedDesc, _ := renderAppDescribe(desc)
	flex.AddItem(tview.NewTextView().SetText(renderedDesc), 0, 1, false)

	return flex
}

func listTable(apps []*ddevapp.DdevApp) *tview.Table {
	table := tview.NewTable().SetBorders(false)
	table.SetSelectable(true, false)

	displayItems := []string{"name", "type", "approot", "primary_url", "status"}
	for tc, title := range displayItems {
		_ = table.SetCell(0, tc, tview.NewTableCell(title).SetAlign(tview.AlignLeft)).SetBorders(true)
	}
	for r, app := range apps {
		desc, err := app.Describe()
		desc["primary_url"] = app.GetPrimaryURL()

		if err != nil {
			util.Error("Failed to describe project %s: %v", app.GetName(), err)
		}
		for c, itemname := range displayItems {
			_ = table.SetCell(r+1, c, tview.NewTableCell(desc[itemname].(string)).SetAlign(tview.AlignLeft))
		}
	}

	return table
}
