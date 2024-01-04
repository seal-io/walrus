package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/strs"
)

func Apply(sc *config.Config) (*cobra.Command, error) {
	flags := &manifest.Option{}

	use := "apply"
	cmd := &cobra.Command{
		Use:     use,
		GroupID: common.GroupAdvanced.ID,
		Example: manifestExample("apply"),
		Short:   "Apply a configuration to a resource using the provided file path or folder.",
		PreRun:  setupServerContextFunc(sc, flags),
		Run: func(cmd *cobra.Command, args []string) {
			// Load from files.
			loader := manifest.DefaultLoader(sc, flags.ValidateParametersSet)
			set, err := loader.LoadFromFiles(flags.Filenames, flags.Recursive)
			if err != nil {
				panic(err)
			}

			resultChan := make(chan manifest.ObjectSet, 4)
			ctx := context.Background()
			wg := gopool.GroupWithContextIn(ctx)

			// Wait for the result.
			if flags.Wait {
				wg.Go(func(ctx context.Context) error {
					waiter := manifest.StatusWaiter(sc, flags.Timeout)
					_, err = waiter.Wait(ctx, set, resultChan)
					return err
				})
			}

			// Apply the files.
			wg.Go(func(ctx context.Context) error {
				operator := manifest.DefaultApplyOperator(sc, flags.Wait)
				r, err := operator.Operate(set)
				operator.PrintResult(r)
				if err != nil {
					return err
				}

				// Send result to wait.
				resultChan <- r.UnChanged
				resultChan <- r.Failed
				return nil
			})

			err = wg.Wait()
			if err != nil {
				panic(err)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd, nil
}

func manifestExample(action string) string {
	title := strs.Title(action)

	return fmt.Sprintf(`  # %s the configuration in the walrus-file.yaml 
  $ walrus %s -f walrus-file.yaml

  # %s configurations from a directory containing yaml files
  $ walrus %s -f dir/

  # %s configurations from a directory recursively
  $ walrus %s -f dir/ --recursive

  # %s configurations with a specific project/environment
  $ walrus %s -f dir/ -p my-project -e my-environment
`,
		title, action,
		title, action,
		title, action,
		title, action)
}
