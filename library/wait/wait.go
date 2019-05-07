package wait

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/xingshanghe/sirius/runtime"
)

// NeverStop 永不停止
var NeverStop <-chan struct{} = make(chan struct{})

// Group 允许启动一组协程并等待它们完成.
type Group struct {
	wg sync.WaitGroup
}

// Wait 阻塞
func (g *Group) Wait() {
	g.wg.Wait()
}

// Start 在一个新的协程中执行函数f.
func (g *Group) Start(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

// StartWithChannel 通过channel启动一组协程
// stopCh 作为参数之一，当stopCh可用时f停止执行.
func (g *Group) StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{})) {
	g.Start(func() {
		f(stopCh)
	})
}

// StartWithContext 通过上下文启动一组协程.
// ctx 作为参数之一，当ctx.Done()可用时f停止执行.
func (g *Group) StartWithContext(ctx context.Context, f func(context.Context)) {
	g.Start(func() {
		f(ctx)
	})
}

// Forever 永远执行f.
func Forever(f func(), period time.Duration) {
	Until(f, period, NeverStop)
}

// Until 周期内循环执行f,直到stop channel 关闭
// sliding = true,f函数执行完成 后 计时
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}

// NonSlidingUntil 周期内循环执行f,直到stop channel 关闭
// sliding = false,f函数执行 前 计时
func NonSlidingUntil(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, false, stopCh)
}

// JitterUntil 在stop channel被结束之前循环执行f
// sliding 计时的起点，如果是TRUE的话，就是f函数执行完成后计时.
func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	var t *time.Timer
	var sawTimeout bool

	for {
		select {
		case <-stopCh:
			return
		default:
		}

		jitteredPeriod := period
		if jitterFactor > 0.0 {
			jitteredPeriod = Jitter(period, jitterFactor)
		}

		if !sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}

		func() {
			defer runtime.HandleCrash()
			f()
		}()

		if sliding {
			t = resetOrReuseTimer(t, jitteredPeriod, sawTimeout)
		}
		// NOTE：B/C在Golang没有优先级选择,这意味着我们可以触发T.C和Stopch，T.C选择失败.
		// 为了减轻压力,我们在每个循环开始时重新检查stopch,以防止额外执行f.
		select {
		case <-stopCh:
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}

// Jitter 抖动返回time.Duration类型介于duration和duration + maxFactor * duration之间.
// 轻微的抖动减少在心跳等操作时误判
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// resetOrReuseTimer 如果一个计时器已经使用，避免分配一个新的计时器.
// 并发不安全.
func resetOrReuseTimer(t *time.Timer, d time.Duration, sawTimeout bool) *time.Timer {
	if t == nil {
		return time.NewTimer(d)
	}
	if !t.Stop() && !sawTimeout {
		<-t.C
	}
	t.Reset(d)
	return t
}
