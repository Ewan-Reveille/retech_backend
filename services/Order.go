package services

import (
    "errors"
    "github.com/Ewan-Reveille/retech/models"
    "github.com/google/uuid"
    "github.com/stripe/stripe-go/v81"
    "github.com/stripe/stripe-go/v81/paymentintent"
    "gorm.io/gorm"
)

type OrderService struct {
    Repo        models.OrderRepository
    PaymentRepo models.PaymentRepository
    ProductRepo models.ProductRepository
    DB          *gorm.DB
}

// CreateWithStripe kicks off a Stripe payment and returns the client secret
func (os *OrderService) CreateWithStripe(
    buyerID uuid.UUID,
    productID uuid.UUID,
    shippingID uuid.UUID,
) (clientSecret string, orderID uuid.UUID, err error) {

    // 1) Fetch the product to get price & stripe priceID
    prod, err := os.ProductRepo.GetByID(productID)
    if err != nil {
        return "", uuid.Nil, err
    }

    // 2) Calculate amount in cents
    amountCents := int64(prod.Price * 100)

    // 3) Start a DB transaction so we don’t get partial state
    err = os.DB.Transaction(func(tx *gorm.DB) error {
        // 4) Create Stripe PaymentIntent
        piParams := &stripe.PaymentIntentParams{
            Amount:   stripe.Int64(amountCents),
            Currency: stripe.String(string(stripe.CurrencyEUR)),
            // You can attach metadata, e.g. our own orderID once we have one
            Metadata: map[string]string{
                "buyer_id":    buyerID.String(),
                "product_id":  productID.String(),
                "shipping_id": shippingID.String(),
            },
        }
        pi, err := paymentintent.New(piParams)
        if err != nil {
            return err
        }

        // 5) Create our local Payment record
        payment := &models.Payment{
            PaymentID: uuid.New(),
            UserID:    buyerID,
            Amount:    prod.Price,
            OrderID:   uuid.Nil, // we’ll fill in after creating the order
            Method:    models.MethodCard,          // or MethodPaypal / Crypto
            Status:    models.StatusPending,       // until webhook confirms
            // You may want a field stripe_intent_id in your model; if so:
            // StripeIntentID: pi.ID,
        }
        if err := os.PaymentRepo.Create(payment); err != nil {
            return err
        }

        // 6) Create the Order record, linking the payment
        order := &models.Order{
            ID:         uuid.New(),
            BuyerID:    buyerID,
            ProductID:  productID,
            ShippingID: shippingID,
            PaymentID:  payment.PaymentID,
            Status:     "pending",
            Commission: 0,
        }
        if err := os.Repo.Create(order); err != nil {
            return err
        }

        // 7) Update the Payment.OrderID now that order exists
        payment.OrderID = order.ID
        if err := os.PaymentRepo.Update(payment); err != nil {
            return err
        }

        // return out values
        clientSecret = pi.ClientSecret
        orderID = order.ID
        return nil
    })

    return clientSecret, orderID, err
}


func (os *OrderService) GetByID(id uuid.UUID) (*models.Order, error) {
	order, err := os.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (os *OrderService) GetByBuyerID(buyerId uuid.UUID) ([]models.Order, error) {
	orders, err := os.Repo.GetByBuyerID(buyerId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// func (os *OrderService) GetBySellerId(sellerId uuid.UUID) ([]models.Order, error) {
// 	orders, err := os.Repo.GetBySellerId(sellerId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

func (os *OrderService) GetByProductID(productID uuid.UUID) ([]models.Order, error) {
	orders, err := os.Repo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *OrderService) Update(order *models.Order) error {
	existing, err := os.Repo.GetByID(order.ID)
	if err != nil {
		return err
	}

	if existing.BuyerID != order.BuyerID {
		return errors.New("buyer ID mismatch")
	}
	if existing.ProductID != order.ProductID {
		return errors.New("product ID mismatch")
	}

	return os.Repo.Update(order)
}

func (os *OrderService) Delete(id uuid.UUID) error {
	order, err := os.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.BuyerID == uuid.Nil {
		return errors.New("buyer ID cannot be empty")
	}
	if order.ProductID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}

	if order.ProductID == uuid.Nil {
		return errors.New("product ID cannot be empty")
	}
	return os.Repo.Delete(id)
}
