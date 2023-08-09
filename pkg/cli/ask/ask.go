package ask

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// This file is copy and adapt for provide default value and value display for password.
// https://github.com/go-survey/survey/blob/fa37277e6394c29db7bcc94062cb30cd7785a126/password.go#L1

/*
Password is like a normal Input but the text shows up as *'s and there is no default. Response
type is a string.

	password := ""
	prompt := &survey.Password{ Message: "Please type your password" }
	survey.AskOne(prompt, &password)
*/
type Password struct {
	survey.Renderer

	Default        string
	DefaultDisplay string
	Message        string
	Help           string
}

type PasswordTemplateData struct {
	Password

	ShowHelp bool
	Config   *survey.PromptConfig
}

// PasswordQuestionTemplate is a template with color formatting.
// See Documentation: https://github.com/mgutz/ansi#style-format
//
//nolint:lll
var PasswordQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ .Config.HelpInput }} for help]{{color "reset"}} {{end}}
{{- if .Default}}{{color "white"}}({{.DefaultDisplay}}) {{color "reset"}}{{end}}`

func (p *Password) Prompt(config *survey.PromptConfig) (any, error) {
	// Render the question template.
	userOut, _, err := core.RunTemplate(
		PasswordQuestionTemplate,
		PasswordTemplateData{
			Password: *p,
			Config:   config,
		},
	)
	if err != nil {
		return "", err
	}

	if _, err := fmt.Fprint(terminal.NewAnsiStdout(p.Stdio().Out), userOut); err != nil {
		return "", err
	}

	rr := p.NewRuneReader()
	_ = rr.SetTermMode()

	defer func() {
		_ = rr.RestoreTermMode()
	}()

	// No help msg?  Just return any response.
	if p.Help == "" {
		line, err := rr.ReadLine(config.HideCharacter)
		if err != nil {
			return "", err
		}

		// If the line is empty.
		if len(line) == 0 {
			return p.Default, err
		}

		return string(line), err
	}

	cursor := p.NewCursor()

	var line []rune
	// Process answers looking for help prompt answer.
	for {
		line, err = rr.ReadLine(config.HideCharacter)
		if err != nil {
			return string(line), err
		}

		if string(line) == config.HelpInput {
			// Terminal will echo the \n so we need to jump back up one row.
			_ = cursor.PreviousLine(1)

			err = p.Render(
				PasswordQuestionTemplate,
				PasswordTemplateData{
					Password: *p,
					ShowHelp: true,
					Config:   config,
				},
			)
			if err != nil {
				return "", err
			}

			continue
		}

		break
	}

	lineStr := string(line)

	// If the line is empty.
	if len(lineStr) == 0 {
		return p.Default, err
	}

	p.AppendRenderedText(strings.Repeat(string(config.HideCharacter), len(lineStr)))

	return lineStr, err
}

// Cleanup hides the string with a fixed number of characters.
func (p *Password) Cleanup(config *survey.PromptConfig, val any) error {
	return nil
}

// Required does not allow an empty value.
func (p *Password) Required(val any) error {
	// The reflect value of the result.
	value := reflect.ValueOf(val)

	// If the value passed in is the zero value of the appropriate type.
	if isZero(value) && value.Kind() != reflect.Bool && p.Default == "" {
		// This error message should render as capitalized
		//nolint:stylecheck
		return errors.New("Value is required")
	}

	return nil
}

// isZero returns true if the passed value is the zero object.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}

	// Compare the types directly with more general coverage.
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
