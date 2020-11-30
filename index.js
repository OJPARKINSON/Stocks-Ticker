const jsdom = require('jsdom');
const fetch = require('node-fetch');
var crypto = require('crypto');
require('dotenv').config()

var apiKey = process.env.apiKey;
var apiSecret = process.env.apiSecret;

exports.handler = async (event, context, callback) => {

    const coinbaseRequest = async (path) => {
        var timestamp = Math.floor(Date.now() / 1000);
        var options = {
            method: 'GET',
            headers: {
                'CB-ACCESS-SIGN': crypto.createHmac("sha256", apiSecret).update(timestamp + 'GET' + path).digest("hex"),
                'CB-ACCESS-TIMESTAMP': timestamp,
                'CB-ACCESS-KEY': apiKey,
                'CB-VERSION': '2015-07-22'
            }
        };

        return fetch(`https://api.coinbase.com/${path}`, options)
            .then(res => res.json())
            .then(json => json.data)
            .then(data => data.native_balance ? data.native_balance.amount : data.amount);
    };

    const stockPrice = async (stockCode) => {
        const data = await fetch(`https://in.finance.yahoo.com/lookup?s=${stockCode}`)
        const pageBody = await data.text()
        const dom = await new jsdom.JSDOM(await pageBody, 'text/html');
        return parseFloat(await dom.window.document.querySelector('td[data-reactid="59"]').textContent.replace(/,/g, ''))
    }

    const main = async () => {
        const cmcsaPrice = await stockPrice('CMCSA') * process.env.averageGBP;

        return ({
            statusCode: 200, 
            body: JSON.stringify({ 
                xrpPrice: `£${await coinbaseRequest('/v2/prices/XRP-GBP/buy')}`,
                xrpProf: `£${await coinbaseRequest('/v2/accounts/7524fa83-38cc-5a0e-a29b-ec9555d2657c')}`,  
                cmcsaPrice: `£${cmcsaPrice}`,
                cmcsaProf: `${(cmcsaPrice * 330 - 27.26 * 330).toLocaleString('en-UK',{style:'currency',currency:'GBP'})}`
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
