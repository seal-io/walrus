package cmd

import (
	"flag"
	"os"
	"strings"

	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

func Init(c *cobra.Command) *cobra.Command {
	gcl, pcl := flag.CommandLine, pflag.CommandLine

	// Support klog configuration.
	klog.InitFlags(gcl)
	{ // Default klog configuration.
		_ = gcl.Set("logtostderr", "true")
		_ = gcl.Set("v", "4")
	}
	pcl.AddGoFlag(gcl.Lookup("v"))           // --v
	pcl.AddGoFlag(gcl.Lookup("vmodule"))     // --vmodule
	pcl.AddGoFlag(gcl.Lookup("logtostderr")) // --logtostderr
	{                                        // Hide klog flags.
		_ = pcl.MarkHidden("logtostderr")
	}

	// Support printing command line.
	printCmdline := pcl.Bool("print-cmdline", false,
		"print cmdline, which includes the arguments retrieved from environment.")
	c.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if *printCmdline {
			c.Printf("%s\n\n", strings.Join(os.Args, " "))
		}
		return nil
	}

	// Silence usage/errors,
	// and return help message if flag error occurs.
	c.SilenceUsage = true
	c.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		_ = c.Help()
		return err
	})

	// Append version.
	c.Version = version.Get()

	// Retrieve args from environment variables.
	osx.RetrieveArgsFromEnvInto(c)

	return c
}
