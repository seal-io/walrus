package doc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/seal-io/seal/pkg/cli/api"
)

// This file is copy and adapt for generate CLI docs with customize format.
// https://github.com/spf13/cobra/blob/v1.7.0/doc/md_docs.go

// GenMarkdownTree for input command with sealFilePrepender and sealDocLinkHandler.
func GenMarkdownTree(cmd *cobra.Command, dir string) error {
	return GenMarkdownTreeCustom(cmd, dir, sealFilePrepender, sealDocLinkHandler)
}

// GenMarkdownTreeCustom is the same as GenMarkdownTree, but
// with custom filePrepender and linkHandler.
func GenMarkdownTreeCustom(
	cmd *cobra.Command,
	dir string,
	filePrepender func(*cobra.Command, string) (string, error),
	linkHandler func(*cobra.Command, string) string,
) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	filename, err := filePrepender(cmd, dir)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}

	return nil
}

// GenMarkdownCustom creates custom markdown output.
func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(*cobra.Command, string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("# " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")

	if len(cmd.Long) > 0 {
		buf.WriteString("## Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("## Examples\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd); err != nil {
		return err
	}

	if hasSeeAlso(cmd) {
		buf.WriteString("## SEE ALSO\n\n")

		if cmd.HasParent() {
			link := linkHandler(cmd, "parent")
			buf.WriteString(link)
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			link := linkHandler(child, "child")
			buf.WriteString(link)
		}

		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)

	return err
}

// printOptions print the flags and parent's flags.
func printOptions(buf *bytes.Buffer, cmd *cobra.Command) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)

	if flags.HasAvailableFlags() {
		buf.WriteString("## Options\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)

	if parentFlags.HasAvailableFlags() {
		buf.WriteString("## Options inherited from parent commands\n\n```\n")
		parentFlags.PrintDefaults()
		buf.WriteString("```\n\n")
	}

	return nil
}

// hasSeeAlso check command has extra related info.
func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		return true
	}

	return false
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// sealDocLinkHandler customize the doc link for seal docs.
func sealDocLinkHandler(cmd *cobra.Command, commandType string) string {
	switch commandType {
	case "parent":
		// Parent.
		var (
			parent = cmd.Parent()
			pname  = parent.CommandPath()
			pres   = parent.Annotations[api.AnnResourceName]
			cres   = cmd.Annotations[api.AnnResourceName]
			link   = strings.ReplaceAll(pname, " ", "_")
		)

		if pres == cres {
			return fmt.Sprintf("* [%s](%s)\t - %s\n", pname, link, parent.Short)
		}

		return fmt.Sprintf("* [%s](../%s)\t - %s\n", pname, link, parent.Short)
	default:
		//  Child.
		var (
			parent = cmd.Parent()
			cname  = parent.CommandPath() + " " + cmd.Name()
			pres   = parent.Annotations[api.AnnResourceName]
			cres   = cmd.Annotations[api.AnnResourceName]
			link   = strings.ReplaceAll(cname, " ", "_")
		)

		if pres == cres {
			return fmt.Sprintf("* [%s](%s)\t - %s\n", cname, link, cmd.Short)
		}

		return fmt.Sprintf("* [%s](%s/%s)\t - %s\n", cname, cres, link, parent.Short)
	}
}

// sealFilePrepender customize the doc file path for seal docs.
func sealFilePrepender(cmd *cobra.Command, dir string) (string, error) {
	res := cmd.Annotations[api.AnnResourceName]

	err := os.MkdirAll(filepath.Join(dir, res), 0o755)
	if err != nil {
		return "", err
	}

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"

	return filepath.Join(dir, res, basename), nil
}
