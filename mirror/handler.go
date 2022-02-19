package mirror

import "fmt"

func (m *mirror) handleHostMsg(msg []byte) {
	fmt.Println("host message: ", string(msg))
}
