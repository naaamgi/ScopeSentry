package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
)

// mergeTaskTargets preserves project scope while adding targets from selected tasks.
func mergeTaskTargets(existing string, tasks []models.Task) string {
	seen := make(map[string]struct{})
	lines := make([]string, 0)
	for _, source := range append([]string{existing}, taskTargets(tasks)...) {
		for _, line := range strings.Split(source, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func taskTargets(tasks []models.Task) []string {
	targets := make([]string, 0, len(tasks))
	for _, task := range tasks {
		targets = append(targets, task.Target)
	}
	return targets
}

func taskNames(tasks []models.Task) []string {
	names := make([]string, 0, len(tasks))
	for _, task := range tasks {
		names = append(names, task.Name)
	}
	return names
}

func rootDomainsForTargets(target, ignore string) ([]string, error) {
	targets, err := helper.GetTargetList(target, ignore)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, target := range targets {
		root := rootDomainForTarget(target)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		result = append(result, root)
	}
	return result, nil
}

func (s *service) SyncTaskTargets(ctx context.Context, projectID, tag, name string, tasks []models.Task) (string, error) {
	if len(tasks) == 0 {
		return "", fmt.Errorf("select at least one scan task")
	}
	if projectID == "" {
		name, tag = strings.TrimSpace(name), strings.TrimSpace(tag)
		if name == "" || tag == "" {
			return "", fmt.Errorf("project name and tag are required")
		}
		exists, err := s.projectRepo.ExistsByName(ctx, name)
		if err != nil {
			return "", err
		}
		if exists {
			return "", fmt.Errorf("project name already exists")
		}
		target := mergeTaskTargets("", tasks)
		if target == "" {
			return "", fmt.Errorf("selected tasks have no targets")
		}
		roots, err := rootDomainsForTargets(target, "")
		if err != nil {
			return "", err
		}
		first := tasks[0]
		newProject := &models.Project{
			Name: name, Tag: tag, Tp: "project", RootDomains: roots,
			Template: first.Template, Node: first.Node, AllNode: first.AllNode,
			Duplicates: first.Duplicates, Hour: 24,
		}
		projectID, err = s.projectRepo.InsertProject(ctx, newProject)
		if err != nil {
			return "", err
		}
		if err := s.projectRepo.UpsertProjectTarget(ctx, projectID, target); err != nil {
			_ = s.projectRepo.DeleteProjects(ctx, []string{projectID})
			return "", err
		}
		if err := s.projectRepo.UpdateAssetsForTasks(ctx, taskNames(tasks), projectID); err != nil {
			return projectID, fmt.Errorf("project created but asset linking failed: %w", err)
		}
		return projectID, nil
	}

	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return "", err
	}
	if project == nil {
		return "", fmt.Errorf("project not found")
	}
	existing, err := s.projectRepo.GetTarget(ctx, projectID)
	if err != nil {
		return "", err
	}
	target := mergeTaskTargets(existing, tasks)
	roots, err := rootDomainsForTargets(target, project.Ignore)
	if err != nil {
		return "", err
	}
	if err := s.projectRepo.UpsertProjectTarget(ctx, projectID, target); err != nil {
		return "", err
	}
	if err := s.projectRepo.UpdateRootDomains(ctx, projectID, roots); err != nil {
		return "", err
	}
	if err := s.projectRepo.UpdateAssetsForTasks(ctx, taskNames(tasks), projectID); err != nil {
		return projectID, fmt.Errorf("targets saved but asset linking failed: %w", err)
	}
	return projectID, nil
}
