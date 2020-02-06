package cmd

import (
	"github.com/gdamore/tcell"
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
				tgui(apps)
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

func tgui(apps []*ddevapp.DdevApp) {
	tguiApp := tview.NewApplication()
	table := list(apps)
	table.SetSelectable(true, false)

	table.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			tguiApp.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		//_ = table.GetCell(row, column).SetTextColor(tcell.ColorRed)

	})
	if err := tguiApp.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}

}

func tguiFlex(apps []*ddevapp.DdevApp) {
	//	tguiApp := tview.NewApplication()
	//
	//	flex := tview.NewFlex().
	//		AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false)
	//	//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
	//	//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
	//	//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
	//	//AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	//	i := flex.AddItem(tview.NewList(), 0, 0, true)
	//	if err := tguiApp.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
	//		panic(err)
	//	}
	//
	//	tview.NewFlex().AddItem()
	//	displayItems := []string{"name", "type", "approot", "primary_url", "status"}
	//	for tc, title := range displayItems {
	//		_ = table.SetCell(0, tc, tview.NewTableCell(title).SetAlign(tview.AlignLeft)).SetBorders(true)
	//	}
	//	for r, app := range apps {
	//		desc, err := app.Describe()
	//		desc["primary_url"] = app.GetPrimaryURL()
	//
	//		if err != nil {
	//			util.Error("Failed to describe project %s: %v", app.GetName(), err)
	//		}
	//		for c, itemname := range displayItems {
	//			_ = table.SetCell(r+1, c, tview.NewTableCell(desc[itemname].(string)).SetAlign(tview.AlignLeft))
	//
	//		}
	//	}
	//	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
	//		if key == tcell.KeyEscape {
	//			tguiApp.Stop()
	//		}
	//		if key == tcell.KeyEnter {
	//			table.SetSelectable(true, true)
	//		}
	//	}).SetSelectedFunc(func(row int, column int) {
	//		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	//		table.SetSelectable(false, false)
	//		d := tview.NewDropDown().SetLabel(table.GetCell(row, column).Text).SetOptions([]string{"stop", "start", "describe", "launch"}, nil)
	//		//modal := tview.NewModal().
	//		//	SetText("Do you want to quit the application?").
	//		//	AddButtons([]string{"Quit", "Cancel"}).
	//		//	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	//		//		//if buttonLabel == "Quit" {
	//		//		//	app.Stop()
	//		//		//}
	//		//	})
	//
	//		c := table.GetCell(row, column)
	//		tview.NewTableCell("xxx")
	//		tview.new
	//		if err := tguiApp.SetRoot(d, false).SetFocus(d).Run(); err != nil {
	//			panic(err)
	//		}
	//
	//	})
	//	if err := tguiApp.SetRoot(table, true).Run(); err != nil {
	//		panic(err)
	//	}

	//for r := 0; r < len(apps); r++ {
	//	for c := 0; c < cols; c++ {
	//		color := tcell.ColorWhite
	//		if c < 1 || r < 1 {
	//			color = tcell.ColorYellow
	//		}
	//		table.SetCell(r, c,
	//			tview.NewTableCell(lorem[word]).
	//				SetTextColor(color).
	//				SetAlign(tview.AlignCenter))
	//		word = (word + 1) % len(lorem)
	//	}
	//}
	//table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
	//	if key == tcell.KeyEscape {
	//		tguiApp.Stop()
	//	}
	//	if key == tcell.KeyEnter {
	//		table.SetSelectable(true, true)
	//	}
	//}).SetSelectedFunc(func(row int, column int) {
	//	table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	//	table.SetSelectable(false, false)
	//})
	//if err := tguiApp.SetRoot(table, true).Run(); err != nil {
	//	panic(err)
	//}
}

func tguiPages(apps []*ddevapp.DdevApp) {
	tguiApp := tview.NewApplication()
	pages := tview.NewPages()

	pages.AddPage("list",
		tview.NewModal().
			SetText("This is page. Choose where to go next.").
			AddButtons([]string{"Next", "Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage("list")
				} else {
					tguiApp.Stop()
				}
			}),
		false,
		true)
	if err := tguiApp.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}

}

func list(apps []*ddevapp.DdevApp) *tview.Table {
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
