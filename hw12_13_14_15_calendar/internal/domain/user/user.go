package user

import "database/sql/driver"

type UID string

func (uid UID) Value() (driver.Value, error) {
	return string(uid), nil
}

func (uid *UID) Scan(src interface{}) error {
	*uid = UID(src.(string))
	return nil
}
