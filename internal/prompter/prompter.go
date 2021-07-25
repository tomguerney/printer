package prompter

import (
	"github.com/manifoldco/promptui"
)

//Prompter gets input from the user
type Prompter struct {
}

// New returns a new prompter
func New() *Prompter {
	return &Prompter{}
}

//Select prompts the user to select an option
func (p *Prompter) Select(label string, table []string) (i int, err error) {

	prompt := promptui.Select{
		Label: label,
		Items: table,
	}

	i, _, err = prompt.Run()

	if err != nil {
		return 0, err
	}

	return i, nil
}

// Confirm gets the user to confirm yes or no
func (p *Prompter) Confirm(label string) bool {

	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()

	return err == nil

}

// In gets the user to provide some input
func (p *Prompter) In(label string) (s string, err error) {

	prompt := promptui.Prompt{
		Label: label,
	}

	input, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return input, nil

}
