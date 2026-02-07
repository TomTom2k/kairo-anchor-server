package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
	projectUC "github.com/tomtom2k/kairo-anchor-server/internal/usecase/project"
)

type ProjectHandler struct {
	createProject *projectUC.CreateProjectUseCase
	updateProject *projectUC.UpdateProjectUseCase
	deleteProject *projectUC.DeleteProjectUseCase
	getProject    *projectUC.GetProjectUseCase
	listProjects  *projectUC.ListProjectsUseCase
}

func NewProjectHandler(
	create *projectUC.CreateProjectUseCase,
	update *projectUC.UpdateProjectUseCase,
	delete *projectUC.DeleteProjectUseCase,
	get *projectUC.GetProjectUseCase,
	list *projectUC.ListProjectsUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		createProject: create,
		updateProject: update,
		deleteProject: delete,
		getProject:    get,
		listProjects:  list,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProjectRequest true "Create Project Request"
// @Success 201 {object} APIResponse{data=ProjectResponse}
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	input := projectUC.CreateProjectInput{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Status:      project.ProjectStatus(req.Status),
		Progress:    req.Progress,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	p, err := h.createProject.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "CREATE_PROJECT_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusCreated, toProjectResponse(p), "Project created successfully")
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Param request body UpdateProjectRequest true "Update Project Request"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	projectID := c.Param("id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	input := projectUC.UpdateProjectInput{
		ID:          projectID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Status:      project.ProjectStatus(req.Status),
		Progress:    req.Progress,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	// Convert DTOs to domain models
	if req.Tasks != nil {
		tasks := make([]project.Task, len(req.Tasks))
		for i, t := range req.Tasks {
			tasks[i] = project.Task{
				ID:       t.ID,
				Title:    t.Title,
				Status:   project.TaskStatus(t.Status),
				Priority: project.TaskPriority(t.Priority),
				DueDate:  t.DueDate,
			}
		}
		input.Tasks = tasks
	}

	if req.Documents != nil {
		documents := make([]project.Document, len(req.Documents))
		for i, d := range req.Documents {
			documents[i] = project.Document{
				ID:        d.ID,
				Name:      d.Name,
				Type:      d.Type,
				Size:      d.Size,
				UpdatedAt: d.UpdatedAt,
			}
		}
		input.Documents = documents
	}

	p, err := h.updateProject.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "UPDATE_PROJECT_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Project updated successfully")
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	projectID := c.Param("id")
	err = h.deleteProject.Execute(c.Request.Context(), projectID, userID)
	if err != nil {
		SendError(c, http.StatusBadRequest, "DELETE_PROJECT_FAILED", err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, nil, "Project deleted successfully")
}

// GetProject godoc
// @Summary Get a project
// @Description Get a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	projectID := c.Param("id")
	p, err := h.getProject.Execute(c.Request.Context(), projectID, userID)
	if err != nil {
		SendError(c, http.StatusNotFound, ErrCodeNotFound, err.Error())
		return
	}

	SendSuccess(c, http.StatusOK, toProjectResponse(p), "")
}

// ListProjects godoc
// @Summary List all projects
// @Description Get all projects for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=[]ProjectResponse}
// @Failure 401 {object} APIErrorResponse
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}

	projects, err := h.listProjects.Execute(c.Request.Context(), userID)
	if err != nil {
		SendInternalError(c, err)
		return
	}

	response := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		response[i] = *toProjectResponse(&p)
	}

	SendSuccess(c, http.StatusOK, response, "")
}

// Helper function to convert domain model to response DTO
func toProjectResponse(p *project.Project) *ProjectResponse {
	tasks := make([]TaskDTO, len(p.Tasks))
	for i, t := range p.Tasks {
		tasks[i] = TaskDTO{
			ID:       t.ID,
			Title:    t.Title,
			Status:   string(t.Status),
			Priority: string(t.Priority),
			DueDate:  t.DueDate,
		}
	}

	documents := make([]DocumentDTO, len(p.Documents))
	for i, d := range p.Documents {
		documents[i] = DocumentDTO{
			ID:        d.ID,
			Name:      d.Name,
			Type:      d.Type,
			Size:      d.Size,
			UpdatedAt: d.UpdatedAt,
		}
	}

	return &ProjectResponse{
		ID:          p.ID.String(),
		UserID:      p.UserID.String(),
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		Progress:    p.Progress,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Tasks:       tasks,
		Documents:   documents,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
