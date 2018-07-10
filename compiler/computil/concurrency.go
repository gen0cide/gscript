package computil

import "sync"

// ExecuteFuncsInParallel is a meta function that takes an array of function pointers (hopefully for each VM)
// and executes them in parallel to decrease compile times. This is setup to handle errors within
// each VM gracefully and not allow a goroutine to fail silently.
func ExecuteFuncsInParallel(fns []func() error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for _, fn := range fns {
		wg.Add(1)
		go func(f func() error) {
			err := f()
			if err != nil {
				errChan <- err
			}
			wg.Done()
		}(fn)
	}
	go func() {
		wg.Wait()
		close(finChan)
	}()
	select {
	case <-finChan:
	case err := <-errChan:
		if err != nil {
			return err
		}
	}
	return nil
}
