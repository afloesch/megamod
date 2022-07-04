package cmd

import (
	"context"

	"github.com/afloesch/megamod/swizzle"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	repo    string
	repoVer string
)

func addPrompt() error {
	prompt := promptui.Prompt{
		Label:   "Repo",
		Default: repo,
	}

	var err error
	repo, err = prompt.Run()
	if err != nil {
		return err
	}

	return nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new mod.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		err := addPrompt()
		if err != nil {
			return err
		}

		mod, err := swizzle.New().ReadFile("./swizzle.yml")
		if err != nil {
			return err
		}

		ver := repoVer
		if ver == "latest" {
			rel, err := swizzle.Repo(repo).LatestManifest(ctx)
			if err != nil {
				return err
			}

			ver = string(rel.Version)
		}

		err = mod.AddDependency(ctx, repo, ver)
		if err != nil {
			return err
		}

		return mod.WriteFile("./swizzle.yml")
	},
}

func init() {
	addCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Github repository.")
	addCmd.PersistentFlags().StringVarP(&repoVer, "version", "v", "latest", "Release version.")
	rootCmd.AddCommand(addCmd)
}
