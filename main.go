package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	discovery "github.com/rancher-sandbox/upgradechannel-discovery/pkg/discovery"
	github "github.com/rancher-sandbox/upgradechannel-discovery/pkg/discovery/type/github"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:        "upgradechannel-discovery",
		Version:     "", // TODO: bind internal.Version to CI while building with ldflags
		Author:      "",
		Usage:       "",
		Description: "",
		Copyright:   "",
		Commands: []cli.Command{
			{
				Name: "github",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "image-prefix",
						Value:  "",
						EnvVar: "IMAGE_PREFIX",
						Usage:  "Image prefix to use when returning json data",
					},
					&cli.StringFlag{
						Name:   "github-token",
						EnvVar: "GITHUB_TOKEN",
						Value:  "",
						Usage:  "Github token used to identify against github for fetching releases",
					},
					&cli.StringFlag{
						Name:   "output-file",
						EnvVar: "OUTPUT_FILE",
						Value:  "",
						Usage:  "File to output the resulting json from",
					},
					&cli.StringFlag{
						Name:   "version-name-prefix",
						EnvVar: "VERSION_NAME_PREFIX",
						Value:  "",
						Usage:  "Version name prefix",
					},
					&cli.StringFlag{
						Name:   "version-name-suffix",
						EnvVar: "VERSION_NAME_SUFFIX",
						Value:  "",
						Usage:  "Version name suffix",
					},
					&cli.StringFlag{
						Name:   "repository",
						EnvVar: "REPOSITORY",
						Value:  "rancher-sandbox/os2",
						Usage:  "Github repository to scan releases against",
					},
				},
				Action: func(c *cli.Context) error {
					outFile := c.String("output-file")

					rf, err := github.NewReleaseFinder(
						github.WithContext(context.Background()),
						github.WithRepository(c.String("repository")),
						github.WithToken(c.String("github-token")),
						github.WithVersionNamePrefix(c.String("version-name-prefix")),
						github.WithVersionNameSuffix(c.String("version-name-suffix")),
						github.WithBaseImage(c.String("image-prefix")),
					)

					if err != nil {
						return err
					}

					b, err := discovery.Versions(rf)
					if err != nil {
						return err
					}

					if outFile != "" {
						return ioutil.WriteFile(outFile, b, os.ModePerm)
					} else {
						fmt.Print(string(b))
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
