package task

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/gin-gonic/gin"
)

type syncToProjectRequest struct {
	IDs     []string `json:"ids" binding:"required"`
	Option  string   `json:"option" binding:"required"`
	Project string   `json:"project"`
	Tag     string   `json:"tag"`
	Name    string   `json:"name"`
}

func SyncToProject(c *gin.Context) {
	var req syncToProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 || (req.Option != "existing" && req.Option != "new") ||
		(req.Option == "existing" && strings.TrimSpace(req.Project) == "") ||
		(req.Option == "new" && (strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Tag) == "")) {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("select tasks and a project, or enter a new project name and tag"))
		return
	}
	projectID, err := taskService.SyncToProject(c, req.IDs, req.Option, req.Project, req.Tag, req.Name)
	if err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	response.Success(c, gin.H{"projectId": projectID}, "api.success")
}
