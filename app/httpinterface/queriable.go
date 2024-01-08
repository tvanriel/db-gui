package httpinterface

import "github.com/tvanriel/db-gui/app/domain"

func getQueriable(conn domain.Connection, databaseName string) (domain.Queriable, error) {
	if databaseName == "" {
		err := conn.Connect()
		if err != nil {
			return nil, err
		}

		return conn.AsQueriable(), nil
	}
	db := conn.GetDatabase(databaseName)
	err := db.Connect()

	if err != nil {
		return nil, err
	}

	return db.AsQueriable(), nil
}
