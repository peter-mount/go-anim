package renderer

import (
	util3 "git.area51.dev/peter/videoident/util"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util"
	"sync"
	"sync/atomic"
	"time"
)

type Parallelize struct {
	threads    int                // number of threads
	jobs       chan Context       // worker queue
	resultWG   sync.WaitGroup     // result wait group
	workerWG   sync.WaitGroup     // worker wait group
	queue      util.PriorityQueue // queue to collate
	task       Renderer           // Renderer to generate work
	output     Renderer           // Renderer to receive results
	nextFrame  atomic.Int32       // frame to receive
	closed     bool               // set when last task has been consumed
	maxPending int                // Number of entries in queue to hold back workers
	pending    atomic.Int32       // Number in queue
	collate    bool               // true to ensure output called in sequence, e.g. ffmpeg
}

func (r Renderer) Parallel(threads, start int, collate bool, output Renderer) Renderer {
	threads = util3.ThreadLimit(threads)

	// If 1 thread then do nothing
	if threads < 2 {
		return r.Then(output)
	}

	log.Printf("Using %d threads", threads)

	p := &Parallelize{
		threads: threads,
		task:    r,
		output:  output,
		collate: collate,
		jobs:    make(chan Context, threads),
	}

	p.nextFrame.Store(int32(start))

	if p.collate {
		p.maxPending = threads
		go p.consumer()
	}

	for i := 0; i < p.threads; i++ {
		go p.worker(i)
	}

	return p.render
}

// Render parallelize the rendering
func (p *Parallelize) render(ctx Context) error {
	p.jobs <- CloneContext(ctx)

	if !ctx.HasNext() {
		p.Close()
		if p.collate {
			p.resultWG.Wait()
		}
	}

	return nil
}

func (p *Parallelize) Close() {
	p.closed = true
	close(p.jobs)
}

func (p *Parallelize) worker(_ int) {
	for {
		ctx, ok := <-p.jobs

		if !ok {
			return
		}

		// Wait until the consumer has caught up
		for p.collate && int(p.pending.Load()) >= p.maxPending {
			time.Sleep(time.Millisecond * 5)
		}

		err := p.task(ctx)
		if err != nil {
			p.add(ctx, err)
			p.Close()
		} else {
			p.add(ctx, ctx)
		}
	}
}

func (p *Parallelize) add(ctx Context, v interface{}) {
	if p.collate {
		p.queue.AddPriority(ctx.Frame(), v)
		p.pending.Add(1)
	} else {
		if err, ok := v.(error); ok {
			panic(err)
		}
		if ctx, ok := v.(Context); ok {
			p.consume(ctx)
		}
	}
}

func (p *Parallelize) consumer() {
	p.resultWG.Add(1)
	defer func() {
		p.resultWG.Add(-1)
	}()

	for p.pending.Load() > 0 || !p.queue.IsEmpty() || !p.closed {
		if r, ok := p.queue.Pop(); ok {

			if err, ok := r.(error); ok {
				// TODO handle errors better
				panic(err)
			}

			if ctx, ok := r.(Context); ok {
				if p.collate && ctx.Frame() > int(p.nextFrame.Load()) {
					p.queue.AddPriority(ctx.Frame(), ctx)
					time.Sleep(time.Second / 4)
				} else {
					p.pending.Add(-1)
					p.consume(ctx)
					p.nextFrame.Add(1)
				}
			}
		} else {
			// sleep a bit for next result
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (p *Parallelize) consume(ctx Context) {
	err := p.output(ctx)
	if err != nil {
		p.Close()
		// TODO handle errors better
		panic(err)
	}
}
