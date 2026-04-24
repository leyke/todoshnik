package user

import "todoshnik/internal/bot/tg"

type UserState struct {
	tg.Command
	tg.State
}

type StateStorage struct {
	states map[int64]UserState
}

func NewStateStorage() *StateStorage {
	return &StateStorage{
		states: make(map[int64]UserState),
	}
}

func (ss *StateStorage) Set(userID int64, command tg.Command, state tg.State) {
	ss.states[userID] = UserState{command, state}
}

func (ss *StateStorage) Get(userID int64) (UserState, bool) {
	state, ok := ss.states[userID]
	return state, ok
}
