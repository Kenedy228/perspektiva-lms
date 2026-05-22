package postgres

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
)

func nullUUID(id uuid.UUID) any {
	if id == uuid.Nil {
		return nil
	}
	return id
}

func uuidArray(ids []uuid.UUID) string {
	values := make([]string, 0, len(ids))
	for _, id := range ids {
		values = append(values, id.String())
	}
	return "{" + strings.Join(values, ",") + "}"
}

func nullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}
