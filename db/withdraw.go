// Package db provides ...
package db

const colWithdraw = "withdraw"

type Withdraw struct {
	UserID   string  `bson:"user_id"`
	UserName string  `bson:"user_name"`
	Amount   float64 `bson:"amount"`
	Address  string  `bson:"address"`
	TxID     string  `bson:"txid"`
}

func (db *DB) SaveWithdraw(userID, userName, address, txid string, amount float64) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(db.database).C(colWithdraw)
	data := &Withdraw{
		UserID:   userID,
		UserName: userName,
		Amount:   amount,
		Address:  address,
		TxID:     txid,
	}
	_, err := col.Upsert(data, data)
	if err != nil {
		log.Errorf("SaveWithdraw Error:%s [%s][%s][%s][%.8f]", err, userID, address, txid, amount)
	}
}
