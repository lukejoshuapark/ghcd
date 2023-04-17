package github

import (
	"fmt"
	"io"
	"os"
)

func WriteResults(results map[string]*string) error {
	if err := writeOutputs(results); err != nil {
		return err
	}

	if err := writeSummary(results); err != nil {
		return err
	}

	return nil
}

func writeOutputs(results map[string]*string) error {
	var w io.Writer

	githubOutputFile, ok := os.LookupEnv("GITHUB_OUTPUT")
	if ok {
		f, err := os.OpenFile(githubOutputFile, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("GITHUB_OUTPUT was provided but could not be opened for appending: %w", err)
		}
		defer f.Close()

		w = f
	} else {
		w = os.Stdout
	}

	for targetName, fileName := range results {
		_, err := fmt.Fprintf(w, "%v=%v\n", targetName, fileName != nil)
		if err != nil {
			return fmt.Errorf("failed to write to result stream: %w", err)
		}
	}

	return nil
}

func writeSummary(results map[string]*string) error {
	var w io.Writer

	githubSummaryFile, ok := os.LookupEnv("GITHUB_STEP_SUMMARY")
	if ok {
		f, err := os.OpenFile(githubSummaryFile, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("GITHUB_STEP_SUMMARY was provided but could not be opened for appending: %w", err)
		}
		defer f.Close()

		w = f
	} else {
		w = os.Stderr
	}

	_, err := fmt.Fprintf(w, "### 🌾 Change Detection Results\n")
	if err != nil {
		return fmt.Errorf("failed to write to summary stream: %w", err)
	}

	for targetName, fileName := range results {
		if fileName != nil {
			_, err := fmt.Fprintf(w, "- 🚀 **%v** was triggered by _%v_\n", targetName, *fileName)
			if err != nil {
				return fmt.Errorf("failed to write to summary stream: %w", err)
			}
		} else {
			_, err := fmt.Fprintf(w, "- 💤 **%v** was not triggered\n", targetName)
			if err != nil {
				return fmt.Errorf("failed to write to summary stream: %w", err)
			}
		}
	}

	return nil
}
