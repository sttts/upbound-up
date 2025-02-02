// Copyright 2022 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package organization

import (
	"context"
	"fmt"

	"github.com/pterm/pterm"

	"github.com/upbound/up-sdk-go/service/organizations"

	"github.com/upbound/up/internal/input"
)

// BeforeApply sets default values for the delete command, before assignment and validation.
func (c *deleteCmd) BeforeApply() error {
	c.prompter = input.NewPrompter()
	return nil
}

// AfterApply accepts user input by default to confirm the delete operation.
func (c *deleteCmd) AfterApply(p pterm.TextPrinter) error {
	if c.Force {
		return nil
	}

	confirm, err := c.prompter.Prompt("Are you sure you want to delete this organization? [y/n]", false)
	if err != nil {
		return err
	}

	if input.InputYes(confirm) {
		p.Printfln("Deleting organization %s. This cannot be undone.", c.Name)
		return nil
	}

	return fmt.Errorf("operation canceled")
}

// deleteCmd deletes an organization on Upbound.
type deleteCmd struct {
	prompter input.Prompter

	Name string `arg:"" required:"" help:"Name of organization." predictor:"orgs"`

	Force bool `help:"Force deletion of the organization." default:"false"`
}

// Run executes the delete command.
func (c *deleteCmd) Run(p pterm.TextPrinter, oc *organizations.Client) error {
	id, err := oc.GetOrgID(context.Background(), c.Name)
	if err != nil {
		return err
	}
	if err := oc.Delete(context.Background(), id); err != nil {
		return err
	}
	p.Printfln("%s deleted", c.Name)
	return nil
}
