package calculator

import (
	"fmt"
	"sync"
)

type Calculator struct {
	memory     float64
	lastResult float64
	history    []string
	mu         sync.RWMutex
}

func NewCalculator() *Calculator {
	return &Calculator{
		memory:     0,
		lastResult: 0,
		history:    make([]string, 0),
	}
}

func (c *Calculator) SetMemory(value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory = value
}

func (c *Calculator) GetMemory() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.memory
}

func (c *Calculator) AddToMemory(value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory += value
}

func (c *Calculator) SubtractFromMemory(value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory -= value
}

func (c *Calculator) ClearMemory() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory = 0
}

func (c *Calculator) SetLastResult(value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastResult = value
}

func (c *Calculator) GetLastResult() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastResult
}

func (c *Calculator) AddToHistory(expression string, result float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := fmt.Sprintf("%s = %.6g", expression, result)
	c.history = append(c.history, entry)

	// Keep only last 100 entries
	if len(c.history) > 100 {
		c.history = c.history[1:]
	}
}

func (c *Calculator) ShowHistory() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.history) == 0 {
		fmt.Println("No history available")
		return
	}

	fmt.Println("\nCalculation History:")
	for i, entry := range c.history {
		fmt.Printf("%3d. %s\n", i+1, entry)
	}
}
