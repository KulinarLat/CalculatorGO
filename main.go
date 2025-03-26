// Scientific Calculator in Go using Fyne GUI
// Author: KulinarLat
// License: MIT
// See README.md for details

package main

import (
	"fmt"
	"image/color"
	"math"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
)

const Version = "1.0.1"

var lastAnswer string = "0"
var isDarkTheme bool = false
var output *canvas.Text

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Инженерный калькулятор v" + Version)
	myWindow.Resize(fyne.NewSize(400, 600))

	isDegrees := true

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

	expression := ""
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
			lbl := label
			button := widget.NewButton(lbl, func() {
				switch lbl {
				case "C":
					expression = ""
				case "DEL":
					if len(expression) > 0 {
						expression = expression[:len(expression)-1]
					}
				case "=":
					res := eval(expression, isDegrees)
					col := color.RGBA{100, 100, 100, 255}
					if isDarkTheme {
						col = color.RGBA{255, 255, 255, 255}
					}
					entry := canvas.NewText(expression+" = "+res, col)
					entry.TextStyle = fyne.TextStyle{Italic: true}
					historyBox.Add(entry)
					historyContainer.ScrollToBottom()
					expression = res
				case "π":
					expression += fmt.Sprintf("%f", math.Pi)
				case "e":
					expression += fmt.Sprintf("%f", math.E)
				case "Ans":
					expression += lastAnswer
				case "sin", "cos", "tan", "log", "ln", "sqrt":
					expression += lbl + "("
				default:
					expression += lbl
				}
				output.Text = expression
				output.Refresh()
			})
			rowContainer.Add(button)
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

func replacePowers(expr string) string {
	re := regexp.MustCompile(`(\d+(\.\d+)?|\w+)\^(\d+(\.\d+)?|\w+)`)
	for re.MatchString(expr) {
		expr = re.ReplaceAllStringFunc(expr, func(s string) string {
			parts := strings.Split(s, "^")
			if len(parts) == 2 {
				return fmt.Sprintf("pow(%s,%s)", parts[0], parts[1])
			}
			return s
		})
	}
	return expr
}

func eval(expr string, isDegrees bool) string {
	expr = replacePowers(expr)

	functions := map[string]govaluate.ExpressionFunction{
		"sin": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x = x * math.Pi / 180
			}
			return math.Sin(x), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x = x * math.Pi / 180
			}
			return math.Cos(x), nil
		},
		"tan": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x = x * math.Pi / 180
			}
			return math.Tan(x), nil
		},
		"log": func(args ...interface{}) (interface{}, error) {
			return math.Log10(toFloat(args[0])), nil
		},
		"ln": func(args ...interface{}) (interface{}, error) {
			return math.Log(toFloat(args[0])), nil
		},
		"sqrt": func(args ...interface{}) (interface{}, error) {
			return math.Sqrt(toFloat(args[0])), nil
		},
		"pow": func(args ...interface{}) (interface{}, error) {
			return math.Pow(toFloat(args[0]), toFloat(args[1])), nil
		},
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		fmt.Println("Ошибка разбора выражения:", err)
		return "Ошибка ввода"
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		fmt.Println("Ошибка вычисления:", err)
		return "Ошибка вычисления"
	}

	lastAnswer = fmt.Sprintf("%v", result)
	return lastAnswer
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		fmt.Printf("toFloat: неподдерживаемый тип %T\n", val)
		return 0
	}
}
