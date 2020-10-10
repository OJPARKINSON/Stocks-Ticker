const jsdom = require('jsdom');
const fetch = require('node-fetch');

exports.handler = async (event, context, callback) => {
    const StocksCalculator = (buyPrice, sellPrice, stockAmount) => {
        return (buyPrice * stockAmount - sellPrice * stockAmount).toLocaleString('en-UK', { style: 'currency', currency: 'GBP' });
    }

    const stockPrice = async (stockCode) => {
        const data = await fetch(`https://in.finance.yahoo.com/quote/${stockCode}`)
        const pageBody = await data.text()
        const dom = await new jsdom.JSDOM(await pageBody, 'text/html');
        return parseFloat(await dom.window.document.querySelector('span[data-reactid="32"]').textContent.replace(/,/g, ''))
    }

    const currencyConverter = async (stockPrice) => {
        const data = await fetch(`https://in.finance.yahoo.com/quote/USDGBP=X`)
        const pageBody = await data.text()
        const { window } = await new jsdom.JSDOM(await pageBody, 'text/html');
        const exchangeRate = parseFloat(await window.document.querySelector('span[data-reactid="32"]').textContent.replace(/,/g, ''));
        return parseInt((await exchangeRate * stockPrice).toFixed(2))
    }

    const main = async () => {
        const xrpPrice = await stockPrice('XRP-GBP');
        const cmcsaPrice = await stockPrice('CMCSA');
        const convertedCmcsaPrice = await currencyConverter(cmcsaPrice);

        return ({ 
            statusCode: 200, 
            body: JSON.stringify({ 
                xrpPrice: `£${await xrpPrice}`, 
                cmcsaPrice: `£${await convertedCmcsaPrice}`, 
                xrpProf: `${await StocksCalculator(6, await xrpPrice, 1000)}`, 
                cmcsaProf: `${await StocksCalculator(await convertedCmcsaPrice, 27.26, 330)}`
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
