package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
	"github.com/seal-io/walrus/utils/gopool"
)

func Delete(sc *config.Config) (*cobra.Command, error) {
	flags := &manifest.Option{}

	use := "delete"
	cmd := &cobra.Command{
		Use:     use,
		GroupID: common.GroupAdvanced.ID,
		Example: manifestExample("delete"),
		Short:   "Delete resource using the provided file path or folder.",
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
					waiter := manifest.DeleteWaiter(sc, flags.Timeout)
					_, err = waiter.Wait(ctx, set, resultChan)
					return err
				})
			}

			// Delete the files.
			wg.Go(func(ctx context.Context) error {
				operator := manifest.DefaultDeleteOperator(sc, flags.Wait)
				r, err := operator.Operate(set)
				operator.PrintResult(r)
				if err != nil {
					return err
				}

				// Send result to wait.
				resultChan <- r.NotFound
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
