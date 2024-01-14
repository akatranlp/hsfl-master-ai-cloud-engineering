package shared_types

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type Code int

const (
	OK Code = iota
	InvalidArgument
	NotFound
	Internal
	Unauthenticated
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case InvalidArgument:
		return "InvalidArgument"
	case NotFound:
		return "NotFound"
	case Internal:
		return "Internal"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return "Unknown"
	}
}

func (c Code) ToHTTPStatusCode() int {
	switch c {
	case OK:
		return http.StatusOK
	case InvalidArgument:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Internal:
		return http.StatusInternalServerError
	case Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (c Code) ToGRPCStatusCode() codes.Code {
	switch c {
	case OK:
		return codes.OK
	case InvalidArgument:
		return codes.InvalidArgument
	case NotFound:
		return codes.NotFound
	case Internal:
		return codes.Internal
	case Unauthenticated:
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
