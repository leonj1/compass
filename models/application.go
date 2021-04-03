package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Application struct {
	Id   int64             `json:"-"`
	Name string            `json:"name"`
	Envs map[string]string `json:"envs"`
}

const APPLICATION_TABLE = "applications"
const ENV_VALUES_TABLE = "env_values"

func (a *Application) Save(db *sql.DB, item Application) (*Application, error) {
	if db == nil {
		return nil, errors.New("db cannot be empty")
	}
	var sqlCmd string
	if item.Id == 0 {
		sqlCmd = fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", APPLICATION_TABLE)
	} else {
		sqlCmd = fmt.Sprintf("UPDATE %s SET name=? WHERE id=%d", APPLICATION_TABLE, item.Id)
	}
	res, err := db.Exec(sqlCmd, item.Name)
	if err != nil {
		return nil, err
	}
	if item.Id == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		item.Id = id
	}
	err = a.deleteEnvsByApplicationId(db, item.Id)
	if err != nil {
		return nil, err
	}
	err = a.setEnvsByApplicationId(db, item.Id, item.Envs)
	if err != nil {
		return nil, err
	}
	return	&item, nil
}

func (a *Application) FindAll(db *sql.DB) ([]Application, error) {
	if db == nil {
		return nil, errors.New("db cannot be empty")
	}
	sqlCmd := fmt.Sprintf("SELECT id, name FROM %s", APPLICATION_TABLE)
	rows, err := db.Query(sqlCmd)
	if err != nil {
		return nil, err
	}
	results := []Application{}
	for rows.Next() {
		t := new(Application)
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			return nil, err
		}
		envs, err := a.findByApplicationId(db, t.Id)
		if err != nil {
			return nil, err
		}
		t.Envs = envs
		results = append(results, *t)
	}
	return results, nil
}

func (a *Application) FindByApplicationName(db *sql.DB, name string) ([]Application, error) {
	if db == nil {
		return nil, errors.New("db cannot be empty")
	} else if name == "" {
		return nil, errors.New("application name cannot be empty")
	}
	sqlCmd := fmt.Sprintf("SELECT id, name FROM %s WHERE name=?", APPLICATION_TABLE)
	rows, err := db.Query(sqlCmd, name)
	if err != nil {
		return nil, err
	}
	results := []Application{}
	for rows.Next() {
		t := new(Application)
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			return nil, err
		}
		envs, err := a.findByApplicationId(db, t.Id)
		if err != nil {
			return nil, err
		}
		t.Envs = envs
		results = append(results, *t)
	}
	return results, nil
}

// privates

func (a *Application) findByApplicationId(db *sql.DB, applicationId int64) (map[string]string, error) {
	if db == nil {
		return nil, errors.New("db cannot be empty")
	} else if applicationId == 0 {
		return nil, errors.New("application id cannot be empty")
	}
	sqlCmd := fmt.Sprintf("SELECT key, value FROM %s WHERE application_id=?", ENV_VALUES_TABLE)
	rows, err := db.Query(sqlCmd, applicationId)
	if err != nil {
		return nil, err
	}
	results := map[string]string{}
	for rows.Next() {
		var key string
		var value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		results[key] = value
	}
	return results, nil
}

func (a *Application) deleteEnvsByApplicationId(db *sql.DB, applicationId int64) error {
	if db == nil {
		return errors.New("db cannot be empty")
	}
	if applicationId == 0 {
		return errors.New("application id cannot be empty")
	}
	sqlCmd := fmt.Sprintf("DELETE FROM %s WHERE application_id=?", ENV_VALUES_TABLE)
	_, err := db.Exec(sqlCmd, applicationId)
	if err != nil {
		return err
	}
	return nil
}

func (a *Application) setEnvsByApplicationId(db *sql.DB, applicationId int64, envs map[string]string) error {
	if db == nil {
		return errors.New("db cannot be empty")
	} else if applicationId == 0 {
		return errors.New("application id cannot be empty")
	}
	if len(envs) == 0 {
		return nil
	}
	sqlCmd := fmt.Sprintf("INSERT INTO %s ('application_id', 'key', 'value') VALUES (?,?,?)", ENV_VALUES_TABLE)
	for k, v := range envs {
		_, err := db.Exec(sqlCmd, applicationId, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
