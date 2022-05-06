package DataBase

import "testing"

const stmt_ = `create table TABLE_USER ("id"	INTEGER NOT NULL,
				"uid"	TEXT NOT NULL UNIQUE,
				"count"	INTEGER NOT NULL,
				"register_time" datetime NOT NULL,
				"month_time" TEXT NOT NULL DEFAULT '',
				"is_forbidden"	INTEGER DEFAULT 0,
				"comment" TEXT DEFAULT '',
				PRIMARY KEY("id"));`

func TestImportDB(t *testing.T) {




}