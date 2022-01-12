package main

import (
	"sync"
)

type Task interface {
	Execute()
}
type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

//NewPool - конструктор пула горутин
func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, 256), // по-умолчанию, создаём буферезированный канал мошностью 256 элементов
		kill:  make(chan struct{}),  // канал для завершения 1 горутины
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	for {
		select {
		case task, isOpened := <-p.tasks:
			if !isOpened {
				continue
			}
			task.Execute()
			p.wg.Done()
		case <-p.kill:
			return
		}
	}
}

//Resize - изменяем размер пула гоуртин
func (p *Pool) Resize(size int) {
	//	требуется изменение кол-ва вопркеров
	if p.size != size {
		p.mu.Lock()
		//	нужно увеличить кол-во воркеров
		for p.size < size {
			// инкрементируем кол-во воркеров
			p.size++
			//	...и запускаем горутину воркера
			go p.worker()
		}
		//	нужно уменьшить кол-во воркеров
		for p.size > size {
			// декрементируем кол-во воркеров
			p.size--
			//	... и отправляем в канал "сигнал", что горуттине надо прекратит своё исполнение
			//	какая горутина первой считает из этого канала, та и завершит своё выпонение
			p.kill <- struct{}{}
		}
		p.mu.Unlock()
	}
}

//Close - закрытие пула горутин воркеров
func (p *Pool) Close() {
	//	в итерации завершаем все запущенные горутины воркеров
	for p.size > 0 {
		p.kill <- struct{}{}
		p.size--
	}
	//	закрываем канал с заданиями для воркеров
	close(p.tasks)
}

//Wait - ожидаем, когда закончит работу WaitGroup
func (p *Pool) Wait() {
	p.wg.Wait()
}

//Exec - добавляем задание для воркера
func (p *Pool) Exec(t Task) {
	p.wg.Add(1)
	p.tasks <- t
}
