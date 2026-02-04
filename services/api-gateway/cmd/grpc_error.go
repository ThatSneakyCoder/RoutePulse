package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isUnauthorized(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	return st.Code() == codes.Unauthenticated
}
