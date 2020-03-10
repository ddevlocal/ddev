package cmd

import (
	fapp "fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"image/color"
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
				rows = append(rows, fyne.NewContainerWithLayout(layout.NewGridLayout(5), canvas.NewText("Project", color.White),
					canvas.NewText("Type", color.White),
					canvas.NewText("Location", color.White),
					canvas.NewText("URL", color.White),
					canvas.NewText("Status", color.White),
				))

				for _, app := range apps {
					desc, err := app.Describe()
					if err != nil {
						util.Error("Failed to describe project %s: %v", app.GetName(), err)
					}
					rows = append(rows, fyne.NewContainerWithLayout(layout.NewGridLayout(5),
						canvas.NewText(desc["name"].(string), color.White),
						canvas.NewText(desc["type"].(string), color.White),
						canvas.NewText(desc["approot"].(string), color.White),
						canvas.NewText(desc["httpsurl"].(string), color.White),
						canvas.NewText(desc["status"].(string), color.White),
					))
				}
				window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1), rows...))

				window.ShowAndRun()

				//output.UserOut.WithField("raw", appDescs).Print(table.String() + "\n" + ddevapp.RenderRouterStatus())
			}

			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	},
}

func loadUI(app fyne.App) {
	//output := widget.NewLabel("")
	//output.Alignment = fyne.TextAlignTrailing
	//output.TextStyle.Monospace = true
	//equals := addButton("=", func() {
	//	evaluate()
	//})
	//equals.Style = widget.PrimaryButton
	//
	window := app.NewWindow("DDEV-Local")
	//window.SetIcon(icon.CalculatorBitmap)
	window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			canvas.NewText("Project", color.White),
			canvas.NewText("Type", color.White),
			canvas.NewText("Location", color.White),
			canvas.NewText("URL", color.White),
		),

		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			canvas.NewText("d8composer", color.White),
			canvas.NewText("drupal8", color.White),
			canvas.NewText("~/workspace/d8composer", color.White),
			canvas.NewText("https://d8composer.ddev.site", color.White),
		),
	))

	//window.Canvas().SetOnTypedRune(typedRune)
	//window.Canvas().SetOnTypedKey(typedKey)
	window.ShowAndRun()
}
func init() {
	RootCmd.AddCommand(GUICommand)
}
