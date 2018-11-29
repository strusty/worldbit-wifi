


function pay(selector, value, id) {
        paypal.Button.render({
            // Configure environment
            env: 'sandbox',
            client: {
                sandbox: 'AWDP8bjlJObrXU2QmM-vD1dTUZYGiCikh6c1gkP1-mh3KL1Svkp3ueNNUMpQQwRWsrbNFox7_Zd_ccsT',
                production: '<insert production client id>'
            },
            // Customize button (optional)
            locale: 'en_US',
            style: {
                layout: 'vertical',  // horizontal | vertical
                size: 'medium',    // medium | large | responsive
                shape: 'rect',      // pill | rect
                color: 'gold'       // gold | blue | silver | black
            },

            funding: {
                // allowed: [paypal.FUNDING.CARD, paypal.FUNDING.CREDIT],
                disallowed: []
            },


            // Set up a payment
            payment: function (data, actions) {
                return actions.payment.create({
                    transactions: [{
                        amount: {
                            total: value,
                            currency: 'USD'
                        }
                    }]
                });
            },
            // Execute the payment
            onAuthorize: function (data, actions) {
                return actions.payment.execute().then(function (response) {
                    console.log(1234567,response);
                    methods.post(host + '/paypal/payment', "POST", PayPalVoucherRequest = {

                        SaleID : response.transactions[0].related_resources[0].sale.id,
                        PricingPlanID: id,
                        PhoneNumber : localStorage.getItem('PhoneNumber'),
                    }, function (el) {
                        console.log('capch', el);

                    });


                    // Show a confirmation message to the buyer
                    window.alert('Thank you for your purchase!');
                });
            }
        }, selector)
    }

