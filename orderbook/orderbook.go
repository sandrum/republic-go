package orderbook

import (
	"context"
	"encoding/base64"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

type Listener interface {
	OnOrderFragmentReceived(order.Fragment)
}

type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server hosted at the identity.MultiAddress.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

type Server interface {
	OpenOrder(context.Context, order.EncryptedFragment) error
}

type Orderbook interface {
	Server
	Syncer

	// OrderFragment stored in this local Orderbook. These are received from
	// other Orderbooks calling Orderbook.OpenOrder to send an
	// order.EncryptedFragment to this local Orderbook.
	OrderFragment(order.ID) (order.Fragment, error)

	// Order that has been reconstructed and stored in this local Orderbook.
	// This only happens for orders that have been matched and confirmed.
	Order(order.ID) (order.Order, error)
}

type orderbook struct {
	crypto.RsaKey

	syncer   Syncer
	storer   Storer
	listener Listener
}

func NewOrderbook(key crypto.RsaKey, syncer Syncer, storer Storer, listener Listener) Orderbook {
	return &orderbook{
		RsaKey: key,

		syncer:   syncer,
		storer:   storer,
		listener: listener,
	}
}

func (book *orderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	fragment, err := orderFragment.Decrypt(*book.RsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if fragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, base64.StdEncoding.EncodeToString(fragment.OrderID[:8]), base64.StdEncoding.EncodeToString(fragment.ID[:8]))
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, base64.StdEncoding.EncodeToString(fragment.OrderID[:8]), base64.StdEncoding.EncodeToString(fragment.ID[:8]))
	}

	book.listener.OnOrderFragmentReceived(fragment)
	return book.storer.InsertOrderFragment(fragment)
}

func (book *orderbook) Sync() (ChangeSet, error) {
	return book.syncer.Sync()
}

func (book *orderbook) OrderFragment(id order.ID) (order.Fragment, error) {
	return book.storer.OrderFragment(id)
}

func (book *orderbook) Order(id order.ID) (order.Order, error) {
	return book.storer.Order(id)
}
