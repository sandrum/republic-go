package adapter_test

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/http/adapter"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Status adapter", func() {
	var statusAdapter adapter.StatusAdapter
	var reader testReader

	Context("when retreiving a status from a status adapter", func() {

		BeforeEach(func() {
			multiAddr, err := testutils.RandomMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			reader = testReader{
				network:   "someNetwork",
				multiAddr: multiAddr,
				ethAddr:   "someString",
				publicKey: []byte("Here is a some random stuff for a public key...."),
				peers:     1337,
			}
			statusAdapter = adapter.NewStatusAdapter(reader)
		})

		It("should return the same information that was stored", func() {
			status, err := statusAdapter.Status()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Network).Should(Equal(reader.network))
			Expect(status.MultiAddress).Should(Equal(reader.multiAddr.String()))
			Expect(status.EthereumAddress).Should(Equal(reader.ethAddr))
			hexPk := "0x" + hex.EncodeToString(reader.publicKey)
			Expect(status.PublicKey).Should(Equal(hexPk))
			Expect(status.Peers).Should(Equal(reader.peers))
		})
	})

})

type testReader struct {
	network   string
	multiAddr identity.MultiAddress
	ethAddr   string
	publicKey []byte
	peers     int
}

func (r testReader) Network() (string, error) {
	return r.network, nil
}

func (r testReader) MultiAddress() (identity.MultiAddress, error) {
	return r.multiAddr, nil
}

func (r testReader) EthereumAddress() (string, error) {
	return r.ethAddr, nil
}

func (r testReader) PublicKey() ([]byte, error) {
	return r.publicKey, nil
}

func (r testReader) Peers() (int, error) {
	return r.peers, nil
}
