// Package resutil provides RES utility functions, complementing the more common ones in the github.com/jirenius/go-res
// package.
package resutil

import (
	"strings"

	"github.com/Autherain/saltyChat/internal/utils/errors"
	"github.com/Autherain/saltyChat/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

func HandleCollectionQueryRequest[Collection ~[]Model, Model any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) (Collection, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			errors.LogAndWriteRESError(logger.Default(), request, err)
			return
		}

		request.Collection(response)
	})
}

func HandleModelQueryRequest[Model any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) (Model, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			errors.LogAndWriteRESError(logger.Default(), request, err)
			return
		}

		request.Model(response)
	})
}

// https://resgate.io/docs/specification/res-service-protocol/#query-request
func handleQueryRequest(
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest),
) error {
	return service.With(rid, func(resource res.Resource) {
		resource.QueryEvent(func(request res.QueryRequest) {
			if request == nil {
				return // https://github.com/jirenius/go-res/blob/372a82d603a13d7601f8b14e74eccaebd325ee61/resource.go#L336-L339
			}

			handler(request)
		})
	})
}

// ParseUUIDPathParam from the resource with the given key.
func ParseUUIDPathParam(resource res.Resource, key string) (uuid.UUID, error) {
	result, err := uuid.Parse(resource.PathParam(key))
	if err != nil {
		return uuid.Nil, &errors.Error{
			Code:            errors.CodeInvalid,
			Message:         "Invalid '" + key + "' path parameter",
			UnderlyingError: err,
		}
	}

	return result, nil
}

func JoinResourcePath(parts ...string) string {
	return strings.Join(parts, ".")
}
