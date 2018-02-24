package connpool

import ("fmt")

const numConnections = 10

type Connection interface {
	Close() error
	Execute() error
}

type ConnectionPool interface {
	getConnection() (Connection, error)
}

type MyConnection struct {

}

func (m *MyConnection) Close() error {
	return nil
}

func (m *MyConnection) Execute() error {
	return nil
}

type MyConnectionPool struct {
	connections []MyConnection //Perhaps should be Connection
	used []bool
}

func (p *MyConnectionPool) GetConnection() (Connection, error){

	for i, c := range p.connections {
		if p.used[i] == false { 
			p.used[i] = true
			return &c, nil
		}
		//fmt.Printf("%d %t\n", i, p.used[i])
	}
	
	return nil, fmt.Errorf("No connections available")
}

func New() MyConnectionPool {
	var p MyConnectionPool
	p.connections = make([]MyConnection, numConnections)
	p.used = make([]bool, numConnections)
	return p
}
