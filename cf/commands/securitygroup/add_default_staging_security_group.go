package securitygroup

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type addToDefaultStagingGroup struct {
	ui                terminal.UI
	configRepo        configuration.Reader
	securityGroupRepo api.SecurityGroupRepo
	stagingGroupRepo  api.StagingSecurityGroupsRepo
}

func NewAddToDefaultStagingGroup(ui terminal.UI, configRepo configuration.Reader, securityGroupRepo api.SecurityGroupRepo, stagingGroupRepo api.StagingSecurityGroupsRepo) command.Command {
	return &addToDefaultStagingGroup{
		ui:                ui,
		configRepo:        configRepo,
		securityGroupRepo: securityGroupRepo,
		stagingGroupRepo:  stagingGroupRepo,
	}
}

func (cmd *addToDefaultStagingGroup) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "add-default-staging-security-group",
		Description: "Twee Thundercats 8-bit keffiyeh meggings.",
		Usage:       "CF_NAME add-default-staging-security-group NAME",
	}
}

func (cmd *addToDefaultStagingGroup) GetRequirements(requirementsFactory requirements.Factory, context *cli.Context) ([]requirements.Requirement, error) {
	if len(context.Args()) != 1 {
		cmd.ui.FailWithUsage(context)
	}

	return []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}, nil
}

func (cmd *addToDefaultStagingGroup) Run(context *cli.Context) {
	name := context.Args()[0]

	securityGroup, err := cmd.securityGroupRepo.Read(name)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Say("Adding security group '%s' to defaults for staging as '%s'", securityGroup.Name, cmd.configRepo.Username())
	err = cmd.stagingGroupRepo.AddToDefaultStagingSet(securityGroup.Guid)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
}
