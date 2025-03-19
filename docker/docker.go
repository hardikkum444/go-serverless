package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func BuildDockerImage(dir, language, handlerfile string) (string, error) {
	var dockerfilecontent string

	template, err := LoadTemplate(language)
	if err != nil {
		return "", fmt.Errorf("failed to load the template: %w", err)
	}

	if language == "python" {
		dockerfilecontent = fmt.Sprintf(template.Dockerfile, handlerfile)
	}

	if language == "golang" {
		dockerfilecontent = fmt.Sprintf(template.Dockerfile)
	}

	// writing the dockerfile to the directory
	dockerfilepath := filepath.Join(dir, "DockerFile")
	if err := os.WriteFile(dockerfilepath, []byte(dockerfilecontent), 0644); err != nil {
		return "", fmt.Errorf("failed to write Dockerfile: %v", err)
	}

	// building the docker image
	imageTag := fmt.Sprintf("go-serverless:%s", language)
	cmd := exec.Command("docker", "build", "-t", imageTag, dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker build failed: %s", output)
	}

	// extracting the image id from the build
	imageID, err := ExtractImageID(string(output))
	if err != nil {
		return "", fmt.Errorf("failed to extract the image ID: %v", err)
	}

	return imageID, nil
}

type Template struct {
	Dockerfile string `yaml:"dockerfile"`
}

func LoadTemplate(language string) (*Template, error) {
	// determining the path to the template file
	templatefile := fmt.Sprintf("templates/%s.yaml", language)

	// reading the template file
	data, err := os.ReadFile(templatefile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %v", err)
	}

	var template Template

	// parsing the YAML content
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal template %v", err)
	}

	return &template, nil
}
