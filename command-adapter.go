package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// GetApplication assembles application.
func GetApplication(args []string) Application {
	envValues := NewEnvValues()
	var fileFlag string
	app := cli.App{
		Flags: configFile(envValues, &fileFlag),
		Commands: []*cli.Command{
			profileCommand(&envValues, &fileFlag),
		},
	}
	return &CliApp{
		Args: args,
		App:  app,
	}
}

func configFile(envValues EnvValues, fileFlag *string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Aliases:     []string{"file", "f"},
			Usage:       "configuration file",
			Required:    false,
			Value:       string(envValues.DefaultProfileFile()),
			Destination: fileFlag,
		},
	}
}

// CliApp is wrapper of cli.App
type CliApp struct {
	Args []string
	App  cli.App
}

// Start starts application
func (app *CliApp) Start() error {
	return app.App.Run(app.Args)
}

func profileCommand(envValues *EnvValues, fileFlag *string) *cli.Command {
	var name string
	var token string
	var dir string
	return &cli.Command{
		Name:    "profile",
		Aliases: []string{"p"},
		Usage:   "add or update profile configuration",
		Action: func(context *cli.Context) error {
			return profileCommandAction(envValues, *fileFlag, name, token, dir)
		},
		Flags: []cli.Flag{
			profileFlag(&name),
			tokenFlag(&token),
			destinationDirectoryFlag(&dir),
		},
	}
}

func profileCommandAction(envValues *EnvValues, fileFlag string, name string, token string, dir string) error {
	ctx, err := envValues.NewContext(ProfileFile(fileFlag))
	if err != nil {
		return fmt.Errorf("ProfileCommand_NewContext: %w", err)
	}
	command := NewProfileCommand(&name, &token, &dir)
	err = command.Run(ctx)
	if err != nil {
		return fmt.Errorf("ProfileCommandAction: %w", err)
	}
	return nil
}

func profileFlag(name *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "profile",
		Aliases:     []string{"p"},
		Usage:       "A name of profile",
		Required:    false,
		Value:       "default",
		Destination: name,
	}
}

func tokenFlag(token *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "token",
		Aliases:     []string{"t"},
		Usage:       "GitHub Access Token for this profile",
		Required:    false,
		DefaultText: "",
		Destination: token,
	}
}

func destinationDirectoryFlag(dir *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "dir",
		Aliases:     []string{"d"},
		Usage:       "Destination directory for this profile",
		Required:    false,
		Value:       "",
		Destination: dir,
	}
}
