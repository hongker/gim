package applications

import "gim/internal/domain/dto"

type MessageApp struct {
	
}


func (app *MessageApp) Send(req *dto.MessageSendRequest) error {
	return nil
}



func (app *MessageApp) Query(req *dto.MessageQueryRequest) (*dto.MessageQueryResponse, error) {
	return nil, nil
}