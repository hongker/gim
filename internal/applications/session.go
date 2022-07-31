package applications

import "gim/internal/domain/entity"

type SessionApp struct {

}

func (app *SessionApp) List() (items []entity.Session, err error) {
	return
}

func (app *SessionApp) Delete(id string) error {
	return nil
}
