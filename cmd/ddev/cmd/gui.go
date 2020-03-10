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
	"fyne.io/fyne/layout"
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

				var rows []fyne.CanvasObject
				//rows := []*fyne.Container{fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				rows = append(rows, widget.NewLabel("Project"),
					widget.NewLabel("Type"),
					widget.NewLabel("Location"),
					widget.NewLabel("URL"),
					widget.NewLabel("Status"),
				)

				for _, app := range apps {
					desc, err := app.Describe()
					if err != nil {
						util.Error("Failed to describe project %s: %v", app.GetName(), err)
					}
					rows = append(rows,
						widget.NewLabel(desc["name"].(string)),
						widget.NewLabel(desc["type"].(string)),
						widget.NewLabel(desc["approot"].(string)),
						widget.NewLabel(desc["httpsurl"].(string)),
						widget.NewLabel(desc["status"].(string)),
					)
				}
				window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(5), rows...))

				window.ShowAndRun()

				//output.UserOut.WithField("raw", appDescs).Print(table.String() + "\n" + ddevapp.RenderRouterStatus())
			}

			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	},
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
