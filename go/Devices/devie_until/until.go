package devie_until

import "database/sql"

func NullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		v := n.Int64
		return &v
	}
	return nil
}
