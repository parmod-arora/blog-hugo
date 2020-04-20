+++
title= "One time checkout with Stripe"
cover = "/images/stripe/social.png"
date= 2020-04-11T18:04:30+08:00
draft= false
tags= [
    "stripe",
    "golang",
		"payments"
]
+++

There are a lot of Online Payment processing services these days and Stripe is one of them. In this article we will have a look at one time payment processing with stripe and how to  integrate checkout flow in our project.

### Overview
I heard about stripe on online and i thought why don't I gave a short and see it myself. Frankly speaking I like stripe because.
1. It has great [documentation](https://stripe.com/docs), which i liked the most and also it helped me to build my sample project.
1. Many other good features like subscription, save credit card etc. which i will talk in abother post.
1. It provide integration with all major cards scheme VISA, MASTERCARD, AMEX, Discover, JCB etc and other wallets like Google Pay and Apple Pay

### Requirement
1. HTTPS on production but can use over http on development env
1. Requires a server to create the session. 

### Checkout flow
- Create Checkout Session
- Redirect to Checkout
- Confirm payment status


![A flowchart of the Checkout flow](/images/stripe/checkout-one-time-client-server.png " ")


Create Checkout Session 
------
Before prcessing any payment stripe want client application to create a session on their server, which represents the intent to purchase by the customer. The Checkout Sessions API provide flexibility for dynamic amounts and line items. 



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

Redirect to Checkout 
------
After succesfully creating a session on stripe from our server, our app will redirect the customer to stripe payment form page using below code. where user can provide its personal and credit card information. This whole process is handled by stripe.js loaded in our app page.Its recomendded to load stripe js from `https://js.stripe.com`.

``` js
  <script src="https://js.stripe.com/v3/"></script>
  <script>
    let stripe = Stripe("stripe_key")
    stripe.redirectToCheckout({
      sessionId: sessionId
    })
  </script>
```

It will load stripe hosted page on testing env (bank OTP page in production envionment) where user can continue to make payment. Challenge flow will be handled by stripe we don't have to worry about the 3DS and chargeback cases, As these cases are handled by stripe.

Confirm payment status
------
When your customer completes a payment, Stripe redirects them to the URL that you specified in the success_url parameter. Typically, this is a page on your website that informs your customer that their payment was successful.
Stripe sends the `checkout.session.completed` event for a successful Checkout payment. 

``` golang
// EventCallback evant callback from stripe
func EventCallback(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// If you are testing your webhook locally with the Stripe CLI you
	// can find the endpoint's secret by running `stripe listen`
	// Otherwise, find your endpoint's secret in your webhook settings in the Developer Dashboard
	endpointSecret := os.Getenv("stripe_webhook_secret")

	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		order.UpdateSessionStatus(r.Context(), session.ID, order.StatusPaid)
	default:
		fmt.Fprintf(os.Stderr, "Unexpected event type: %s\\n", event.Type)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
```

### Test cards
|NUMBER|BRAND|CVC|DATE|REMARKS|
|------|-------|---|----|-----|
|4242424242424242|	Visa|	Any 3 digits	|Any future date|Default U.S. card|
|4000000000003220|	Visa (debit)|	Any 3 digits	|Any future date|Authenticate with 3D Secure|


[Source Code](https://github.com/parmod-arora/blog-hugo/tree/master/projects/stripe_integration)

``` sh
$ git clone https://github.com/parmod-arora/blog-hugo.git 
$ cd ./projects/stripe_integration/
$ make run 
```

