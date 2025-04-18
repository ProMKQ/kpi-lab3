package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури через виконання операцій.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлена останнього разу

	mq   messageQueue
	wg   sync.WaitGroup
	stop chan struct{}
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Запустіть цей метод перед використанням Post.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.mq.ops = make(chan Operation, 100)
	l.stop = make(chan struct{})

	go func() {
		for {
			select {
			case op := <-l.mq.ops:
				if op.Do(l.next) {
					l.Receiver.Update(l.next)
					l.next, l.prev = l.prev, l.next
				}
				l.wg.Done()
			}
		}
	}()
}

// Post додає нову операцію до внутрішньої черги.
func (l *Loop) Post(op Operation) {
	l.wg.Add(1)
	l.mq.ops <- op
}

// StopAndWait зупиняє цикл і чекає його завершення.
func (l *Loop) StopAndWait() {
	l.wg.Wait()
}

// messageQueue реалізує просту чергу операцій через канал.
type messageQueue struct {
	ops chan Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.ops <- op
}

func (mq *messageQueue) pull() Operation {
	return <-mq.ops
}

func (mq *messageQueue) empty() bool {
	return len(mq.ops) == 0
}
