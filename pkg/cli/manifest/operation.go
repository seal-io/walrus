package manifest

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/config"
)

type OperateOption struct {
	// Context flags.
	config.ScopeContext

	// File path or folder path.
	Filenames []string `json:"filenames,omitempty"`

	// Recursive apply.
	Recursive bool `json:"recursive,omitempty"`
}

func (f *OperateOption) AddFlags(cmd *cobra.Command) {
	f.ScopeContext.AddFlags(cmd)

	cmd.Flags().StringSliceVarP(&f.Filenames, "filenames", "f", nil, "File path or folder path")
	cmd.Flags().BoolVarP(&f.Recursive, "recursive", "r", false, "Recursive apply")
}

func (f *OperateOption) Context(sc *config.Config) config.ScopeContext {
	c := config.ScopeContext{
		Project:     sc.Project,
		Environment: sc.Environment,
	}

	if f.Project != "" {
		c.Project = f.Project
	}

	if f.Environment != "" {
		c.Environment = f.Environment
	}

	return c
}

func Apply(sc *config.Config, opt *OperateOption) error {
	objs, err := LoadObjects(opt.Context(sc), opt.Filenames, opt.Recursive)
	if err != nil {
		return err
	}

	for _, group := range APIGroupCreateSequence {
		resObjs := objs[group]
		if len(resObjs) == 0 {
			continue
		}

		toPatch, toCreate, err := GetObjects(sc, group, resObjs)
		if err != nil {
			return err
		}

		// Patch.
		err = PatchObjects(sc, group, toPatch)
		if err != nil {
			return err
		}

		// Batch create.
		err = BatchCreateObjects(sc, group, toCreate)
		if err != nil {
			return err
		}
	}

	return nil
}

func Delete(sc *config.Config, opt *OperateOption) error {
	objs, err := LoadObjects(opt.Context(sc), opt.Filenames, opt.Recursive)
	if err != nil {
		return err
	}

	// Delete in reverse order.
	for i := len(APIGroupCreateSequence) - 1; i >= 0; i-- {
		group := APIGroupCreateSequence[i]

		resObjs := objs[group]
		if len(resObjs) == 0 {
			continue
		}

		toDelete, _, err := GetObjects(sc, group, resObjs)
		if err != nil {
			return err
		}

		if len(toDelete) == 0 {
			continue
		}

		err = DeleteObjects(sc, group, toDelete)
		if err != nil {
			return err
		}
	}

	return nil
}
