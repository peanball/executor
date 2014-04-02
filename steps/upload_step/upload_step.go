package upload_step

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"

	steno "github.com/cloudfoundry/gosteno"
	"github.com/vito/gordon"

	"github.com/cloudfoundry-incubator/executor/backend_plugin"
	"github.com/cloudfoundry-incubator/executor/compressor"
	"github.com/cloudfoundry-incubator/executor/log_streamer"
	"github.com/cloudfoundry-incubator/executor/uploader"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type UploadStep struct {
	containerHandle string
	model           models.UploadAction
	uploader        uploader.Uploader
	compressor      compressor.Compressor
	tempDir         string
	backendPlugin   backend_plugin.BackendPlugin
	wardenClient    gordon.Client
	streamer        log_streamer.LogStreamer
	logger          *steno.Logger
}

func New(
	containerHandle string,
	model models.UploadAction,
	uploader uploader.Uploader,
	compressor compressor.Compressor,
	tempDir string,
	wardenClient gordon.Client,
	streamer log_streamer.LogStreamer,
	logger *steno.Logger,
) *UploadStep {
	return &UploadStep{
		containerHandle: containerHandle,
		model:           model,
		uploader:        uploader,
		compressor:      compressor,
		tempDir:         tempDir,
		wardenClient:    wardenClient,
		streamer:        streamer,
		logger:          logger,
	}
}

func (step *UploadStep) Perform() error {
	step.logger.Infod(
		map[string]interface{}{
			"handle": step.containerHandle,
		},
		"runonce.handle.upload-step",
	)

	tempDir, err := ioutil.TempDir(step.tempDir, "upload")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempDir)

	currentUser, err := user.Current()
	if err != nil {
		panic("existential failure: " + err.Error())
	}

	_, err = step.wardenClient.CopyOut(
		step.containerHandle,
		step.model.From,
		tempDir,
		currentUser.Username,
	)
	if err != nil {
		return err
	}

	url, err := url.Parse(step.model.To)
	if err != nil {
		return err
	}

	tempFile, err := ioutil.TempFile(step.tempDir, "compressed")
	if err != nil {
		return err
	}

	tempFile.Close()

	finalFileLocation := tempFile.Name()

	defer os.RemoveAll(finalFileLocation)

	err = step.compressor.Compress(tempDir, finalFileLocation)
	if err != nil {
		return err
	}

	if step.streamer != nil {
		step.streamer.StreamStdout(fmt.Sprintf("Uploading %s", step.model.Name))
	}

	return step.uploader.Upload(finalFileLocation, url)
}

func (step *UploadStep) Cancel() {}

func (step *UploadStep) Cleanup() {}
