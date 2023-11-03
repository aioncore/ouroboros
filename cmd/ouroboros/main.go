package main

import (
	"github.com/aioncore/ouroboros/pkg/cmd"
	"github.com/aioncore/ouroboros/pkg/config"
	"github.com/aioncore/ouroboros/pkg/log"
	"github.com/aioncore/ouroboros/pkg/service/utils"
	"github.com/spf13/cobra"
)

func main() {
	rootPath := utils.GetRootPath()
	log.InitLogger(rootPath, "core")
	config.InitConfig(rootPath, "core")
	rootCmd := &cobra.Command{
		Use:   "ouroboros",
		Short: "multiply blockchain core",
	}
	rootCmd.AddCommand(
		cmd.NewStartCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
