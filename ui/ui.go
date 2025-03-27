package ui

import (
	"Test/calculator"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"strings"
)

var (
	lastAnswer  string = "0"
	isDarkTheme bool   = false
	output      *canvas.Text
)

func StartApp() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Инженерный калькулятор v1.0.2")
	myWindow.Resize(fyne.NewSize(400, 600))

	isDegrees := true
	expression := ""
	wasError := false

	modeLabel := widget.NewLabel("Режим: Degrees")
	var toggleMode *widget.Button
	toggleMode = widget.NewButton("Переключить на Radians", func() {
		isDegrees = !isDegrees
		if isDegrees {
			modeLabel.SetText("Режим: Degrees")
			toggleMode.SetText("Переключить на Radians")
		} else {
			modeLabel.SetText("Режим: Radians")
			toggleMode.SetText("Переключить на Degrees")
		}
	})

	historyBox := container.NewVBox()
	historyContainer := container.NewVScroll(historyBox)
	historyContainer.SetMinSize(fyne.NewSize(380, 120))

	var themeButton *widget.Button
	themeButton = widget.NewButton("Тема: Светлая", func() {
		isDarkTheme = !isDarkTheme
		if isDarkTheme {
			myApp.Settings().SetTheme(theme.DarkTheme())
			themeButton.SetText("Тема: Тёмная")
			output.Color = color.RGBA{255, 255, 255, 255}
		} else {
			myApp.Settings().SetTheme(theme.LightTheme())
			themeButton.SetText("Тема: Светлая")
			output.Color = color.RGBA{0, 0, 0, 255}
		}
		output.Refresh()
		for _, obj := range historyBox.Objects {
			if text, ok := obj.(*canvas.Text); ok {
				if isDarkTheme {
					text.Color = color.RGBA{255, 255, 255, 255}
				} else {
					text.Color = color.RGBA{100, 100, 100, 255}
				}
				text.Refresh()
			}
		}
	})

	output = canvas.NewText("", color.RGBA{0, 0, 0, 255})
	output.TextStyle = fyne.TextStyle{Bold: true}
	output.Alignment = fyne.TextAlignTrailing
	output.TextSize = 20

	buttons := [][]string{
		{"7", "8", "9", "/", "sqrt"},
		{"4", "5", "6", "*", "^"},
		{"1", "2", "3", "-", "log"},
		{"0", ".", "(", ")", "ln"},
		{"π", "e", "sin", "cos", "tan"},
		{"Ans", "C", "=", "+", "DEL"},
	}

	grid := container.NewGridWithRows(len(buttons))
	for _, row := range buttons {
		rowContainer := container.NewGridWithColumns(len(row))
		for _, label := range row {
			rowContainer.Add(widget.NewButton(label, func(lbl string) func() {
				return func() {
					if wasError && lbl != "=" && lbl != "C" && lbl != "DEL" {
						expression = ""
						wasError = false
					}
					switch lbl {
					case "C":
						expression = ""
					case "DEL":
						if len(expression) > 0 {
							expression = expression[:len(expression)-1]
						}
					case "=":
						res := calculator.Eval(expression, isDegrees)
						col := color.RGBA{100, 100, 100, 255}
						if isDarkTheme {
							col = color.RGBA{255, 255, 255, 255}
						}
						historyEntry := canvas.NewText(expression+" = "+res, col)
						historyEntry.TextStyle = fyne.TextStyle{Italic: true}
						historyBox.Add(historyEntry)
						historyContainer.ScrollToBottom()
						if strings.HasPrefix(res, "Ошибка") {
							wasError = true
						} else {
							expression = res
							wasError = false
						}
					case "π":
						expression += fmt.Sprintf("%f", math.Pi)
					case "e":
						expression += fmt.Sprintf("%f", math.E)
					case "Ans":
						expression += fmt.Sprintf("(%v)", lastAnswer)
					case "sin", "cos", "tan", "log", "ln", "sqrt":
						expression += lbl + "("
					default:
						expression += lbl
					}
					output.Text = expression
					output.Refresh()
				}
			}(label)))
		}
		grid.Add(rowContainer)
	}

	myWindow.SetContent(container.NewVBox(
		themeButton,
		modeLabel,
		toggleMode,
		historyContainer,
		layout.NewSpacer(),
		output,
		grid,
	))

	myWindow.ShowAndRun()
}
