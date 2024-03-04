package main

import (
	"flag"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	generatorargs "github.com/seal-io/code-generator/cmd/webhook-gen/args"
	"github.com/seal-io/code-generator/cmd/webhook-gen/generators"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := generatorargs.NewDefaults()

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	_ = flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := generatorargs.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	klog.V(2).Info("Completed successfully.")
}
