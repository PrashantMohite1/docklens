package logger

import (
	"fmt"
	"log/slog"
)

func Generate_report(checked int, Matched int, MisMatched int) {
	slog.Info("\nVerification Report")
	slog.Info("---------------------")
	slog.Info(fmt.Sprintf("Checked: %d", checked))
	slog.Info(fmt.Sprintf("Matched: %d", Matched))
	slog.Info(fmt.Sprintf("Mismatched: %d", MisMatched))
}
