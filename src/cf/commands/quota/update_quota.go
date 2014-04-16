package quota

import (
	"cf/api"
	"cf/command_metadata"
	"cf/configuration"
	"cf/flag_helpers"
	"cf/formatters"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
)

type updateQuota struct {
	ui        terminal.UI
	config    configuration.Reader
	quotaRepo api.QuotaRepository
}

func NewUpdateQuota(ui terminal.UI, config configuration.Reader, quotaRepo api.QuotaRepository) *updateQuota {
	return &updateQuota{
		ui:        ui,
		config:    config,
		quotaRepo: quotaRepo,
	}
}

func (command *updateQuota) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "update-quota",
		Description: "Update an existing resource quota",
		Usage:       "CF_NAME update-quota QUOTA [-m MEMORY] [-n NEW_NAME] [-r ROUTES] [-s SERVICE_INSTANCES]",
		Flags: []cli.Flag{
			flag_helpers.NewStringFlag("m", "Total amount of memory (e.g. 1024M, 1G, 10G)"),
			flag_helpers.NewStringFlag("n", "New name"),
			flag_helpers.NewIntFlag("r", "Total number of routes"),
			flag_helpers.NewIntFlag("s", "Total number of service instances"),
		},
	}
}

func (cmd *updateQuota) GetRequirements(requirementsFactory requirements.Factory, context *cli.Context) ([]requirements.Requirement, error) {
	if len(context.Args()) != 1 {
		cmd.ui.FailWithUsage(context, "update-quota")
	}

	return []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}, nil
}

func (cmd *updateQuota) Run(c *cli.Context) {
	oldQuotaName := c.Args()[0]
	quota, err := cmd.quotaRepo.FindByName(oldQuotaName)

	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	if c.String("m") != "" {
		memory, err := formatters.ToMegabytes(c.String("m"))

		if err != nil {
			cmd.ui.FailWithUsage(c, "update-quota")
		}

		quota.MemoryLimit = memory
	}

	if c.String("n") != "" {
		quota.Name = c.String("n")
	}

	if c.IsSet("s") {
		quota.ServicesLimit = c.Int("s")
	}

	if c.IsSet("r") {
		quota.RoutesLimit = c.Int("r")
	}

	cmd.ui.Say("Updating quota %s as %s...",
		terminal.EntityNameColor(oldQuotaName),
		terminal.EntityNameColor(cmd.config.Username()))

	err = cmd.quotaRepo.Update(quota)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}
	cmd.ui.Ok()
}
