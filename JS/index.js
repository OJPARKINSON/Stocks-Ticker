const fetch = require('node-fetch');
var crypto = require('crypto');
if (process.env.NODE_ENV === "Production") require('dotenv').config()
else require('dotenv').config({ path: '../.env' })

const { apiKey, apiSec, token, acountID } = process.env;

exports.handler = async (event, context, callback) => {
    const coinbaseRequest = async (path) => {
        var timestamp = Math.floor(Date.now() / 1000);
        var options = {
            method: 'GET',
            headers: {
                'CB-ACCESS-SIGN': crypto.createHmac("sha256", apiSec).update(timestamp + 'GET' + path).digest("hex"),
                'CB-ACCESS-TIMESTAMP': timestamp,
                'CB-ACCESS-KEY': apiKey,
                'CB-VERSION': '2015-07-22'
            }
        };

        return fetch(`https://api.coinbase.com/${path}`, options)
        .then(res => res.json())
        .then(json => json.data)
    };

    const cmcsaRequest = async () => {
        const options = "&types=quote"
        return fetch(`https://cloud.iexapis.com/stable/stock/cmcsa/batch?token=${token}${options}`)
            .then(res => res.json())
    }

    const main = async () => {
        const cmcsaPrice = await cmcsaRequest();
        const exchangeRate = await coinbaseRequest('/v2/exchange-rates?currency=USD')
        const GPBCMCSA = exchangeRate.rates.GBP * cmcsaPrice.quote.latestPrice
        const portfolio = await coinbaseRequest(`/v2/accounts/${acountID}`)

        return ({
            statusCode: 200, 
            body: JSON.stringify({
                xrpProf: `£${portfolio.native_balance.amount}`,  
                cmcsaPrice: cmcsaPrice.quote.latestPrice.toLocaleString('en-UK',{style:'currency',currency:'USD'}),
                cmcsaProf: (GPBCMCSA * 330 - 27.26 * 330).toLocaleString('en-UK',{style:'currency',currency:'GBP'})
            })
        })
    }

    try {
        return main();
    } catch(error) {
        return {
            "statusCode": error.statusCode,
            "headers": {
              "Content-Type": "text/plain",
              "x-amzn-ErrorType": error.code
            },
            "isBase64Encoded": false,
            "body": error.code + ": " + error.message
          }
    }
};
