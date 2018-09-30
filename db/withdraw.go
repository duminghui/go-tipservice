// Package db provides ...
package db

func (db *DB) SaveWithdraw(userID, address, txid string, amount float64) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(db.database).C(colWithdraw)
	data := &Withdraw{
		UserID:  userID,
		Amount:  amount,
		Address: address,
		TxID:    txid,
	}
	_, err := col.Upsert(data, data)
	if err != nil {
		log.Errorf("SaveWithdraw Error:%s [%s][%s][%s][%.8f]", err, userID, address, txid, amount)
	}
}
