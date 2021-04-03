package services

import (
	"database/sql"
	"github.com/leonj1/compass/exceptions"
	"github.com/leonj1/compass/models"
)

type Compass struct {
	db *sql.DB
}

func NewCompass(db *sql.DB) *Compass {
	return &Compass{db: db}
}

func (a *Compass) SetApplicationEnv(name, env, value string) (*models.Application, error) {
	t := models.Application{}
	apps, err := t.FindByApplicationName(a.db, name)
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		t.Name = name
		t.Envs = map[string]string{
			env: value,
		}
		return t.Save(a.db, t)
	}
	existingApp := new(models.Application)
	for _, app := range apps {
		if app.Name == name {
			existingApp = &app
			break
		}
	}
	if existingApp == nil {
		return nil, exceptions.NewNotFound(name)
	}
	existingApp.Envs[env] = value
	return t.Save(a.db, *existingApp)
}

func (a *Compass) FetchAll() ([]models.Application, error) {
	t := models.Application{}
	return t.FindAll(a.db)
}

func (a *Compass) FetchApplicationByName(name string) ([]models.Application, error) {
	t := models.Application{}
	return t.FindByApplicationName(a.db, name)
}

func (a *Compass) FetchApplicationByNameAndEnv(name, env string) (*string, error) {
	apps, err := a.FetchApplicationByName(name)
	if err != nil {
		return nil, err
	}
	if len(apps) == 0 {
		return nil, exceptions.NewNotFound(name)
	} else if len(apps) > 1 {
		return nil, exceptions.NewConflict("App %s expected 1 returned %d", name, len(apps))
	}
	app := apps[0]
	for k, v := range app.Envs {
		if k == env {
			return &v, nil
		}
	}
	return nil, exceptions.NewNotFound("Nothing found for App %s and env %s", name, env)
}
