package dispatch

import (
	"fmt"
	"reflect"
	"sync"
)

// Dispatch functions onto goroutine in the background. Returns a channel that
// is closed when all goroutines have terminated.
func Dispatch(fs ...func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		var wg sync.WaitGroup
		for _, f := range fs {
			wg.Add(1)
			go func(f func()) {
				defer wg.Done()
				f()
			}(f)
		}
		wg.Wait()
	}()

	return done
}

// Wait waits for multiple signal channels to end
func Wait(chs ...chan struct{}) {
	for _, ch := range chs {
		for range ch {
			// Pass
		}
	}
}

// Close closes multiple channels
func Close(chs ...interface{}) {
	for _, ch := range chs {
		if reflect.TypeOf(ch).Kind() == reflect.Chan {
			reflect.ValueOf(ch).Close()
		}
	}
}

// Split splits a channel into multiple channel
// The input and output channels should be of the same type
func Split(chIn interface{}, chsOut ...interface{}) {
	if reflect.TypeOf(chIn).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot split from value of type %T", chIn))
	}
	for {
		msg, ok := reflect.ValueOf(chIn).Recv()
		if !ok {
			return
		}
		for _, chOut := range chsOut {
			Send(chOut, msg)
		}
	}
}

// Send sends a msg to a channel or an array of channels
func Send(chOut interface{}, msgValue reflect.Value) {
	switch reflect.TypeOf(chOut).Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < reflect.ValueOf(chOut).Len(); i++ {
			if reflect.ValueOf(chOut).Index(i).Kind() != reflect.Chan {
				panic(fmt.Sprintf("cannot send to type %T", chOut))
			}
			reflect.ValueOf(chOut).Index(i).Send(msgValue)
		}
	case reflect.Chan:
		reflect.ValueOf(chOut).Send(msgValue)
	default:
		panic(fmt.Sprintf("cannot send to type %T", chOut))
	}
}

// Merge merges multiple channels of into a channel
// The input and output channels should be of the same type
func Merge(chOut interface{}, chsIn ...interface{}) {
	if reflect.TypeOf(chOut).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot merge to type %T", chOut))
	}

	var wg sync.WaitGroup

	mergeCh := func(chIn interface{}) {
		defer wg.Done()
		for {
			msg, ok := reflect.ValueOf(chIn).Recv()
			if !ok {
				return
			}
			reflect.ValueOf(chOut).Send(msg)
		}
	}

	for _, chIn := range chsIn {
		switch reflect.TypeOf(chIn).Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < reflect.ValueOf(chIn).Len(); i++ {
				if reflect.ValueOf(chIn).Index(i).Kind() != reflect.Chan {
					panic(fmt.Sprintf("cannot merge from value of type %T", chIn))
				}
				wg.Add(1)
				go mergeCh(reflect.ValueOf(chIn).Index(i).Interface())
			}
		case reflect.Chan:
			wg.Add(1)
			go mergeCh(chIn)
		default:
			panic(fmt.Sprintf("cannot merge from value of type %T", chOut))
		}
	}

	wg.Wait()
}

// Pipe all values from a producer channel to a consumer channel until the
// producer is closed, and empty, or until the done channel is closed.
// The consumer channel must not be closed until the Pipe function has
// returned.
func Pipe(done <-chan struct{}, producer interface{}, consumer interface{}) {
	// Type guard the interface inputs
	if reflect.TypeOf(producer).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot pipe from type %T", producer))
	}
	if reflect.TypeOf(consumer).Kind() != reflect.Chan {
		panic(fmt.Sprintf("cannot pipe to type %T", consumer))
	}
	for {
		cases := [2]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(producer)},
		}
		i, val, ok := reflect.Select(cases[:])
		if i == 0 || !ok {
			return
		}

		cases = [2]reflect.SelectCase{
			reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)},
			reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(consumer), Send: val},
		}
		i, val, ok = reflect.Select(cases[:])
		if i == 0 {
			return
		}
	}
}