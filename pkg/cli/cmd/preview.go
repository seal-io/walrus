package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
	"github.com/seal-io/walrus/utils/gopool"
)

func Preview(sc *config.Config) (*cobra.Command, error) {
	flags := &manifest.PreviewOption{}

	use := "preview"
	cmd := &cobra.Command{
		Use:     use,
		GroupID: common.GroupAdvanced.ID,
		Example: previewExample("preview"),
		Short:   "Generate or apply resource previews using the provided file path or folder.",
		PreRun:  setupServerContextForPreviewFunc(sc, flags),
		Run: func(cmd *cobra.Command, args []string) {
			// Load from files.
			loader := manifest.DefaultLoader(sc, flags.ValidateParametersSet, flags.RunLabels)
			set, err := loader.LoadFromFiles(flags.Filenames, flags.Recursive)
			if err != nil {
				panic(err)
			}

			// Apply.
			if flags.Apply {
				err = previewApply(sc, set, flags)
				if err != nil {
					panic(err)
				}
				return
			}

			// Preview.
			err = preview(sc, set, flags)
			if err != nil {
				panic(err)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd, nil
}

func setupServerContextForPreviewFunc(sc *config.Config, opts *manifest.PreviewOption) func(*cobra.Command, []string) {
	return setupServerContextFunc(sc, &opts.CommonOption)
}

func preview(sc *config.Config, set manifest.ObjectSet, flags *manifest.PreviewOption) error {
	ctx := context.Background()

	// Apply the files.
	operator := manifest.DefaultPreviewOperator(sc, flags)
	r, err := operator.Operate(set)
	operator.PrintResult(r)
	if err != nil {
		return err
	}

	// Wait for the result.
	if flags.Wait {
		waiter := manifest.DefaultPreviewWaiter(sc, flags.Timeout)
		_, _ = waiter.Wait(ctx, set, nil)
	}

	return nil
}

func previewApply(sc *config.Config, set manifest.ObjectSet, flags *manifest.PreviewOption) error {
	resultChan := make(chan manifest.ObjectSet, 4)
	ctx := context.Background()
	wg := gopool.GroupWithContextIn(ctx)

	// Wait for the result.
	if flags.Wait {
		wg.Go(func(ctx context.Context) error {
			waiter := manifest.DefaultPreviewApplyWaiter(sc, flags.Timeout)
			_, err := waiter.Wait(ctx, set, resultChan)
			return err
		})
	}

	// Apply the files.
	wg.Go(func(ctx context.Context) error {
		operator := manifest.DefaultPreviewApplyOperator(sc, flags)
		r, err := operator.Operate(set)
		operator.PrintResult(r)
		if err != nil {
			return err
		}

		resultChan <- r.Success
		return nil
	})

	return wg.Wait()
}

func previewExample(action string) string {
	return fmt.Sprintf(`  # Generate preview of the configuration in the walrus-file.yaml 
  $ walrus %s -f walrus-file.yaml --run-labels key=value

  # Generate preview of yaml files from a directory
  $ walrus %s -f dir/ --run-labels key=value

  # Generate preview of yamls files from a directory recursively
  $ walrus %s -f dir/ --recursive --run-labels key=value

  # Generate preview of yaml files with a specific project/environment
  $ walrus %s -f dir/ --run-labels key=value -p my-project -e my-environment 

  # Apply generated preview with specific labels
  $ walrus %s -f walrus-file.yaml --run-labels key=value --apply
`,
		action,
		action,
		action,
		action,
		action)
}
