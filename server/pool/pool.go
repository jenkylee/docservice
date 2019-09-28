// Package pool 通用资源池定义
package pool

import (
	"errors"
	"io"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)

type factory func() (io.Closer, error)

type Pool interface {
	Acquire() (io.Closer, error) // 获取资源
	Release(io.Closer) error     // 释放资源
	Close(io.Closer) error       // 关闭资源
	Shutdown() error             // 关闭池
}

type CommonPool struct {
	sync.Mutex
	pool        chan io.Closer
	maxOpen     int  // 池中最大资源数
	numOpen     int  // 当前池中资源数
	minOpen     int  // 池中最少资源数
	closed      bool // 池是否已关闭
	maxLifetime time.Duration
	factory     factory // 创建连接的方法
}

func NewCommonPool(minOpen, maxOpen int, maxLifetime time.Duration, factory factory) (*CommonPool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}

	p := &CommonPool{
		maxOpen:        maxOpen,
		minOpen:        minOpen,
		maxLifetime:    maxLifetime,
		factory:        factory,
		pool: make(chan io.Closer, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		closer, err := factory()
		if err != nil {
			continue
		}

		p.numOpen++
		p.pool <- closer
	}

	return p, nil
}
// 获取资源
func (p *CommonPool) Acquire() (io.Closer, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}

	for {
		closer, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		// todo maxLifttime处理
		return closer, nil
	}
}

// 实际获取或者创建资源
func (p *CommonPool) getOrCreate() (io.Closer, error) {
	select {
	case closer := <-p.pool:
		return closer, nil
	default:
	}

	p.Lock()
	if p.numOpen >= p.maxOpen {
		closer := <-p.pool
		p.Unlock()
		return closer, nil
	}

	// 新建连接
	closer, err := p.factory()
	if err != nil {
		p.Unlock()
		return nil, err
	}

	p.numOpen++
	p.Unlock()

	return closer, nil
}

// 释放单个资源到资源池
func (p *CommonPool) Release(closer io.Closer) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	p.pool <- closer
	p.Unlock()

	return nil
}

// 关闭单个资源
func (p CommonPool) Close(closer io.Closer) error {
	p.Lock()
	closer.Close()
	p.numOpen--
	p.Unlock()

	return nil
}

// 关闭连接池， 释放所有资源
func (p CommonPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	close(p.pool)
	for closer := range p.pool {
		closer.Close()
		p.numOpen--
	}
	p.closed = true
	p.Unlock()

	return nil
}