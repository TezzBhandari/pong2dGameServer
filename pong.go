package pong

import "fmt"

func ReportPanic(v any) {
	fmt.Printf("[ReportPanic] %q", v)
}
