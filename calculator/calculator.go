package calculator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

var powerRegex = regexp.MustCompile(`(\d+(\.\d+)?|\w+)\^(\d+(\.\d+)?|\w+)`)
var scientificNotationRegex = regexp.MustCompile(`([\d\.]+)e([+-]?\d+)`)

func normalizeScientificNotation(expr string) string {
	return scientificNotationRegex.ReplaceAllStringFunc(expr, func(s string) string {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return s // вернуть как есть, пусть govaluate отловит
		}
		return fmt.Sprintf("%f", val)
	})
}

func replacePowers(expr string) string {
	for powerRegex.MatchString(expr) {
		expr = powerRegex.ReplaceAllStringFunc(expr, func(s string) string {
			parts := strings.Split(s, "^")
			if len(parts) == 2 {
				return fmt.Sprintf("pow(%s,%s)", parts[0], parts[1])
			}
			return s
		})
	}
	return expr
}

func formatResult(val float64) string {
	if math.Abs(val) < 1e-10 {
		return "0"
	}
	return fmt.Sprintf("%.10g", val)
}

func Eval(expr string, isDegrees bool) string {
	expr = normalizeScientificNotation(expr)
	expr = replacePowers(expr)

	functions := map[string]govaluate.ExpressionFunction{
		"sin": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x *= math.Pi / 180
			}
			return math.Sin(x), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x *= math.Pi / 180
			}
			return math.Cos(x), nil
		},
		"tan": func(args ...interface{}) (interface{}, error) {
			x := toFloat(args[0])
			if isDegrees {
				x *= math.Pi / 180
			}
			return math.Tan(x), nil
		},
		"log": func(args ...interface{}) (interface{}, error) {
			val := toFloat(args[0])
			if val <= 0 {
				return nil, fmt.Errorf("логарифм недопустим")
			}
			return math.Log10(val), nil
		},
		"ln": func(args ...interface{}) (interface{}, error) {
			val := toFloat(args[0])
			if val <= 0 {
				return nil, fmt.Errorf("натуральный логарифм недопустим")
			}
			return math.Log(val), nil
		},
		"sqrt": func(args ...interface{}) (interface{}, error) {
			val := toFloat(args[0])
			if val < 0 {
				return nil, fmt.Errorf("корень из отрицательного")
			}
			return math.Sqrt(val), nil
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

	return formatResult(result.(float64))
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
