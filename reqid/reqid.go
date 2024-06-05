package reqid

import (
	"github.com/google/uuid"
)

func GenRequestID() string {
	oid, _ := uuid.NewUUID()
	return oid.String()
}
