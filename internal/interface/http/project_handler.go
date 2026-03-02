package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
	projectUC "github.com/tomtom2k/kairo-anchor-server/internal/usecase/project"
)

type ProjectHandler struct {
	createProject   *projectUC.CreateProjectUseCase
	updateProject   *projectUC.UpdateProjectUseCase
	deleteProject   *projectUC.DeleteProjectUseCase
	getProject      *projectUC.GetProjectUseCase
	listProjects    *projectUC.ListProjectsUseCase
	addTask         *projectUC.AddTaskUseCase
	updateTask      *projectUC.UpdateTaskUseCase
	deleteTask      *projectUC.DeleteTaskUseCase
	reorderTasks    *projectUC.ReorderTasksUseCase
	addDocument     *projectUC.AddDocumentUseCase
	updateDocument  *projectUC.UpdateDocumentUseCase
	deleteDocument  *projectUC.DeleteDocumentUseCase
}

func NewProjectHandler(
	create *projectUC.CreateProjectUseCase,
	update *projectUC.UpdateProjectUseCase,
	delete *projectUC.DeleteProjectUseCase,
	get *projectUC.GetProjectUseCase,
	list *projectUC.ListProjectsUseCase,
	addTask *projectUC.AddTaskUseCase,
	updateTask *projectUC.UpdateTaskUseCase,
	deleteTask *projectUC.DeleteTaskUseCase,
	reorderTasks *projectUC.ReorderTasksUseCase,
	addDocument *projectUC.AddDocumentUseCase,
	updateDocument *projectUC.UpdateDocumentUseCase,
	deleteDocument *projectUC.DeleteDocumentUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		createProject:  create,
		updateProject:  update,
		deleteProject:  delete,
		getProject:     get,
		listProjects:   list,
		addTask:        addTask,
		updateTask:    updateTask,
		deleteTask:    deleteTask,
		reorderTasks:  reorderTasks,
		addDocument:   addDocument,
		updateDocument: updateDocument,
		deleteDocument: deleteDocument,
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
		Progress:    0, // tiến độ tính từ task hoàn thành
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
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
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

// AddTask godoc
// @Summary Add a task to a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body CreateTaskRequest true "Create Task Request"
// @Success 201 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/tasks [post]
func (h *ProjectHandler) AddTask(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	input := projectUC.AddTaskInput{
		ProjectID: projectID,
		UserID:    userID,
		Title:     req.Title,
		Status:    project.TaskStatus(req.Status),
		Priority:  project.TaskPriority(req.Priority),
		DueDate:   req.DueDate,
	}
	p, err := h.addTask.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "ADD_TASK_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusCreated, toProjectResponse(p), "Task added")
}

// UpdateTask godoc
// @Summary Update a task
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param taskId path string true "Task ID"
// @Param request body UpdateTaskRequest true "Update Task Request"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/tasks/{taskId} [put]
func (h *ProjectHandler) UpdateTask(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	taskID := c.Param("taskId")
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	input := projectUC.UpdateTaskInput{
		ProjectID: projectID,
		UserID:    userID,
		TaskID:    taskID,
		Title:     req.Title,
		Status:    ptrToTaskStatus(req.Status),
		Priority:  ptrToTaskPriority(req.Priority),
		DueDate:   req.DueDate,
	}
	p, err := h.updateTask.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "UPDATE_TASK_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Task updated")
}

// DeleteTask godoc
// @Summary Delete a task
// @Tags projects
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param taskId path string true "Task ID"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/tasks/{taskId} [delete]
func (h *ProjectHandler) DeleteTask(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	taskID := c.Param("taskId")
	p, err := h.deleteTask.Execute(c.Request.Context(), projectID, userID, taskID)
	if err != nil {
		SendError(c, http.StatusBadRequest, "DELETE_TASK_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Task deleted")
}

// ReorderTasks godoc
// @Summary Reorder tasks
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body ReorderTasksRequest true "Ordered list of task IDs"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/tasks/order [put]
func (h *ProjectHandler) ReorderTasks(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	var req ReorderTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	p, err := h.reorderTasks.Execute(c.Request.Context(), projectID, userID, req.TaskIDs)
	if err != nil {
		SendError(c, http.StatusBadRequest, "REORDER_TASKS_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Tasks reordered")
}

// AddDocument godoc
// @Summary Add a document to a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body CreateDocumentRequest true "Create Document Request"
// @Success 201 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/documents [post]
func (h *ProjectHandler) AddDocument(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	input := projectUC.AddDocumentInput{
		ProjectID: projectID,
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Size:      req.Size,
	}
	p, err := h.addDocument.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "ADD_DOCUMENT_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusCreated, toProjectResponse(p), "Document added")
}

// UpdateDocument godoc
// @Summary Update a document
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param docId path string true "Document ID"
// @Param request body UpdateDocumentRequest true "Update Document Request"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/documents/{docId} [put]
func (h *ProjectHandler) UpdateDocument(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	docID := c.Param("docId")
	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	input := projectUC.UpdateDocumentInput{
		ProjectID:  projectID,
		UserID:     userID,
		DocumentID: docID,
		Name:       req.Name,
		Type:       req.Type,
		Size:       req.Size,
	}
	p, err := h.updateDocument.Execute(c.Request.Context(), input)
	if err != nil {
		SendError(c, http.StatusBadRequest, "UPDATE_DOCUMENT_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Document updated")
}

// DeleteDocument godoc
// @Summary Delete a document
// @Tags projects
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param docId path string true "Document ID"
// @Success 200 {object} APIResponse{data=ProjectResponse}
// @Router /projects/{id}/documents/{docId} [delete]
func (h *ProjectHandler) DeleteDocument(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized")
		return
	}
	projectID := c.Param("id")
	docID := c.Param("docId")
	p, err := h.deleteDocument.Execute(c.Request.Context(), projectID, userID, docID)
	if err != nil {
		SendError(c, http.StatusBadRequest, "DELETE_DOCUMENT_FAILED", err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, toProjectResponse(p), "Document deleted")
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

	// Tiến độ luôn tính từ task hoàn thành
	progress := project.ProgressFromTasks(p.Tasks)

	return &ProjectResponse{
		ID:          p.ID.String(),
		UserID:      p.UserID.String(),
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		Progress:    progress,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Tasks:       tasks,
		Documents:   documents,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ptrToTaskStatus(s *string) *project.TaskStatus {
	if s == nil {
		return nil
	}
	v := project.TaskStatus(*s)
	return &v
}

func ptrToTaskPriority(s *string) *project.TaskPriority {
	if s == nil {
		return nil
	}
	v := project.TaskPriority(*s)
	return &v
}
