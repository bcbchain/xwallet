package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	HomeFlag	= "home"
	TraceFlag	= "trace"
	OutputFlag	= "output"
	EncodingFlag	= "encoding"
)

type Executable interface {
	Execute() error
}

func PrepareBaseCmd(cmd *cobra.Command, envPrefix, defaultHome string) Executor {
	cobra.OnInitialize(func() { initEnv(envPrefix) })
	cmd.PersistentFlags().StringP(HomeFlag, "", defaultHome, "directory for config and data")
	cmd.PersistentFlags().Bool(TraceFlag, false, "print out full stack trace on errors")
	cmd.PersistentPreRunE = concatCobraCmdFuncs(bindFlagsLoadViper, cmd.PersistentPreRunE)
	return Executor{cmd, os.Exit}
}

func PrepareMainCmd(cmd *cobra.Command, envPrefix, defaultHome string) Executor {
	cmd.PersistentFlags().StringP(EncodingFlag, "e", "hex", "Binary encoding (hex|b64|btc)")
	cmd.PersistentFlags().StringP(OutputFlag, "o", "text", "Output format (text|json)")
	cmd.PersistentPreRunE = concatCobraCmdFuncs(validateOutput, cmd.PersistentPreRunE)
	return PrepareBaseCmd(cmd, envPrefix, defaultHome)
}

func initEnv(prefix string) {
	copyEnvVars(prefix)

	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}

func copyEnvVars(prefix string) {
	prefix = strings.ToUpper(prefix)
	ps := prefix + "_"
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			k, v := kv[0], kv[1]
			if strings.HasPrefix(k, prefix) && !strings.HasPrefix(k, ps) {
				k2 := strings.Replace(k, prefix, ps, 1)
				os.Setenv(k2, v)
			}
		}
	}
}

type Executor struct {
	*cobra.Command
	Exit	func(int)
}

type ExitCoder interface {
	ExitCode() int
}

func (e Executor) Execute() error {
	e.SilenceUsage = true
	e.SilenceErrors = true
	err := e.Command.Execute()
	if err != nil {
		if viper.GetBool(TraceFlag) {
			fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		}

		exitCode := 1
		if ec, ok := err.(ExitCoder); ok {
			exitCode = ec.ExitCode()
		}
		e.Exit(exitCode)
	}
	return err
}

type cobraCmdFunc func(cmd *cobra.Command, args []string) error

func concatCobraCmdFuncs(fs ...cobraCmdFunc) cobraCmdFunc {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range fs {
			if f != nil {
				if err := f(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func bindFlagsLoadViper(cmd *cobra.Command, args []string) error {

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	homeDir := viper.GetString(HomeFlag)
	viper.Set(HomeFlag, homeDir)
	viper.SetConfigName("config")
	viper.AddConfigPath(homeDir)
	viper.AddConfigPath(filepath.Join(homeDir, "config"))

	if err := viper.ReadInConfig(); err == nil {

	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {

		return err
	}
	return nil
}

func validateOutput(cmd *cobra.Command, args []string) error {

	output := viper.GetString(OutputFlag)
	switch output {
	case "text", "json":
	default:
		return errors.Errorf("Unsupported output format: %s", output)
	}
	return nil
}
