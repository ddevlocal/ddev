package cmd

import (
	fapp "fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"time"

	"github.com/drud/ddev/pkg/ddevapp"
	"github.com/drud/ddev/pkg/output"
	"github.com/drud/ddev/pkg/util"
	"github.com/spf13/cobra"

	"fyne.io/fyne"
)

var sleepTime = 2

// GuiCommand provides the `ddev gui` command
var GUICommand = &cobra.Command{
	Use:     "gui",
	Short:   "Show experimental GUI",
	Example: `ddev gui`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			apps, err := ddevapp.GetProjects(false)
			if err != nil {
				util.Failed("failed getting GetProjects: %v", err)
			}
			appDescs := make([]map[string]interface{}, 0)

			if len(apps) < 1 {
				output.UserOut.WithField("raw", appDescs).Println("No ddev projects were found.")
			} else {
				fyneApp := fapp.New()
				window := fyneApp.NewWindow("DDEV-Local")

				headings := []string{"Name", "Type", "Location", "URL", "Status"}
				var rows [][]string

				for _, app := range apps {
					desc, err := app.Describe()
					if err != nil {
						util.Error("Failed to describe project %s: %v", app.GetName(), err)
					}
					status := desc["status"].(string)
					//switch {
					//case strings.Contains(status, ddevapp.SitePaused):
					//	status = color.YellowString(status)
					//case strings.Contains(status, ddevapp.SiteStopped):
					//	status = color.RedString(status)
					//case strings.Contains(status, ddevapp.SiteDirMissing):
					//	status = color.RedString(status)
					//case strings.Contains(status, ddevapp.SiteConfigMissing):
					//	status = color.RedString(status)
					//default:
					//	status = color.CyanString(status)
					//}

					rows = append(rows, []string{desc["name"].(string), desc["type"].(string), desc["approot"].(string), desc["httpsurl"].(string), status})
				}
				t := makeTable(headings, rows)

				window.SetContent(t)

				window.ShowAndRun()

				//output.UserOut.WithField("raw", appDescs).Print(table.String() + "\n" + ddevapp.RenderRouterStatus())
			}

			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	},
}

// From https://github.com/fyne-io/fyne/issues/157#issuecomment-597319590
func makeTable(headings []string, rows [][]string) *widget.Box {

	columns := rowsToColumns(headings, rows)

	objects := make([]fyne.CanvasObject, len(columns))
	for k, col := range columns {
		box := widget.NewVBox(widget.NewLabelWithStyle(headings[k], fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, val := range col {
			box.Append(widget.NewLabel(val))
		}
		objects[k] = box
	}
	return widget.NewHBox(objects...)
}

func rowsToColumns(headings []string, rows [][]string) [][]string {
	columns := make([][]string, len(headings))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}

//
//func loadUI(app fyne.App) {
//	//output := widget.NewLabel("")
//	//output.Alignment = fyne.TextAlignTrailing
//	//output.TextStyle.Monospace = true
//	//equals := addButton("=", func() {
//	//	evaluate()
//	//})
//	//equals.Style = widget.PrimaryButton
//	//
//	layout.NewAdaptiveGridLayout()
//	window := app.NewWindow("DDEV-Local")
//	window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
//		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
//			widget.NewLabel("Project"),
//			widget.NewLabel("Type"),
//			widget.NewLabel("Location"),
//			widget.NewLabel("URL"),
//		),
//
//		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
//			widget.NewLabel("d8composer"),
//			widget.NewLabel("drupal8"),
//			widget.NewLabel("~/workspace/d8composer"),
//			widget.NewLabel("https://d8composer.ddev.site"),
//		),
//	))
//
//	//window.Canvas().SetOnTypedRune(typedRune)
//	//window.Canvas().SetOnTypedKey(typedKey)
//	window.ShowAndRun()
//}
func init() {
	RootCmd.AddCommand(GUICommand)
}
