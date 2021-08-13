
var tab_sdk  = require('../../gamedata/data/sdk.json');
var tab_prod = require('../../gamedata/data/billProduct.json');
var tab_rate = require('../../gamedata/data/billRate.json');

// ============================================================================

module.exports = {

    sdk: (function () {
        var m = {};
        tab_sdk.forEach(v => {
            m[v.name] = v;
        });
        return m;
    })(),

    bill_product: (function () {
        var m = {};
        tab_prod.forEach(v => {
            m[v.id] = v
        });
        return m;
    })(),

    bill_rate: (function () {
        var m = {};
        tab_rate.forEach(v => {
            m[v.key] = v;
        });
        return m;
    })(),

};
