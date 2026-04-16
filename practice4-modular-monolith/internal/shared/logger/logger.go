package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
)

var std = log.New(os.Stdout, "", log.LstdFlags)

func formatMessage(ctx context.Context, format string, v ...any) string {
	corrID := ctxutil.GetCorrelationID(ctx)
	msg := fmt.Sprintf(format, v...)
	return fmt.Sprintf("[CorrID: %s] %s", corrID, msg)
}

func Printf(ctx context.Context, format string, v ...any) {
	std.Print(formatMessage(ctx, format, v...))
}

func Println(ctx context.Context, v ...any) {
	corrID := ctxutil.GetCorrelationID(ctx)
	msg := fmt.Sprint(v...)
	std.Printf("[CorrID: %s] %s", corrID, msg)
}

func Fatalf(ctx context.Context, format string, v ...any) {
	std.Fatal(formatMessage(ctx, format, v...))
}

func Fatal(ctx context.Context, v ...any) {
	corrID := ctxutil.GetCorrelationID(ctx)
	msg := fmt.Sprint(v...)
	std.Fatalf("[CorrID: %s] %s", corrID, msg)
}
