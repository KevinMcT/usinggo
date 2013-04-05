package SlotList

import (
	//"lab5Merge/model/net/msg"
	//"fmt"
	"sync"
)

/*
The slot list. 
*/
type Slots struct {
	List       []interface{}
	SizeOfList int
	mutex      sync.RWMutex
}

func NewSlots() *Slots {
	return &Slots{SizeOfList: 0, List: make([]interface{}, 10)}
}

func (s *Slots) Add(message interface{}, slotNumber int) bool {
	//Check if a slot is taken
	if s.List[slotNumber] == nil {
		s.List[slotNumber] = message
		s.SizeOfList = s.SizeOfList + 1

		for i := s.SizeOfList; i < s.SizeOfList*2; i++ {
			s.List = append(s.List, nil)
		}
		//fmt.Println(s.List)
		return true
	}
	return false
}

func (s *Slots) Next() {

}
