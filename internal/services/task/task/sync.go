package task

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *service) SyncToProject(ctx *gin.Context, ids []string, option, projectID, tag, name string) (string, error) {
	if len(ids) == 0 {
		return "", fmt.Errorf("select at least one scan task")
	}
	if option != "existing" && option != "new" {
		return "", fmt.Errorf("invalid project option")
	}
	if option == "existing" {
		if _, err := primitive.ObjectIDFromHex(projectID); err != nil {
			return "", fmt.Errorf("select a project, not a tag")
		}
		name, tag = "", ""
	} else {
		projectID = ""
		name, tag = strings.TrimSpace(name), strings.TrimSpace(tag)
		if name == "" || tag == "" {
			return "", fmt.Errorf("enter the new project name and tag")
		}
	}
	tasks := make([]models.Task, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", fmt.Errorf("invalid task ID: %s", id)
		}
		task, err := s.taskRepo.FindByID(ctx.Request.Context(), objectID)
		if err != nil {
			return "", err
		}
		if task == nil {
			return "", fmt.Errorf("task not found: %s", id)
		}
		tasks = append(tasks, *task)
	}
	return project.NewService().SyncTaskTargets(ctx.Request.Context(), projectID, tag, name, tasks)
}
