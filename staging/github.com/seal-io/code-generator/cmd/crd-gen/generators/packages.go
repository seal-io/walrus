package generators

import (
	"fmt"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/klog/v2"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(0),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("//go:build !%s\n// +build !%s\n\n", arguments.GeneratedBuildTag, arguments.GeneratedBuildTag)), boilerplate...)

	for i := range context.Inputs {
		klog.V(5).Infof("considering pkg %q", context.Inputs[i])
		pkg := context.Universe[context.Inputs[i]]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		path := pkg.Path
		// If the source path is within a /vendor/ directory (for example,
		// k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/apis/meta/v1), allow
		// generation to output to the proper relative path (under vendor).
		// Otherwise, the generator will create the file in the wrong location
		// in the output directory.
		if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
			expandedPath := strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
			if strings.HasPrefix(expandedPath, "/vendor/") {
				path = expandedPath
			}
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: pkg.Name,
				PackagePath: path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) []generator.Generator {
					return []generator.Generator{
						NewCRDGen(arguments.OutputFileBaseName, path),
					}
				},
			})
	}

	return packages
}
