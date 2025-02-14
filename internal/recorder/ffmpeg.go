package recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shinosaki/namagent/internal/recorder/types"
)

func FFmpeg(
	url string,
	cookies []types.WS_StreamCookie,
	ffmpegPath string,
	outputPath string,
) *exec.Cmd {
	// mkdir
	if abs, err := filepath.Abs(outputPath); err != nil {
		log.Fatalf("invalid output path: %v", err)
	} else {
		dir := filepath.Dir(abs)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("failed to creating output path dir: %v", err)
		}
	}

	command := []string{ffmpegPath}

	// if exists cookies
	if len(cookies) > 0 {
		formatted := make([]string, len(cookies))
		for i, c := range cookies {
			formatted[i] = fmt.Sprintf("%s=%s; domain=%s; path=%s",
				c.Name,
				c.Value,
				c.Domain,
				c.Path,
			)
		}

		command = append(command, "-cookies", strings.Join(formatted, ";\n"))
	}

	command = append(command,
		"-nostdin",
		"-loglevel", "warning",
		"-i", url,
		"-c", "copy",
		outputPath,
	)

	proc := exec.Command(command[0], command[1:]...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc
}
