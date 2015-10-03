package srm

import "fmt"

func deleteall(db *DbMap, table string) error {
	_, err := db.Conn.Exec(fmt.Sprintf("DELETE FROM %s", table))
	if err != nil {
		return err
	}
	return nil
}
