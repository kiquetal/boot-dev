package main

func (u user) doBattles(subCh <-chan move) []piece {
	battle := make([]piece,0)
	for m:= range subCh {
	    for _, mo:= range u.pieces {
			if (mo.location == m.piece.location) {
				battle = append(battle,mo)
			}
		}	
	}
	return battle
}

// don't touch below this line

type user struct {
	name   string
	pieces []piece
}

type move struct {
	userName string
	piece    piece
}

type piece struct {
	location string
	name     string
}

func (u user) march(p piece, publishCh chan<- move) {
	publishCh <- move{
		userName: u.name,
		piece:    p,
	}
}

func distributeBattles(publishCh <-chan move, subChans []chan move) {
	for mv := range publishCh {
		for _, subCh := range subChans {
			subCh <- mv
		}
	}
}

