package cmd

import (
	"fmt"

	"github.com/afloesch/megamod/swizzle"
	"github.com/afloesch/semver"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	manifestFile string
	modName      string
	modDesc      string
	gameExe      string
	gameVer      string

	defModName string = "Enter mod name"
	defModDesc string = "Short mod description text."
	defGameExe string = "Game.exe"
	defGameVer string = ">=v0.0.0"
)

func runPrompts() error {
	promptName := promptui.Prompt{
		Label:   "Name",
		Default: modName,
	}
	promptDesc := promptui.Prompt{
		Label:   "Description",
		Default: modDesc,
	}
	promptGame := promptui.Prompt{
		Label:   "Game Executable",
		Default: gameExe,
	}

	var err error
	if modName == defModName {
		modName, err = promptName.Run()
		if err != nil {
			return err
		}
	}
	if modDesc == defModDesc {
		modDesc, err = promptDesc.Run()
		if err != nil {
			return err
		}
	}
	if gameExe == defGameExe {
		gameExe, err = promptGame.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new swizzle manifest file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runPrompts()
		if err != nil {
			fmt.Println(err)
		}

		m := swizzle.New()
		m.Name = modName
		m.Description = modDesc
		m.Game.Executable = gameExe
		m.Game.Version = semver.String(gameVer)

		err = m.WriteFile(manifestFile)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Manifest created:", manifestFile)
		return nil
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&manifestFile, "file", "f", "swizzle.yml", "Swizzle manifest file.")
	initCmd.PersistentFlags().StringVarP(&modName, "name", "n", defModName, "Mod name.")
	initCmd.PersistentFlags().StringVarP(&modDesc, "desc", "d", defModDesc, "Mod short description text.")
	initCmd.PersistentFlags().StringVarP(&gameExe, "game", "g", defGameExe, "The game executable the mod is for.")
	initCmd.PersistentFlags().StringVarP(&gameVer, "version", "v", defGameVer, "The game verion this mod is for. Defaults to all versions.")

	rootCmd.AddCommand(initCmd)
}
