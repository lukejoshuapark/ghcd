package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/lukejoshuapark/ghcd/config"
	"github.com/lukejoshuapark/ghcd/git"
	"github.com/lukejoshuapark/ghcd/github"
	"github.com/urfave/cli/v2"
)

func detect(cctx *cli.Context) error {
	file := cctx.String("f")
	token := cctx.String("token")
	repository := cctx.String("repository")
	start := cctx.String("start")
	end := cctx.String("end")

	var filesDiffFiles []string

	cfg, err := config.FromFile(file)
	if err != nil {
		return err
	}

	results := map[string]bool{}

	for targetName, target := range cfg.Detect {
		switch target.Mode {
		case "FilesDiff":
			if filesDiffFiles == nil {
				filesDiffFiles, err = git.DiffFiles(start, end)
				if err != nil {
					return fmt.Errorf(`failed on target "%v": %w`, targetName, err)
				}
			}

			results[targetName] = hasChanges(&target, filesDiffFiles)
		case "EnvironmentDiff":
			environmentCommit, err := github.EnvironmentCommit(repository, target.Enviroment, token)
			if err != nil {
				return fmt.Errorf(`failed on target "%v": %w`, targetName, err)
			}

			environmentDiffFiles, err := git.DiffFiles(environmentCommit, end)
			if err != nil {
				return fmt.Errorf(`failed on target "%v": %w`, targetName, err)
			}

			results[targetName] = hasChanges(&target, environmentDiffFiles)
		default:
			return fmt.Errorf(`failed on target "%v": unknown mode "%v"`, targetName, target.Mode)
		}
	}

	for targetName, result := range results {
		fmt.Printf("%v=%v\n", targetName, result)
	}

	return nil
}

func hasChanges(target *config.Target, diffFiles []string) bool {
	for _, checkPath := range target.Paths {
		cleaned := path.Clean(checkPath)
		cleaned = strings.ReplaceAll(cleaned, "\\", "/")
		if !strings.HasPrefix(cleaned, "/") {
			cleaned = "/" + cleaned
		}
		if !strings.HasSuffix(cleaned, "/") {
			cleaned += "/"
		}

		for _, diffFile := range diffFiles {
			if strings.HasPrefix(diffFile, cleaned) {
				return true
			}
		}
	}

	return false
}
