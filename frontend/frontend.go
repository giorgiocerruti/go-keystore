package frontend

import (
	"fmt"

	v1 "github.com/giorgiocerruti/go-keystore/frontend/rest/v1"
)

func NewFrontEnd(frontend string) (Frontend, error) {
	switch frontend {
	case "rest":
		return v1.NewRestFrontend(), nil
	case "rpc":
		//@giorgiocerruti To be implemented
		return nil, nil
	case "":
		return nil, fmt.Errorf("frontend type is needed")
	default:
		return nil, fmt.Errorf("no such frontend")
	}
}
