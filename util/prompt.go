package util

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
)

type PromptContent struct {
	Label     string
	ErrorMsg  string
	Validate  promptui.ValidateFunc
	Templates *promptui.PromptTemplates
	Default   string
	IsConfirm bool
	AllowEdit bool
}

type PromptContentOption func(pc *PromptContent)

var DefaultTemplates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

var ValidateWrapper = func(errorMsg string, validate Validator) promptui.ValidateFunc {
	return func(input string) error {
		if err := validate(input); err != nil {
			if errorMsg != "" {
				return errors.New(errorMsg)
			}
			return err
		}
		return nil
	}
}

func OptionSetValidator(validator Validator) PromptContentOption {
	return func(pc *PromptContent) {
		pc.Validate = ValidateWrapper(pc.ErrorMsg, validator)
	}
}

func OptionSetTemplates(template *promptui.PromptTemplates) PromptContentOption {
	return func(pc *PromptContent) {
		pc.Templates = template
	}
}

func OptionSetDefault(defaultValue string) PromptContentOption {
	return func(pc *PromptContent) {
		pc.Default = defaultValue
	}
}

func OptionSetIsConfirm(isConfirm bool) PromptContentOption {
	return func(pc *PromptContent) {
		pc.IsConfirm = isConfirm
	}
}

func NewPromptContent(label string, errorMsg string, options ...PromptContentOption) *PromptContent {
	pc := &PromptContent{
		Label:     label,
		ErrorMsg:  errorMsg,
		Validate:  ValidateWrapper(errorMsg, RequiredValidator),
		Templates: DefaultTemplates,
	}
	for _, option := range options {
		option(pc)
	}
	// Put cursor at the end of the line when providing default value
	// https://github.com/manifoldco/promptui/issues/146
	if pc.Default != "" {
		pc.AllowEdit = true
	}
	return pc
}

func PromptGetInput(pc *PromptContent) string {
	prompt := promptui.Prompt{
		Label:     pc.Label,
		Templates: pc.Templates,
		Validate:  pc.Validate,
		Default:   pc.Default,
		IsConfirm: pc.IsConfirm,
		AllowEdit: pc.AllowEdit,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}
