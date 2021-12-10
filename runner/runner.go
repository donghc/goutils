package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

//Runner 在给定的超时时间内执行一组任务
// 并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	//从操作系统发送信号
	interrupt chan os.Signal
	//报告处理任务已完成
	complete chan error
	//报告处理任务已经超时
	timeout <-chan time.Time
	//持有一组以索引为顺序的依次执行的以 int 类型 id 为参数的函数
	tasks []func(id int)
}

// 统一错误处理

var (
	// ErrTimeOut 超时错误信息，会在任务执行超时的时候返回
	ErrTimeOut = errors.New("received timeout")
	// ErrInterrupt 中断错误信号，会在接收到操作系统的事件时返回
	ErrInterrupt = errors.New("received interrupt")
)

//New 函数返回一个新的准备使用的Runner， d：自定义分配的时间
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		//会在另一个线程经过时间段d 后想返回值发送当时的时间点
		timeout: time.After(d),
	}
}

// Add 将任务添加到Runner中
func (r *Runner) Add(task ...func(id int)) {
	r.tasks = append(r.tasks, task...)
}

//检测是否收到了中断信号
func (r *Runner) gotInterrupt() bool {
	select {
	//当终端事件被触发
	case <-r.interrupt:
		//停止接收后续的任何信号
		signal.Stop(r.interrupt)
		return true
		//继续执行
	default:
		return false
	}
}

//run 执行每一个已注册的任务
func (r *Runner) run() error {
	for id, task := range r.tasks {
		//检测操作系统的中断信号
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

// Start 开始执行所有任务，并监控通道事件
func (r *Runner) Start() error {
	//监控所有的中断信号
	signal.Notify(r.interrupt, os.Interrupt)
	//使用不同的 goroutine 执行不同的任务
	go func() {
		r.complete <- r.run()
	}()
	//使用select 语句来监控goroutine的通信
	select {
	//等待任务完成
	case err := <-r.complete:
		return err
	//任务超时
	case <-r.timeout:
		return ErrTimeOut
	}
}
