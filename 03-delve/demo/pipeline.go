package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime/pprof"
	"sync"
	"time"
)

// Order represents a customer order to process
type Order struct {
	ID       int
	Priority int
	Items    []string
}

// ProcessedOrder represents an order after processing
type ProcessedOrder struct {
	Order
	ProcessedAt time.Time
	ProcessedBy string
}

// Pipeline processes orders through multiple stages
type Pipeline struct {
	incoming   chan Order
	validation chan Order
	processing chan Order
	shipping   chan ProcessedOrder
	done       chan bool
	wg         sync.WaitGroup
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		incoming:   make(chan Order),          // BUG: unbuffered channel
		validation: make(chan Order),          // BUG: unbuffered channel
		processing: make(chan Order),          // BUG: unbuffered channel
		shipping:   make(chan ProcessedOrder), // BUG: unbuffered channel
		done:       make(chan bool),
	}
}

// Start initializes all pipeline workers
func (p *Pipeline) Start() {
	// Start validation worker
	p.wg.Add(1)
	go p.validator()

	// Start processing workers
	for i := 1; i <= 2; i++ {
		p.wg.Add(1)
		go p.processor(fmt.Sprintf("processor-%d", i))
	}

	// Start shipping worker
	p.wg.Add(1)
	go p.shipper()

	// Start order receiver
	p.wg.Add(1)
	go p.receiver()
}

// receiver takes orders from incoming and sends to validation
func (p *Pipeline) receiver() {
	ls := pprof.Labels("job", "receiver")
	pprof.Do(context.Background(), ls, func(_ context.Context) {
		defer p.wg.Done()
		for order := range p.incoming {
			log.Printf("[RECEIVER] Received order %d with priority %d\n", order.ID, order.Priority)
			// BUG: This can block if validator is busy
			p.validation <- order
			log.Printf("[RECEIVER] Sent order %d to validation\n", order.ID)
		}
		close(p.validation)
	})
}

// validator checks orders and sends to processing
func (p *Pipeline) validator() {
	ls := pprof.Labels("job", "validator")
	pprof.Do(context.Background(), ls, func(_ context.Context) {
		defer p.wg.Done()
		for order := range p.validation {
			log.Printf("[VALIDATOR] Validating order %d\n", order.ID)

			// Simulate validation work
			time.Sleep(100 * time.Millisecond)

			if len(order.Items) == 0 {
				log.Printf("[VALIDATOR] Order %d rejected: no items\n", order.ID)
				continue
			}

			log.Printf("[VALIDATOR] Order %d passed validation\n", order.ID)
			// BUG: This can block if all processors are busy
			p.processing <- order
		}
		close(p.processing)
	})
}

// processor handles order processing
func (p *Pipeline) processor(name string) {
	ls := pprof.Labels("job", "processor")
	pprof.Do(context.Background(), ls, func(_ context.Context) {
		defer p.wg.Done()
		for order := range p.processing {
			log.Printf("[%s] Processing order %d\n", name, order.ID)

			// Simulate processing work (varies by priority)
			duration := time.Duration(500-order.Priority*100) * time.Millisecond
			time.Sleep(duration)

			processed := ProcessedOrder{
				Order:       order,
				ProcessedAt: time.Now(),
				ProcessedBy: name,
			}

			log.Printf("[%s] Completed order %d, sending to shipping\n", name, order.ID)
			// BUG: This can block if shipper is busy
			p.shipping <- processed
		}
	})
}

// shipper handles the final shipping stage
func (p *Pipeline) shipper() {
	pprof.Label(context.TODO(), "shipper")
	defer p.wg.Done()
	defer close(p.done)

	shipped := 0
	// BUG: shipper closes after 10 orders, but we might send more
	for shipped < 10 {
		select {
		case order := <-p.shipping:
			log.Printf("[SHIPPER] Shipping order %d (processed by %s at %s)\n",
				order.ID, order.ProcessedBy, order.ProcessedAt.Format("15:04:05"))

			// Simulate shipping work
			time.Sleep(200 * time.Millisecond)

			shipped++
			log.Printf("[SHIPPER] Shipped %d orders so far\n", shipped)
		case <-time.After(5 * time.Second):
			log.Printf("[SHIPPER] Timeout waiting for orders, shutting down\n")
			return
		}
	}
	log.Printf("[SHIPPER] Shipped maximum orders (%d), shutting down\n", shipped)
}

// SendOrder sends a new order to the pipeline
func (p *Pipeline) SendOrder(order Order) {
	log.Printf("[MAIN] Sending order %d to pipeline\n", order.ID)
	p.incoming <- order // BUG: This will block if receiver is blocked
	log.Printf("[MAIN] Order %d sent successfully\n", order.ID)
}

// Wait waits for all pipeline workers to complete
func (p *Pipeline) Wait() {
	close(p.incoming)
	p.wg.Wait()
}

func generateOrders(count int) []Order {
	orders := make([]Order, count)
	for i := 0; i < count; i++ {
		numItems := rand.Intn(5) // Some orders might have 0 items (invalid)
		items := make([]string, numItems)
		for j := 0; j < numItems; j++ {
			items[j] = fmt.Sprintf("item-%d", j+1)
		}

		orders[i] = Order{
			ID:       i + 1,
			Priority: rand.Intn(5) + 1,
			Items:    items,
		}
	}
	return orders
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	log.Println("Starting order processing pipeline...")

	pipeline := NewPipeline()
	pipeline.Start()

	// Generate and send orders
	orders := generateOrders(15) // BUG: Generating 15 orders but shipper stops at 10

	// Send orders concurrently to simulate real load
	for _, order := range orders {
		// Send each order in a goroutine to avoid blocking main
		pipeline.SendOrder(order)

		// Small delay between orders
		time.Sleep(50 * time.Millisecond)
	}

	log.Println("All orders submitted, waiting for pipeline to complete...")

	// Wait for pipeline to finish or timeout
	done := make(chan bool)
	go func() {
		pipeline.Wait()
		done <- true
	}()

	select {
	case <-done:
		log.Println("Pipeline completed successfully")
	case <-time.After(30 * time.Second):
		log.Println("Pipeline timeout - possible deadlock!")
		// Program will hang here, perfect for Delve inspection
		select {} // Hang forever to allow debugging
	}
}
