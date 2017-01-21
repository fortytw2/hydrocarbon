package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/internal/httputil"
	"github.com/fortytw2/hydrocarbon/internal/log"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/event"
	"github.com/stripe/stripe-go/sub"
)

func activateAccount(s *hydrocarbon.Store, l log.Logger) httputil.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		stripe.Key = os.Getenv("STRIPE_PRIVATE_KEY")

		user := loggedIn(r)
		if user == nil {
			http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
			return nil
		}
		token := r.FormValue("stripeToken")
		if token == "" {
			return httputil.Wrap(errors.New("no stripeToken present"), http.StatusForbidden)
		}
		// Create a Customer:
		customerParams := &stripe.CustomerParams{
			Email: user.Email,
		}
		err := customerParams.SetSource(token)
		if err != nil {
			return httputil.Wrap(err, http.StatusBadRequest)
		}

		customer, err := customer.New(customerParams)
		if err != nil {
			return httputil.Wrap(err, http.StatusBadRequest)
		}

		_, err = sub.New(&stripe.SubParams{
			Customer: customer.ID,
			Plan:     "standard",
		})
		if err != nil {
			return httputil.Wrap(err, http.StatusInternalServerError)
		}

		err = s.Users.SetStripeCustomerID(user.ID, customer.ID)
		if err != nil {
			return httputil.Wrap(err, http.StatusInternalServerError)
		}

		http.Redirect(w, r, settingsURL, http.StatusSeeOther)

		return nil
	}
}

func stripeWebhookHandler(s *hydrocarbon.Store, l log.Logger) httputil.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {

		var hookEvent stripe.Event
		if err := json.NewDecoder(io.LimitReader(r.Body, 16384)).Decode(&hookEvent); err != nil {
			return err
		}

		// Check authenticity with stripe
		checkedEvent, err := event.Get(hookEvent.ID, nil)
		if err != nil {
			// client error or event doesn't exist
			if serr, ok := err.(*stripe.Error); ok && serr.HTTPStatusCode == http.StatusNotFound {
				return httputil.Wrap(err, http.StatusForbidden)
			}
			l.Log("msg", "could not get stripe event", "err", err)
			return httputil.Wrap(err, http.StatusInternalServerError)
		}
		if checkedEvent == nil {
			// I don't think this should occur
			l.Log("msg", "could not verify stripe event, checked event is nil", "err", err)
			return httputil.Wrap(err, http.StatusForbidden)
		}

		stripeLog := log.NewContext(l).With("stripe_event_type", checkedEvent.Type,
			"stripe_event_id", checkedEvent.ID,
			"stripe_event_live", checkedEvent.Live,
			"stripe_event_req", checkedEvent.Req,
			"stripe_event_user_id", checkedEvent.UserID)

		switch checkedEvent.Type {
		case "customer.subscription.deleted":
			var sub stripe.Sub
			err := json.Unmarshal(checkedEvent.Data.Raw, &sub)
			if err != nil {
				return httputil.Wrap(err, http.StatusBadRequest)
			}

			stripeLog.Log("msg", "subscription canceled", "stripe_sub_id", sub.ID, "stripe_customer_id", sub.Customer.ID)
			// err = s.CancelSubscription(sub)
			// if err != nil {
			// 	return httputil.Wrap(err, http.StatusInternalServerError)
			// }
		case "customer.subscription.created":

		case "customer.subscription.trial_will_end":

		default:
			stripeLog.Log("msg", "unhandled stripe callback!", "data", string(checkedEvent.Data.Raw))
		}
		return nil
	}
}
