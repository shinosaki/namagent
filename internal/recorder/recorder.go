package recorder

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shinosaki/namagent/utils"
)

func Recorder(
	sc *utils.SignalContext,
	client *http.Client,
	programId string,
	ffmpegPath string,
	outputPath string,
) error {
	// Exists check
	if !utils.IsExecutable(ffmpegPath, "-version") {
		return fmt.Errorf("does not executable ffmpeg")
	}

	// output path
	if utils.IsExistsFile(outputPath) {
		return fmt.Errorf("already exists output path")
	}

	program, err := FetchProgramData(programId, client)
	if err != nil {
		return err
	}

	log.Println("Program Status:", program.Program.Status)
	log.Println("Program Stream Type:", program.Stream.Type)
	log.Println("Program Supplier ID:", program.Program.Supplier.ProgramProviderId)
	log.Println("Program Supplier Name:", program.Program.Supplier.Name)
	log.Println("Program Title:", program.Program.Title)

	if err := WebSocket(sc, program, ffmpegPath, outputPath); err != nil {
		return err
	}

	return nil
}
