package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}
func Wait() {
	globalCloser.Wait()
}
func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex     // защита от одновременного доступа
	once  sync.Once      // гарантирует, что CloseAll() сработает только один раз
	done  chan struct{}  // канал для ожидания завершения
	funcs []func() error // список функций для вызова при завершении
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(sig) > 0 { // Если переданы сигналы (например, os.Interrupt), запускаем горутину
		go func() {
			ch := make(chan os.Signal, 1) // Канал для получения сигналов от ОС
			signal.Notify(ch, sig...)     // Подписываемся на сигналы ОС
			<-ch                          // Ожидаем получения сигнала (блокировка)
			signal.Stop(ch)               // Отписываемся от сигналов (освобождаем ресурсы)
			c.CloseAll()                  // Запускаем завершение при получении сигнала
		}()
	}
	return c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done // Блокируем выполнение до тех пор, пока канал done не будет закрыт
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		//собираем ошибки в канал
		errs := make(chan error, len(funcs))
		//пробегаем по переданным функициям
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
