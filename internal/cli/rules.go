package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/auth0/auth0-cli/internal/ansi"
	"github.com/auth0/auth0-cli/internal/auth0"
	"github.com/auth0/auth0-cli/internal/prompt"
	"github.com/spf13/cobra"
	"gopkg.in/auth0.v5/management"
)

func rulesCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "manage rules for clients.",
	}

	cmd.SetUsageTemplate(resourceUsageTemplate())
	cmd.AddCommand(listRulesCmd(cli))
	cmd.AddCommand(enableRuleCmd(cli))
	cmd.AddCommand(disableRuleCmd(cli))
	cmd.AddCommand(createRulesCmd(cli))
	cmd.AddCommand(deleteRulesCmd(cli))

	return cmd
}

func listRulesCmd(cli *cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists your rules",
		Long:  `Lists the rules in your current tenant.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var rules *management.RuleList
			err := ansi.Spinner("Getting rules", func() error {
				var err error
				rules, err = getRules(cli)
				return err
			})

			if err != nil {
				return err
			}

			cli.renderer.RulesList(rules)
			return nil
		},
	}

	return cmd
}

func enableRuleCmd(cli *cli) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ansi.Spinner("Enabling rule", func() error {
				var err error
				data, err := getRules(cli)
				if err != nil {
					return err
				}

				rule := findRuleByName(name, data.Rules)
				if rule != nil {
					err := enableRule(rule, cli)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("No rule found with name: %q", name)
				}
				return nil
			})

			if err != nil {
				return err
			}

			// @TODO Only display modified rules
			rules, err := getRules(cli)

			if err != nil {
				return err
			}

			cli.renderer.RulesList(rules)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "rule name")
	mustRequireFlags(cmd, "name")
	return cmd
}

func disableRuleCmd(cli *cli) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ansi.Spinner("Disabling rule", func() error {
				var err error
				data, err := getRules(cli)
				if err != nil {
					return err
				}

				rule := findRuleByName(name, data.Rules)
				if rule != nil {
					if err := disableRule(rule, cli); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("No rule found with name: %q", name)
				}
				return nil
			})

			if err != nil {
				return err
			}

			// @TODO Only display modified rules
			rules, err := getRules(cli)

			if err != nil {
				return err
			}

			cli.renderer.RulesList(rules)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "rule name")
	mustRequireFlags(cmd, "name")

	return cmd
}

func createRulesCmd(cli *cli) *cobra.Command {
	var flags struct {
		Name    string
		Script  string
		Order   int
		Enabled bool
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rule",
		Long: `Create a new rule:

    auth0 rules create --name "My Rule" --script "function (user, context, callback) { console.log( 'Hello, world!' ); return callback(null, user, context); }"
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			checkFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !hasFlags(cmd) {
				name := prompt.TextInput(
					"name", "Name:", 
					"Name of the rule. You can change the rule name later in the rule settings.", 
					"", 
					true)

				script := prompt.TextInput("script", "Script:", "Script of the rule.", "", true)
				order := prompt.TextInput("order", "Order:", "Order of the rule.", "0", false)
				enabled := prompt.BoolInput("enabled", "Enabled:", "Enable the rule.", false)

				if err := prompt.Ask([]*survey.Question {name, script, order, enabled}, &flags); err != nil {
					return err
				}
			}

			r := &management.Rule{
				Name:    &flags.Name,
				Script:  &flags.Script,
				Order:   &flags.Order,
				Enabled: &flags.Enabled,
			}

			err := ansi.Spinner("Creating rule", func() error {
				return cli.api.Rule.Create(r)
			})

			if err != nil {
				return err
			}

			cli.renderer.Infof("Your rule `%s` was successfully created.", flags.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "Name of this rule (required)")
	cmd.Flags().StringVarP(&flags.Script, "script", "s", "", "Code to be executed when this rule runs (required)")
	cmd.Flags().IntVarP(&flags.Order, "order", "o", 0, "Order that this rule should execute in relative to other rules. Lower-valued rules execute first.")
	cmd.Flags().BoolVarP(&flags.Enabled, "enabled", "e", false, "Whether the rule is enabled (true), or disabled (false).")
	mustRequireFlags(cmd, "name", "script")
	return cmd
}

func deleteRulesCmd(cli *cli) *cobra.Command {
	var flags struct {
		id    string
		force bool
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a rule",
		Long: `Delete a rule:

	auth0 rules delete --id "12345" --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !flags.force {
				if confirmed := prompt.Confirm("Are you sure you want to proceed?"); !confirmed {
					return nil
				}
			}

			// TODO: Should add validation of rule
			r := &management.Rule{ID: &flags.id}

			err := ansi.Spinner("Deleting rule", func() error {
				return cli.api.Rule.Delete(*r.ID)
			})

			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flags.id, "id", "", "ID of the rule to delete (required)")
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Do not ask for confirmation.")
	mustRequireFlags(cmd, "id")

	return cmd
}

// @TODO move to rules package
func getRules(cli *cli) (list *management.RuleList, err error) {
	return cli.api.Rule.List()
}

func findRuleByName(name string, rules []*management.Rule) *management.Rule {
	for _, r := range rules {
		if auth0.StringValue(r.Name) == name {
			return r
		}
	}
	return nil
}

func enableRule(rule *management.Rule, cli *cli) error {
	return cli.api.Rule.Update(rule.GetID(), &management.Rule{Enabled: auth0.Bool(true)})
}

func disableRule(rule *management.Rule, cli *cli) error {
	return cli.api.Rule.Update(rule.GetID(), &management.Rule{Enabled: auth0.Bool(false)})
}
