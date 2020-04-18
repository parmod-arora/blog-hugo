---
title: "One time checkout with Stripe"
date: 2020-04-11T18:04:30+08:00
draft: false
---

There are a lot of Online Payment processing services these days and Stripe is one of them. In this article we will have a look on one time payment processing with stripe and how to  integrate checkout flow in our project.

### Overview
I heard about stripe on online and i thought why don't I gave a short and see it myself. Frankly speaking I like stripe because.
1. It has great [documentation](https://stripe.com/docs), which i liked the most and also it helped me to build my sample project.
1. Many other good features like subscription, save credit card etc. which i will talk in abother post.
1. It provide integration with all major cards scheme VISA, MASTERCARD, AMEX, Discover, JCB etc and other wallets like Google Pay and Apple Pay

### Requirement
1. HTTPS on production but can use over http on development env
1. Requires a server to create the session. 

### Checkout flow

![A flowchart of the Checkout flow](/images/checkout-one-time-client-server.png " ")

Before prcessing any one time payment stripe want client application to create a session on their server, which represents the intent to purchase by the customer. The Checkout Sessions API provide flexibility for dynamic amounts and line items. 



<!-- [code](projects/stripe_integration/cmd/serverd/v1/checkout/session.go) -->
``` golang
  stripe.Key = os.Getenv("stripe_secret")
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SuccessURL: stripe.String("http://localhost:3000/client/success.html?session_id={CHECKOUT_SESSION_ID}&order_id=" + model.OrderID),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
	}
	for _, orderItem := range model.Items {
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(orderItem.Name),
			Description: stripe.String(orderItem.Description),
			Amount:      stripe.Int64(orderItem.Amount),
			Currency:    stripe.String(string(stripe.CurrencySGD)),
			Quantity:    stripe.Int64(1),
		})
	}
	session, err := session.New(params)
```

After succesfully creating a session on stripe from our server, our app will redirect the customer to stripe payment form page using below code. where user can provide its personal and credit card information. This whole process is handled by stripe.js loaded in our app page.Its recomendded to load stripe js from `https://js.stripe.com` instead of your web server.

``` js
  <script src="https://js.stripe.com/v3/"></script>
  <script>
    let stripe = Stripe("stripe_key")
    stripe.redirectToCheckout({
      sessionId: sessionId
    })
  </script>
```


### Test cards
|NUMBER|BRAND|CVC|DATE|
|------|-------|---|----|
|4242424242424242|	Visa|	Any 3 digits	|Any future date|
|4000056655665556|	Visa (debit)|	Any 3 digits	|Any future date|
|5555555555554444|	Mastercard|	Any 3 digits	|Any future date|
|2223003122003222|	Mastercard (2-series)|	Any 3 digits	|Any future date|
|5200828282828210|	Mastercard (debit)|	Any 3 digits	|Any future date|
|5105105105105100|	Mastercard (prepaid)|	Any 3 digits	|Any future date|
|378282246310005|	American Express|	Any 4 digits	|Any future date|
|371449635398431|	American Express|	Any 4 digits	|Any future date|
|6011111111111117|	Discover|	Any 3 digits	|Any future date|
|6011000990139424|	Discover|	Any 3 digits	|Any future date|
|3056930009020004|	Diners Club	|Any 3 digits	|Any future date|
|36227206271667|	Diners Club (14 digit card)|	Any 3 digits	|Any future date|
|3566002020360505|	JCB|	Any 3 digits	|Any future date|
|6200000000000005|	UnionPay|	Any 3 digits	|Any future date|

The CheckoutSessions API allows for dynamic amounts and line items but requires a server to create the session. 


