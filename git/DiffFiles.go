package git

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"
)

func DiffFiles(from, to string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("%v..%v", from, to))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(stdout)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	files := []string{}

	for {
		file, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		file = path.Clean(strings.TrimSpace(file))
		if !strings.HasPrefix(file, "/") {
			file = "/" + file
		}

		files = append(files, file)
	}

	return files, nil
}
