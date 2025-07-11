package sid

import (
	"fmt"
	"sync"
)

var g = &Generator{
	kindID: 0,
	nodeID: 0,
}

type Generator struct {
	sync.Map
	kindID int64
	nodeID int64
}

func (g *Generator) next(name string) (id ID, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	gen, _ := g.LoadOrStore(name, g.mustNewSnowflake())

	id = gen.(*Snowflake).Next()
	return
}

func (g *Generator) mustNewSnowflake() *Snowflake {
	gen, err := NewSnowflake(g.kindID, g.nodeID)
	if err != nil {
		panic(err)
	}
	return gen
}

func Set(kindID, nodeID int64) error {
	_, err := NewSnowflake(kindID, nodeID)
	if err != nil {
		return err
	}
	g.kindID = kindID
	g.nodeID = nodeID
	return nil
}

func NextID(name string) (int64, error) {
	id, err := g.next(name)
	if err != nil {
		return 0, err
	}
	return id.Int64(), nil
}
