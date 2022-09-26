package bucket

import "github.com/ebar-go/ego/server/ws"

type Bucket interface {
}

func FindUserConnection(uid string) (ws.Conn, error) {
	return nil, nil
}

func SaveUserConnection(uid string, conn ws.Conn) {

}
