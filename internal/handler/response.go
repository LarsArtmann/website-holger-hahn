package handler

import (
	"context"

	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/domain"
	"holger-hahn-website/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// ResponseHandler provides standardized response handling for HTTP handlers.
type ResponseHandler struct{}

// NewResponseHandler creates a new response handler.
func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

// HandleError sends a standardized error response.
func (r *ResponseHandler) HandleError(c *gin.Context, err error) {
	switch {
	case domain.IsValidationError(err):
		c.JSON(constants.HTTPBadRequest, gin.H{"error": err.Error()})
	case domain.IsNotFoundError(err):
		c.JSON(constants.HTTPNotFound, gin.H{"error": err.Error()})
	case domain.IsConflictError(err):
		c.JSON(constants.HTTPConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(constants.HTTPInternalServerError, gin.H{"error": "Internal server error"})
	}
}

// HandleSuccess sends a standardized success response with data.
func (r *ResponseHandler) HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(constants.HTTPOKStatus, data)
}

// HandleCreated sends a standardized created response with data.
func (r *ResponseHandler) HandleCreated(c *gin.Context, data interface{}) {
	c.JSON(constants.HTTPCreated, data)
}

// HandleNoContent sends a standardized no content response.
func (r *ResponseHandler) HandleNoContent(c *gin.Context) {
	c.Status(constants.HTTPNoContent)
}

// HandleTemplateError sends a standardized template error response.
func (r *ResponseHandler) HandleTemplateError(c *gin.Context, err error) {
	c.JSON(constants.HTTPInternalServerError, gin.H{"error": domain.ErrRenderTemplate.Error()})
}

// RenderTemplate renders a templ component and handles errors.
func (r *ResponseHandler) RenderTemplate(c *gin.Context, component templ.Component) {
	c.Header("Content-Type", "text/html")
	
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		r.HandleTemplateError(c, err)
		return
	}
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// HandleHealthCheck sends a standardized health check response.
func (r *ResponseHandler) HandleHealthCheck(c *gin.Context) {
	r.HandleSuccess(c, HealthResponse{
		Status:  "healthy",
		Version: constants.AppVersion,
	})
}

// ListResponse provides standardized list response handling.
type ListResponse[T any] struct {
	Data  []T `json:"data"`
	Count int `json:"count"`
}

// HandleListSuccess sends a standardized list response with metadata.
func (r *ResponseHandler) HandleListSuccess[T any](c *gin.Context, data []T) {
	response := ListResponse[T]{
		Data:  data,
		Count: len(data),
	}
	r.HandleSuccess(c, response)
}

// ServiceOperationContext provides context for service operations.
type ServiceOperationContext struct {
	Context context.Context
	Handler *ResponseHandler
}

// NewServiceOperationContext creates a new service operation context.
func NewServiceOperationContext() *ServiceOperationContext {
	return &ServiceOperationContext{
		Context: context.Background(),
		Handler: NewResponseHandler(),
	}
}

// ExecuteServiceOperation executes a service operation with standardized error handling.
func (s *ServiceOperationContext) ExecuteServiceOperation[T any](
	c *gin.Context,
	operation func(context.Context) (T, error),
	successHandler func(*gin.Context, T),
) {
	result, err := operation(s.Context)
	if err != nil {
		s.Handler.HandleError(c, err)
		return
	}
	
	successHandler(c, result)
}

// CommonHandlerPatterns provides common handler operation patterns.
type CommonHandlerPatterns struct {
	responseHandler *ResponseHandler
	serviceContext  *ServiceOperationContext
}

// NewCommonHandlerPatterns creates a new common handler patterns helper.
func NewCommonHandlerPatterns() *CommonHandlerPatterns {
	return &CommonHandlerPatterns{
		responseHandler: NewResponseHandler(),
		serviceContext:  NewServiceOperationContext(),
	}
}

// HandleListOperation executes a list operation with standardized response.
func (c *CommonHandlerPatterns) HandleListOperation[T any](
	ginCtx *gin.Context,
	operation func(context.Context) ([]T, error),
) {
	c.serviceContext.ExecuteServiceOperation(
		ginCtx,
		operation,
		func(ctx *gin.Context, data []T) {
			c.responseHandler.HandleListSuccess(ctx, data)
		},
	)
}

// HandleSingleOperation executes a single item operation with standardized response.
func (c *CommonHandlerPatterns) HandleSingleOperation[T any](
	ginCtx *gin.Context,
	operation func(context.Context) (T, error),
) {
	c.serviceContext.ExecuteServiceOperation(
		ginCtx,
		operation,
		func(ctx *gin.Context, data T) {
			c.responseHandler.HandleSuccess(ctx, data)
		},
	)
}

// HandleCreateOperation executes a create operation with standardized response.
func (c *CommonHandlerPatterns) HandleCreateOperation[T any](
	ginCtx *gin.Context,
	operation func(context.Context) (T, error),
) {
	c.serviceContext.ExecuteServiceOperation(
		ginCtx,
		operation,
		func(ctx *gin.Context, data T) {
			c.responseHandler.HandleCreated(ctx, data)
		},
	)
}

// HandleDeleteOperation executes a delete operation with standardized response.
func (c *CommonHandlerPatterns) HandleDeleteOperation(
	ginCtx *gin.Context,
	operation func(context.Context) error,
) {
	_, err := operation(c.serviceContext.Context)
	if err != nil {
		c.responseHandler.HandleError(ginCtx, err)
		return
	}
	
	c.responseHandler.HandleNoContent(ginCtx)
}