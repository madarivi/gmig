package gmig

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/emicklei/tre"
)

// GCS = Google Cloud Storage
type GCS struct {
	onDiskAccess FileStateProvider
}

// NewGCS returns a new GCS
func NewGCS(cfg Config) GCS {
	return GCS{onDiskAccess: FileStateProvider{Configuration: cfg}}
}

// LoadState implements StateProvider
func (g GCS) LoadState() (string, error) {
	cmdline := []string{"gsutil", "-q", "cp",
		"gs://" + filepath.Join(g.Config().Bucket, LastMigrationObjectName),
		LastMigrationObjectName}
	if err := g.gsutil(cmdline); err != nil {
		// see if there was no last applied state
		if strings.Contains(err.Error(), "No URLs matched") { // lame detection method TODO
			if g.Config().Verbose {
				log.Println("no last applied migration found.")
			}
			return "", nil
		}
		return "", err
	}
	return g.onDiskAccess.LoadState()
}

// SaveState implements StateProvider
func (g GCS) SaveState(filename string) error {
	if err := g.onDiskAccess.SaveState(filename); err != nil {
		return err
	}
	cmdline := []string{"gsutil", "-q", "cp",
		LastMigrationObjectName,
		"gs://" + filepath.Join(g.Config().Bucket, LastMigrationObjectName)}
	return g.gsutil(cmdline)
}

// Config implements StateProvider
func (g GCS) Config() Config {
	return g.onDiskAccess.Config()
}

func (g GCS) gsutil(cmdline []string) error {
	if g.onDiskAccess.Config().Verbose {
		log.Println(strings.Join(cmdline, " "))
	}
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return tre.New(err, "gsutil cp failed", "output:", string(stdoutStderr))
	}
	return nil
}
