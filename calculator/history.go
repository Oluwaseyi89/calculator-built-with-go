package calculator

import (
	"time"
)

type HistoryEntry struct {
	Expression string
	Result     float64
	Timestamp  time.Time
}

func (c *Calculator) GetHistory() []HistoryEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entries := make([]HistoryEntry, len(c.history))
	for i, entry := range c.history {
		// Parse entry back to HistoryEntry (simplified)
		entries[i] = HistoryEntry{
			Expression: entry,
			Timestamp:  time.Now(),
		}
	}
	return entries
}

func (c *Calculator) ClearHistory() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history = make([]string, 0)
}

func (c *Calculator) SaveHistoryToFile(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Implementation for saving history to file
	// ...
	return nil
}
