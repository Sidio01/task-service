package email

import (
	"context"
	"log"
	"sync"
	"time"

	"gitlab.com/g6834/team26/task/internal/domain/models"
)

type EmailToStdOut struct {
	EmailChan     chan models.Email
	ResultChan    chan map[models.Email]bool
	RateLimitChan chan struct{}
	ctx           context.Context
	cancelFunc    context.CancelFunc
	wg            *sync.WaitGroup
	nWorkers      int
}

func New(ctx context.Context, nWorkers, nRate int) (*EmailToStdOut, error) {
	newCtx, cancelFunc := context.WithCancel(ctx)
	return &EmailToStdOut{

		EmailChan:     make(chan models.Email, nWorkers),
		ResultChan:    make(chan map[models.Email]bool, nWorkers),
		RateLimitChan: make(chan struct{}, nRate),
		ctx:           newCtx,
		cancelFunc:    cancelFunc,
		nWorkers:      nWorkers,
		wg:            &sync.WaitGroup{},
	}, nil
}

func (etso *EmailToStdOut) Stop() error {
	etso.cancelFunc()
	etso.wg.Wait()
	return nil
}

func (etso *EmailToStdOut) SendEmail(e models.Email) error {
	if e.Type == "approve" {
		log.Printf("Task %s - sending email to %s, type: you need to %s the task\n", e.TaskUUID, e.Reciever, e.Type)
	} else {
		log.Printf("Task %s - sending email to %s, type: task was %s\n", e.TaskUUID, e.Reciever, e.Type)
	}
	return nil
}

func (etso *EmailToStdOut) PushEmailToChan(e models.Email) {
	etso.EmailChan <- e
}

func (etso *EmailToStdOut) GetEmailResultChan() chan map[models.Email]bool {
	return etso.ResultChan
}

func (etso *EmailToStdOut) StartEmailWorkers() {
	for i := 0; i < etso.nWorkers; i++ {
		etso.wg.Add(1)
		go etso.EmailWorker()
	}
}

func (etso *EmailToStdOut) EmailWorker() {
	defer etso.wg.Done()
	log.Println("Worker started")
	for {
		select {
		case <-etso.ctx.Done():
			log.Println("Recieved signal to stop worker")
			return
		case etso.RateLimitChan <- struct{}{}:
			select {
			case <-etso.ctx.Done():
				log.Println("Recieved signal to stop worker")
				return
			case email := <-etso.EmailChan:
				err := etso.SendEmail(email)
				if err != nil {
					time.Sleep(60 * time.Second)
					continue
				}
				result := make(map[models.Email]bool)
				result[email] = true
				etso.ResultChan <- result
				<-etso.RateLimitChan
			}
		}
	}
}
