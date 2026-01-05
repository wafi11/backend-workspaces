package products

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Docker Operations
func (s *Service) StartContainer(c context.Context, containerID int) error {
	container, err := s.repo.GetContainerByID(c, containerID)
	if err != nil {
		return err
	}

	if container.Status == "running" {
		return fmt.Errorf("container already running")
	}

	// Build docker run command
	args := []string{"run", "-d", "--name", container.ContainerName}

	if container.HostPort > 0 && container.ContainerPort > 0 {
		args = append(args, "-p", fmt.Sprintf("%d:%d", container.HostPort, container.ContainerPort))
	}

	for key, value := range container.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	for _, volume := range container.Volumes {
		args = append(args, "-v", volume)
	}

	imageTag := fmt.Sprintf("%s:%s", container.Image, container.Tag)
	args = append(args, imageTag)

	if len(container.Command) > 0 {
		args = append(args, container.Command...)
	}

	cmd := exec.CommandContext(c, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.repo.UpdateContainerStatus(c, containerID, "error", "")
		return fmt.Errorf("failed to start container: %w, output: %s", err, string(output))
	}

	dockerID := strings.TrimSpace(string(output))
	return s.repo.UpdateContainerStatus(c, containerID, "running", dockerID)
}

func (s *Service) StopContainer(c context.Context, containerID int) error {
	container, err := s.repo.GetContainerByID(c, containerID)
	if err != nil {
		return err
	}

	if container.Status != "running" {
		return fmt.Errorf("container is not running")
	}

	cmd := exec.CommandContext(c, "docker", "stop", container.ContainerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return s.repo.UpdateContainerStatus(c, containerID, "stopped", container.DockerContainerID)
}

func (s *Service) RestartContainer(c context.Context, containerID int) error {
	if err := s.StopContainer(c, containerID); err != nil {
		return err
	}

	container, _ := s.repo.GetContainerByID(c, containerID)

	// Remove old container
	exec.CommandContext(c, "docker", "rm", container.ContainerName).Run()

	return s.StartContainer(c, containerID)
}

func (s *Service) GetContainerLogs(c context.Context, containerID int, tail int) (string, error) {
	container, err := s.repo.GetContainerByID(c, containerID)
	if err != nil {
		return "", err
	}

	args := []string{"logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	args = append(args, container.ContainerName)

	cmd := exec.CommandContext(c, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}

	return string(output), nil
}

func (s *Service) GetContainerStatus(c context.Context, containerID int) (map[string]interface{}, error) {
	container, err := s.repo.GetContainerByID(c, containerID)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(c, "docker", "ps", "-a", "--filter",
		fmt.Sprintf("name=%s", container.ContainerName),
		"--format", "{{.ID}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return map[string]interface{}{
			"container_name": container.ContainerName,
			"status":         "not found",
			"running":        false,
		}, nil
	}

	parts := strings.Split(result, "|")
	status := map[string]interface{}{
		"container_id":   parts[0],
		"container_name": container.ContainerName,
		"status":         parts[1],
		"running":        strings.Contains(parts[1], "Up"),
	}

	if len(parts) > 2 {
		status["ports"] = parts[2]
	}

	dbStatus := "stopped"
	if strings.Contains(parts[1], "Up") {
		dbStatus = "running"
	}
	s.repo.UpdateContainerStatus(c, containerID, dbStatus, parts[0])

	return status, nil
}
