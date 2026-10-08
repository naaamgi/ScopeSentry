package project

import (
	"context"
	"reflect"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	repo "github.com/Autumn-27/ScopeSentry/internal/repositories/project"
)

type syncRepositoryStub struct {
	repo.Repository
	project       *models.Project
	target        string
	linkedTasks   []string
	linkedProject string
	roots         []string
}

func (r *syncRepositoryStub) FindByID(context.Context, string) (*models.Project, error) {
	return r.project, nil
}

func (r *syncRepositoryStub) GetTarget(context.Context, string) (string, error) {
	return r.target, nil
}

func (r *syncRepositoryStub) UpsertProjectTarget(_ context.Context, _ string, target string) error {
	r.target = target
	return nil
}

func (r *syncRepositoryStub) UpdateRootDomains(_ context.Context, _ string, roots []string) error {
	r.roots = roots
	return nil
}

func (r *syncRepositoryStub) UpdateAssetsForTasks(_ context.Context, names []string, projectID string) error {
	r.linkedTasks = names
	r.linkedProject = projectID
	return nil
}

func TestMergeTaskTargetsPreservesExistingScopeAndDeduplicates(t *testing.T) {
	tasks := []models.Task{
		{Target: "api.example.com\nexample.com"},
		{Target: "api.example.com\n192.0.2.10"},
	}
	got := mergeTaskTargets("example.com\nlegacy.example.org", tasks)
	want := "example.com\nlegacy.example.org\napi.example.com\n192.0.2.10"
	if got != want {
		t.Fatalf("merged targets = %q, want %q", got, want)
	}
}

func TestRootDomainsForTargets(t *testing.T) {
	got, err := rootDomainsForTargets("example.com\napi.example.com\n192.0.2.10", "")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"example.com", "192.0.2.10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("root domains = %v, want %v", got, want)
	}
}

func TestSyncTaskTargetsKeepsExistingTargetsAndLinksSelectedTask(t *testing.T) {
	r := &syncRepositoryStub{
		project: &models.Project{Name: "test"},
		target:  "legacy.example.org",
	}
	s := &service{projectRepo: r}
	id, err := s.SyncTaskTargets(context.Background(), "project-id", "", "", []models.Task{
		{Name: "scan-1", Target: "api.example.com"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "project-id" || r.target != "legacy.example.org\napi.example.com" {
		t.Fatalf("project ID = %q, target = %q", id, r.target)
	}
	if !reflect.DeepEqual(r.linkedTasks, []string{"scan-1"}) || r.linkedProject != id {
		t.Fatalf("linked tasks = %v, project = %q", r.linkedTasks, r.linkedProject)
	}
	if !reflect.DeepEqual(r.roots, []string{"example.org", "example.com"}) {
		t.Fatalf("root domains = %v", r.roots)
	}
}
