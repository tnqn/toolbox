package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type options struct {
	baseImage       string
	sourceFile      string
	destinationFile string
}

func newOption() *options {
	return &options{
		baseImage: "alpine",
	}
}

func (o *options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.baseImage, "base-image", "b", o.baseImage, "Base Docker image")
	fs.StringVarP(&o.sourceFile, "source-file", "s", o.sourceFile, "Source file to be added to the image")
	fs.StringVarP(&o.destinationFile, "destination-file", "d", o.destinationFile, "Destination file path inside the image")
}

func newCommand() *cobra.Command {
	o := newOption()
	cmd := &cobra.Command{
		Use:     "image-builder IMAGE:TAG [-b BASEIMAGE] -s SOURCE -d DESTINATION",
		Example: "image-builder myimage:v1 -b ubuntu:latest -s ./app -d /app",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "No tags is specified, the built image will not be tagged")
			}
			return run(o, args)
		},
	}
	flags := cmd.Flags()
	o.AddFlags(flags)
	return cmd
}

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running image-builder: %v\n", err)
		os.Exit(1)
	}
}

func run(o *options, tags []string) error {
	if o.sourceFile == "" || o.destinationFile == "" {
		return fmt.Errorf("source and destination file are required")
	}

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	// Build Docker context (tar archive)
	dockerContext, err := createDockerContext(o.sourceFile, o.destinationFile, o.baseImage)
	if err != nil {
		return fmt.Errorf("error creating Docker context: %w", err)
	}

	// Build the Docker image
	buildResponse, err := cli.ImageBuild(context.Background(), dockerContext, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		return fmt.Errorf("error building Docker image: %w", err)
	}
	defer buildResponse.Body.Close()

	// Output build logs
	_, err = io.Copy(os.Stderr, buildResponse.Body)
	if err != nil {
		return fmt.Errorf("error streaming build logs: %w", err)
	}
	return nil
}

// createDockerContext creates a tar archive with a Dockerfile and the source file.
func createDockerContext(sourceFile, destinationFile, baseImage string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Add Dockerfile
	dockerfileContent := fmt.Sprintf(`FROM %s
COPY %s %s
`, baseImage, filepath.Base(sourceFile), destinationFile)

	addFileToTar(tw, "Dockerfile", dockerfileContent, nil)

	// Add source file
	sourceFileContent, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %v", err)
	}
	sourceFileInfo, err := os.Stat(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %v", err)
	}
	addFileToTar(tw, filepath.Base(sourceFile), string(sourceFileContent), sourceFileInfo)

	return buf, nil
}

// addFileToTar adds a file to the tar archive.
func addFileToTar(tw *tar.Writer, name, content string, fileInfo os.FileInfo) error {
	header := &tar.Header{
		Name: name,
		Size: int64(len(content)),
	}
	if fileInfo != nil {
		header.Mode = int64(fileInfo.Mode().Perm())
		header.ModTime = fileInfo.ModTime()
	}
	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file content to tar: %v", err)
	}
	return nil
}
