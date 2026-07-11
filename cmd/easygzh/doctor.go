package main

// Doctor — phase 4 will implement full diagnostics. Stub for now.
// (placeholder so the CLI compiles; the command reports version + Go runtime.)

import (
	"fmt"
	"os"
	"runtime"

	"github.com/easygzh/easygzh/internal/cli"
	"github.com/easygzh/easygzh/internal/memory"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run local environment diagnostics (no network, no auth).",
		RunE: func(cmd *cobra.Command, args []string) error {
			memoryStatus, memoryErr := memory.ValidateStore(memory.DefaultDir())
			if memoryErr != nil {
				memoryStatus.Issues = append(memoryStatus.Issues, memoryErr.Error())
			}
			_, appIDSet := os.LookupEnv("WECHAT_APPID")
			_, secretSet := os.LookupEnv("WECHAT_SECRET")
			data := map[string]any{
				"version":  Version,
				"go":       runtime.Version(),
				"platform": runtime.GOOS + "/" + runtime.GOARCH,
				"themes":   themesAvailable(cmd),
				"memory":   memoryStatus,
				"wechat": map[string]bool{
					"appid_configured":  appIDSet,
					"secret_configured": secretSet,
				},
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "easygzh %s on %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
			return emit(cmd, cli.OK("DOCTOR_OK", data))
		},
	}
	return cmd
}

func themesAvailable(_ *cobra.Command) []string {
	names, err := themeManager().List()
	if err != nil {
		return nil
	}
	return names
}
