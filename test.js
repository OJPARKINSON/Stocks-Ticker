var crypto = require('crypto');
const fetch = require('node-fetch');
require('dotenv').config()

var apiKey = process.env.apiKey;
var apiSecret = process.env.apiSecret;

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
        .then(data => console.log({data}));
};

coinbaseRequest('/v2/accounts/7524fa83-38cc-5a0e-a29b-ec9555d2657c');
coinbaseRequest('/v2/prices/XRP-GBP/buy');