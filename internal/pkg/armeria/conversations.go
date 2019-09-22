package armeria

import (
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type ConversationManager struct {
	sync.RWMutex
	unsafeConversations []*Conversation
}

type Conversation struct {
	sync.RWMutex
	unsafeCharacter   *Character
	unsafeMobInstance *MobInstance
	unsafeTickCount   int
	ticker            *time.Ticker
	doneCh            chan bool
}

// NewConversationManager returns a new ConversationManager.
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		unsafeConversations: []*Conversation{},
	}
}

// NewConversation starts a new conversation between a unsafeCharacter and a mob.
func (m *ConversationManager) NewConversation() *Conversation {
	m.Lock()
	defer m.Unlock()

	c := &Conversation{}

	m.unsafeConversations = append(m.unsafeConversations, c)

	return c
}

// Delete removes a converastion from the conversation manager.
func (m *ConversationManager) Delete(c *Conversation) {
	m.Lock()
	defer m.Unlock()

	for i, convo := range m.unsafeConversations {
		if convo.Character().ID() == c.Character().ID() {
			m.unsafeConversations[i] = m.unsafeConversations[len(m.unsafeConversations)-1]
			m.unsafeConversations = m.unsafeConversations[:len(m.unsafeConversations)-1]
			break
		}
	}
}

// Character returns the unsafeCharacter in the conversation.
func (convo *Conversation) Character() *Character {
	convo.RLock()
	defer convo.RUnlock()

	return convo.unsafeCharacter
}

// SetCharacter sets the unsafeCharacter having the conversation.
func (convo *Conversation) SetCharacter(c *Character) {
	convo.Lock()
	defer convo.Unlock()

	convo.unsafeCharacter = c
}

// MobInstance returns the mob instance in the conversation.
func (convo *Conversation) MobInstance() *MobInstance {
	convo.RLock()
	defer convo.RUnlock()

	return convo.unsafeMobInstance
}

// SetMobInstance sets the mob instance having the conversation.
func (convo *Conversation) SetMobInstance(mi *MobInstance) {
	convo.Lock()
	defer convo.Unlock()

	convo.unsafeMobInstance = mi
}

// TickCount returns the current tick count.
func (convo *Conversation) TickCount() int {
	convo.RLock()
	defer convo.RUnlock()

	return convo.unsafeTickCount
}

// IncTickCount increments the tick count by 1.
func (convo *Conversation) IncTickCount() {
	convo.Lock()
	defer convo.Unlock()

	convo.unsafeTickCount = convo.unsafeTickCount + 1
}

// Start starts the conversation.
func (convo *Conversation) Start() {
	convo.Lock()
	defer convo.Unlock()

	convo.ticker = time.NewTicker(time.Second)
	convo.doneCh = make(chan bool)

	go func() {
		for {
			select {
			case <-convo.doneCh:
				return
			case <-convo.ticker.C:
				convo.IncTickCount()
				go CallMobFunc(
					convo.Character(),
					convo.MobInstance(),
					"conversation_tick",
					lua.LNumber(convo.TickCount()),
				)
			}
		}
	}()
}

// Cancel stops the conversation, removes it from the unsafeCharacter and the manager.
func (convo *Conversation) Cancel() {
	convo.ticker.Stop()
	convo.doneCh <- true
	Armeria.convoManager.Delete(convo)
	convo.unsafeCharacter.SetMobConvo(nil)
}
