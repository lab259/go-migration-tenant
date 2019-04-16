package mtnt

import (
	"sync"

	"github.com/lab259/go-prdcsm"
)

// AccountProducer should list all the accounts that should be migrated.
type AccountProducer interface {
	Get() ([]Account, error)
	HasNext() bool
	Total() int
}

// accountProducerProxy bridges between a AccountProducer and a
// `prdcsm.Producer`.
type accountProducerProxy struct {
	executor         *MigrationExecutor
	producer         AccountProducer
	ch               chan interface{}
	listAccountsOnce sync.Once
	running          bool
}

func (p *accountProducerProxy) Stop() {
	if p.running {
		p.running = false
		close(p.ch)
	}
}

// newAccountProducerProxy returns a new instance of a `accountProducerProxy`.
func newAccountProducerProxy(executor *MigrationExecutor, producer AccountProducer, channelLen int) *accountProducerProxy {
	return &accountProducerProxy{
		executor: executor,
		producer: producer,
		ch:       make(chan interface{}, channelLen),
	}
}

// listAccounts will iterate through the whole set of accounts adding it to the
// the `.ch`.
func (p *accountProducerProxy) listAccounts() {
	p.running = true
	go func() {
		for p.producer.HasNext() {
			list, err := p.producer.Get()
			if err != nil {
				panic(err) // TODO What should be done here?
			}
			for _, account := range list {
				p.ch <- account
			}
		}

		// Reached the end.
		p.running = false
		close(p.ch)
	}()
}

// Produce uses the `.listAccountsOnce` to extract all accounts from the
// `AccountProducer` queueing it to `.ch`.
func (p *accountProducerProxy) Produce() interface{} {
	p.listAccountsOnce.Do(p.listAccounts)

	acc, more := <-p.ch
	if !more {
		return prdcsm.EOF
	}
	return acc
}
