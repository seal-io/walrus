package clis

import (
	"reflect"

	"github.com/urfave/cli/v2"

	"github.com/seal-io/seal/utils/strs"
)

func MutateAppEnvs(app *App) {
	app.Commands = MutateCommandsEnvs(app.Name, app.Commands)
	app.Flags = MutateFlagsEnvs(app.Name, app.Flags)
}

func MutateCommandsEnvs(prefix string, cmds Commands) Commands {
	if len(cmds) == 0 {
		return cmds
	}

	for i := 0; i < len(cmds); i++ {
		if len(cmds[i].Subcommands) != 0 {
			n := cmds[i].Name
			if prefix != "" {
				n = prefix + "-" + n
			}
			cmds[i].Subcommands = MutateCommandsEnvs(n, cmds[i].Subcommands)
		}
		cmds[i].Flags = MutateFlagsEnvs(prefix, cmds[i].Flags)
	}

	return cmds
}

func MutateFlagsEnvs(prefix string, flags Flags) Flags {
	if len(flags) == 0 {
		return flags
	}

	for i := range flags {
		fv := reflect.ValueOf(flags[i])
		for fv.Kind() == reflect.Pointer || fv.Kind() == reflect.Interface {
			fv = fv.Elem()
		}

		fvf := fv.FieldByName("EnvVars")
		if !fvf.IsValid() || fvf.Kind() != reflect.Slice {
			continue
		}

		n := flags[i].Names()[0]
		if prefix != "" {
			n = prefix + "-" + n
		}
		n = strs.UnderscoreUpper(n)

		var (
			fvt = fv.Type()
			cfv = reflect.New(fvt)
		)

		for k := 0; k < fv.NumField(); k++ {
			cfvf := cfv.Elem().Field(k)
			if !cfvf.CanSet() {
				continue
			}

			if fvt.Field(k).Name != "EnvVars" || !fv.Field(k).IsZero() {
				cfvf.Set(fv.Field(k))
				continue
			}

			cfvf.Set(reflect.ValueOf([]string{n}))
		}
		flags[i] = cfv.Interface().(cli.Flag)
	}

	return flags
}
