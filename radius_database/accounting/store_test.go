package accounting

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.Exec("DROP TABLE radacct;")
	db.Exec(`CREATE TABLE radacct (
	RadAcctId               bigserial PRIMARY KEY,
		AcctSessionId           text NOT NULL,
		AcctUniqueId            text NOT NULL UNIQUE,
		UserName                text,
		Realm                   text,
		NASIPAddress            inet NOT NULL,
		NASPortId               text,
		NASPortType             text,
		AcctStartTime           timestamp with time zone,
		AcctUpdateTime          timestamp with time zone,
		AcctStopTime            timestamp with time zone,
		AcctInterval            bigint,
		acctsessiontime         bigint,
		AcctAuthentic           text,
		ConnectInfo_start       text,
		ConnectInfo_stop        text,
		AcctInputOctets         bigint,
		AcctOutputOctets        bigint,
		CalledStationId         text,
		CallingStationId        text,
		AcctTerminateCause      text,
		ServiceType             text,
		FramedProtocol          text,
		FramedIPAddress         inet
	);`)
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewStore(db)

		if assert.Equal(t, testStore, store) {
			if assert.NoError(t, db.Exec(`
				INSERT INTO radacct 
				(acctsessionid, acctuniqueid, nasipaddress, username, acctsessiontime) 
				VALUES('id', 'uid', '0.0.0.0', 'CssMHdz5VPVsKj62', 120);`).Error,
			) && assert.NoError(t, db.Exec(`
				INSERT INTO radacct 
				(acctsessionid, acctuniqueid, nasipaddress, username, acctsessiontime) 
				VALUES('id1', 'uid1', '0.0.0.0', 'CssMHdz5VPVsKj62', 20);`).Error,
			) && assert.NoError(t, db.Exec(`
				INSERT INTO radacct 
				(acctsessionid, acctuniqueid, nasipaddress, username, acctsessiontime) 
				VALUES('id2', 'uid2', '0.0.0.0', 'CssMHdz5VPVsKa49', 20);`).Error,
			) {
				totalTime, err := store.SessionTimeSum("CssMHdz5VPVsKj62")
				if assert.NoError(t, err) {
					assert.Equal(t, int64(140), totalTime)
				}

				assert.NoError(t, store.DeleteAccountingByUsername("CssMHdz5VPVsKj62"))

				totalTime, err = store.SessionTimeSum("CssMHdz5VPVsKj62")
				if assert.NoError(t, err) {
					assert.Equal(t, int64(0), totalTime)
				}

				totalTime, err = store.SessionTimeSum("CssMHdz5VPVsKa49")
				if assert.NoError(t, err) {
					assert.Equal(t, int64(20), totalTime)
				}
			}
		}

	}
}
