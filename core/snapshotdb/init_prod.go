// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

// +build !test

package snapshotdb

const (
	//DBPath path of db
	DBPath = "snapshotdb"
	//DBBasePath path of basedb
	DBBasePath  = "base"
	currentPath = "current"
)

// New  create a new snapshotDB,will clear old snapshotDB data
/*func New(path string) (DB, error) {
	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}
	s, err := openFile(path, false)
	if err != nil {
		logger.Error("open db file fail", "error", err, "path", dbpath)
		return nil, err
	}

	logger.Info("begin new", "path", path)
	db, err := newDB(s)
	if err != nil {
		logger.Error(fmt.Sprint("new db fail:", err))
		return nil, err
	}
	return db, nil
}
*/
