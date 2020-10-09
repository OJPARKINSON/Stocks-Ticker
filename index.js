const jsdom = require('jsdom');
const fetch = require('node-fetch');
require('dotenv').config();

// exports.handler = async (event) => {
    const StocksCalculator = (buyPrice, sellPrice, stockAmount) => 
        (buyPrice * stockAmount - sellPrice * stockAmount).toLocaleString('en-UK', { style: 'currency', currency: 'GBP' });
    
    const stockPrice = (stockCode) => 
        fetch(`https://in.finance.yahoo.com/quote/${stockCode}`)
            .then(response => response.text())
            .then(pageBody => new jsdom.JSDOM(pageBody, 'text/html'))
            .then(({window}) => parseFloat(window.document.querySelector('span[data-reactid="32"]').textContent.replace(/,/g, '')))
            .then(bonk => { console.log(bonk); return bonk})
            .catch(err => console.error(err));

    const currencyConverter = (stockPrice) => 
        fetch(`https://finance.yahoo.com/quote/USDGBP=X`)
            .then(response => response.text())
            .then(pageBody => new jsdom.JSDOM(pageBody, 'text/html'))
            .then(({window}) => parseFloat(window.document.querySelector('span[data-reactid="32"]').textContent.replace(/,/g, '')))
            .then(USD_GBP => USD_GBP * stockPrice)
            .then(ukStockPrice => parseInt(ukStockPrice.toFixed(2)))
            .catch(err => console.error(err));
    
    const main = async () => {
        const xrpPrice = await stockPrice('XRP-GBP');
        const cmcsaPrice = await stockPrice('CMCSA');
        const convertedCmcsaPrice = await currencyConverter(cmcsaPrice);

        return { 
            xrpPrice: `£${xrpPrice}`, 
            cmcsaPrice: `£${convertedCmcsaPrice}`, 
            xrpProf: `${StocksCalculator(process.env.XRPSELLPRICE, xrpPrice, 1000)}`, 
            cmcsaProf: `${StocksCalculator(convertedCmcsaPrice, process.env.CMCSABUYPRICE, 330)}`
        }
    
        // console.log(` XRP💸: ${StocksCalculator(6, xrpPrice, 1000)}`)
        // console.log(` Sharesave 💸: ${StocksCalculator(cmcsaPrice, 27.26, 330)}`)
    }
     main()
      .then(response => console.log({ statusCode: 200, body: JSON.stringify(response) }))
        .catch(err => console.error(err));
// };
