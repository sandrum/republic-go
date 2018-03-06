package compute_test

import (
	"math/big"
	"time"

	"github.com/republicprotocol/go-order-compute"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Computations", func() {

	n := int64(35)
	k := int64(24)
	prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

	Context("when orders are an exact match", func() {

		It("should find a match", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

	})

	Context("when orders use different currencies", func() {

		It("should not find a match for the same currencies in reverse", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeETH, compute.CurrencyCodeBTC, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the first currencies differ", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeETH, compute.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the second currencies differ", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when both currencies differ", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeREN, compute.CurrencyCodeDGD, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when prices vary", func() {

		It("should find a match when the buy price is higher", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(12), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the buy price is lower", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(12), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when volumes vary", func() {

		It("should find a match when the maximum buy volume is higher than the maximum sell volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the maximum sell volume is higher than the maximum buy volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum buy volume is higher than the minimum sell volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum sell volume is higher than the maximum buy volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the maximum buy volume is lower than the minimum sell volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the maximum sell volume is lower than the minimum buy volume", func() {
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result, err := computeFinalFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})
})