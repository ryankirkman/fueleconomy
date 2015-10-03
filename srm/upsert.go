package srm

// TODO
// Write single query upsert using postgres function and test perfomance

func multiQueryUpsert(db *DbMap, table string, updateOnField string, ptr interface{}) (insertedId int, err error) {
	rowsAffected, err := db.UpdateOne(table, updateOnField, ptr)
	if err != nil {
		return insertedId, err
	}
	if rowsAffected == 0 {
		insertedId, err = db.InsertOne(table, ptr)
		if err != nil {
			return insertedId, err
		}
	}
	return insertedId, nil
}
