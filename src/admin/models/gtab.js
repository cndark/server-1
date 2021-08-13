
var config = require('../../config.json');

var tab_bill = require('../../gamedata/data/billProduct.json');
var tab_rate = require('../../gamedata/data/billRate.json');
var tab_ccy = require('../../gamedata/data/currency.json');
var tab_item = require('../../gamedata/data/item.json');
var tab_from = require('../../gamedata/data/dictFrom.json');
var tab_csext = require('../../gamedata/data/dictCsExt.json');
var tab_lang = require('../../gamedata/data/language.json');
var tab_hero = require('../../gamedata/data/monster.json');
var tab_relic = require('../../gamedata/data/relic.json');

// ============================================================================

module.exports = {

    conf: (function () {
        return [
            "billProduct",
            "actBillFirst",
            "item",
        ].map(v => [v, v]);
    })(),

    lang: (function () {
        return tab_lang.map(v => [v.key, v.text]);
    })(),

    // --------------------------------

    toCNY: (function () {
        let m = {};
        tab_rate.forEach(v => {
            m[v.key] = v;
        });

        return (v) => v * m[config.admin.ccy].CNY / 100;
    })(),

    // --------------------------------

    res: (function () {
        let ccy = tab_ccy.map(v => [v.id, v.txt_Name]);
        let item = tab_item.map(v => [v.id, v.txt_Name]);
        let hero = tab_hero.map(v => {
            return (v.fragment && v.fragment.length > 0) ? [v.id, v.name] : [];
        });
        let relic = tab_relic.map(v => [v.id, v.txt_Name]);

        return ccy.concat(item).concat(hero).concat(relic);
    })(),

    hero: (function () {
        return [[0, "--"]].concat(tab_hero.map(v => {
            return (v.fragment && v.fragment.length > 0) ? [v.id, v.name] : [];
        }).sort((a, b) => a[0] - b[0]));
    })(),

    // --------------------------------

    dict_res: (function (id) {
        let m = new Map();

        tab_ccy.forEach(v => {
            m.set(v.id, v.txt_Name);
        });
        tab_item.forEach(v => {
            m.set(v.id, v.txt_Name);
        });
        tab_hero.forEach(v => {
            if (v.fragment && v.fragment.length > 0) {
                m.set(v.id, v.name);
            }
        });
        tab_relic.forEach(v => {
            m.set(v.id, v.txt_Name);
        });

        return (id) => {
            let name = m.get(id);
            return name ? name : id;
        }
    })(),

    dict_bill: (function () {
        let m = new Map();

        tab_bill.forEach(v => {
            m.set(v.id, v.txt_Name);
        });
        tab_csext.forEach(v => {
            m.set(v.id, v.txt_Name);
        });

        return (id) => {
            let name = m.get(id);
            return name ? name : id;
        }
    })(),

    dict_from: (function () {
        let m = new Map();

        tab_from.forEach(v => {
            m.set(v.id, v.name);
        });

        return (id) => {
            let name = m.get(id);
            return name ? name : id;
        }
    })(),

};
