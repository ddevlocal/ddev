package cmd

import (
	"fmt"
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

				s := fyne.NewSize(1, 0)
				nameLabel := widget.NewLabel("Project")
				nameLabel.Resize(s)

				nameCol := []fyne.CanvasObject{nameLabel}
				typeCol := []fyne.CanvasObject{widget.NewLabel("Type")}
				locCol := []fyne.CanvasObject{widget.NewLabel("Location")}
				urlCol := []fyne.CanvasObject{widget.NewLabel("URL")}
				statusCol := []fyne.CanvasObject{widget.NewLabel("Status")}
				//var nameCol, typeCol, locCol, urlCol, statusCol []fyne.CanvasObject

				for _, app := range apps {
					desc, err := app.Describe()
					if err != nil {
						util.Error("Failed to describe project %s: %v", app.GetName(), err)
					}

					nameLabel := widget.NewLabel(desc["name"].(string))
					nameLabel.Resize(s)
					ms := nameLabel.MinSize()
					fmt.Printf("nameLabel ms=%v", ms)
					nameCol = append(nameCol, nameLabel)
					//nameCol = append(nameCol, widget.NewLabel(desc["name"].(string)))
					typeCol = append(typeCol, widget.NewLabel(desc["type"].(string)))
					locCol = append(locCol, widget.NewLabel(desc["approot"].(string)))
					x := widget.NewHyperlink(desc["httpsurl"].(string), nil)
					_ = x.SetURLFromString(desc["httpsurl"].(string))

					urlCol = append(urlCol, x)
					statusCol = append(statusCol, widget.NewLabel(desc["status"].(string)))
				}

				nameCont := fyne.NewContainerWithLayout(layout.NewGridLayout(1), nameCol...)
				fmt.Printf("nameCont ms=%v nameCont size=%v", nameCont.MinSize(), nameCont.Size())

				typeCont := fyne.NewContainerWithLayout(layout.NewGridLayout(1), typeCol...)
				locCont := fyne.NewContainerWithLayout(layout.NewGridLayout(1), locCol...)
				urlCont := fyne.NewContainerWithLayout(layout.NewGridLayout(1), urlCol...)
				statusCont := fyne.NewContainerWithLayout(layout.NewGridLayout(1), statusCol...)

				window.SetContent(
					fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(1),
						nameCont, typeCont, locCont, urlCont, statusCont,
						//fyne.NewContainerWithLayout(layout.NewGridLayout(1), nameCol...),
						//fyne.NewContainerWithLayout(layout.NewGridLayout(1), typeCol...),
						//fyne.NewContainerWithLayout(layout.NewGridLayout(1), locCol...),
						//fyne.NewContainerWithLayout(layout.NewGridLayout(1), urlCol...),
						//fyne.NewContainerWithLayout(layout.NewGridLayout(1), statusCol...)),
					))

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
