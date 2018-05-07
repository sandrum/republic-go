package dispatch_test

import (
	"errors"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Error channels", func() {

	Context("Merge errors", func() {

		It("should merge multiple error channels", func() {
			errCh1 := make(chan error, 1)
			errCh2 := make(chan error, 1)
			errCh3 := make(chan error, 1)

			err1 := errors.New("1")
			err2 := errors.New("2")
			err3 := errors.New("3")

			errCh1 <- err1
			errCh2 <- err2
			errCh3 <- err3

			errCh := MergeErrors(errCh1, errCh2, errCh3)
			time.Sleep(time.Second)

			Close(errCh1, errCh2, errCh3)
			Ω(len(errCh)).Should(Equal(3))
		})

		It("should be able to read the errors originated from all the error channels", func() {
			errCh1 := make(chan error, 1)
			errCh2 := make(chan error, 1)
			errCh3 := make(chan error, 1)

			err1 := errors.New("1")
			err2 := errors.New("2")
			err3 := errors.New("3")

			errCh1 <- err1
			errCh2 <- err2
			errCh3 <- err3

			errCh := MergeErrors(errCh1, errCh2, errCh3)

			time.Sleep(time.Second)
			Ω(len(errCh)).Should(Equal(3))

			Close(errCh1, errCh2, errCh3)

			for err := range errCh {
				Ω(err == err1 || err == err2 || err == err3).Should(BeTrue())
			}

		})

	})

	Context("Filter errors", func() {
		It("should filter errors using a predicate", func() {
			errCh := make(chan error, 3)

			err1 := errors.New("1")
			err2 := errors.New("20")
			err3 := errors.New("300")

			errCh <- err1
			errCh <- err2
			errCh <- err3

			predicate := func(err error) bool {
				if len(err.Error()) == 2 {
					return true
				}
				return false
			}

			filteredErrCh := FilterErrors(errCh, predicate)

			time.Sleep(1 * time.Second)
			Ω(len(filteredErrCh)).Should(Equal(1))

			err := <-filteredErrCh

			Ω(err.Error()).Should(Equal("20"))

			Close(errCh)
		})
	})

	Context("Consume errors", func() {
		It("should be able to process an error", func() {
			errCh := make(chan error, 3)
			defer close(errCh)

			err1 := errors.New("1")
			err2 := errors.New("20")
			err3 := errors.New("300")

			errCh <- err1
			errCh <- err2
			errCh <- err3

			consumeFn := func(err error) {
				log.Println("Processing the error", err.Error())
			}

			go ConsumeErrors(errCh, consumeFn)
			time.Sleep(1 * time.Second)
		})
	})
})