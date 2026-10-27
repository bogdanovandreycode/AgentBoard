//go:build !windows

package httpapi

import (
	"context"
	"fmt"
)

func selectWindowsFolder(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("Windows folder picker is unavailable")
}
