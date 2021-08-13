
var express = require('express');
var router = express.Router();

var axios = require('axios');
var md5 = require('md5');
var qs = require('querystring');
var crypto = require('crypto');

var dbpool = require('../lib/dbpool');
var gtab = require('../models/gtab');
var bill = require('../models/bill');

var config = require('../../config.json');

// ============================================================================

const ERR_Ok = 0;
const ERR_Info = 1;
const ERR_Sdk = 2;
const ERR_Sig = 3;
const ERR_Order = 4;

// ============================================================================

function extract_info(v) {
    let arr = v.split("|");
    if (arr.length < 5) {
        return
    }

    var info = {
        sdk: arr[0],
        userid: arr[1],
        prod_id: arr[2],
        cp_orderid: arr[3],
        csext: arr[4],
    };

    return info;
}

// ============================================================================

async function bill_flow(opt) {
    // extract info
    let info;
    try {
        info = opt.info();
        if (!info) throw -1;
    } catch {
        return ERR_Info;
    }

    // check sdk
    let sdkobj;
    try {
        sdkobj = gtab.sdk[info.sdk];
        if (!sdkobj) throw -1;
    } catch {
        return ERR_Sdk;
    }

    // verify signature
    try {
        if (!opt.sig(sdkobj)) throw -1;
    } catch {
        return ERR_Sig;
    }

    // save order: {orderid, amount, ccy, [discount]}
    try {
        // prepare doc
        let doc = opt.doc;

        if (!doc.orderid || !doc.amount || !doc.ccy)
            return ERR_Order;

        doc._id = `${info.sdk}-${doc.orderid}`;
        doc.area = config.common.area.id;
        doc.sdk = info.sdk;
        doc.userid = info.userid;
        doc.prod_id = tonumber(info.prod_id);
        doc.csext = info.csext;
        doc.discount = doc.discount || 1;
        doc.create_ts = new Date();
        doc.sync_ts = new Date();
        doc.status = 'payed';
        doc.cp_orderid = info.cp_orderid;

        // save
        try {
            let db = dbpool.get("bill");
            await db.collection('order').insertOne(doc);
            await bill.give_items(doc);
        } catch (e) {
            if (e.code != 11000) throw e;
        }

        // ok
        return ERR_Ok;

    } catch {
        return ERR_Order;
    }
}

// ============================================================================
// chuxin proxy
// ============================================================================

async function check_chuxin_proxy(req, res) {
    let q = req.body;

    // check proxy url
    if (!q.payExt) return false;
    let m = q.payExt.match(/\|(https:\/\/[^|]+)$/);
    if (!m) return false;

    // proxy url found, go proxy.
    // check sign, remove proxy info, re-sign
    let removed = false;
    try {
        // check sign
        let info = extract_info(q.payExt);
        if (!info) throw 'no info'

        let sdkobj = gtab.sdk[info.sdk];
        if (!sdkobj) throw 'no sdk';

        let sign = q.sign;
        let cols = ["game", "orderId", "amount", "uid", "zone", "goodsId", "payTime", "payChannel", "payExt"];

        let sign2 = md5(`${cols.map(v => q[v]).join("")}#${sdkobj.bill_key}`);
        if (sign != sign2) throw 'err sign';

        // remove proxy info
        q.payExt = q.payExt.slice(0, -m[0].length);
        removed = true;

        // re-sign
        q.sign = md5(`${cols.map(v => q[v]).join("")}#${sdkobj.bill_key}`);
    } finally {
        if (!removed) {
            q.payExt = q.payExt.slice(0, -m[0].length);
        }
    }

    // proxy now
    try {
        let url = m[1];
        let { data } = await axios.post(url, q);
        res.json(data).end();
    } catch {
        res.status(500).end();
    }

    return true;
}


// ============================================================================
// chuxinh5 proxy
// ============================================================================

async function check_chuxinh5_proxy(req, res) {
    let q = req.body;

    // check proxy url
    if (!q.pay_ext) return false;
    let m = q.pay_ext.match(/\|(https:\/\/[^|]+)$/);
    if (!m) return false;

    // proxy url found, go proxy.
    // check sign, remove proxy info, re-sign
    let removed = false;
    let headers = req.headers;
    try {
        // check sign
        let info = extract_info(q.pay_ext);
        if (!info) throw 'no info'

        let sdkobj = gtab.sdk[info.sdk];
        if (!sdkobj) throw 'no sdk';

        let signStr = JSON.stringify(q) + sdkobj.secretId + headers.algorithm + headers.timestamp;

        var hmac = crypto.createHmac('SHA256', sdkobj.secretKey);
        hmac.update(signStr);
        sign = hmac.digest('hex');
        if (sign != headers.signature) throw 'err sign';

        // remove proxy info
        q.pay_ext = q.pay_ext.slice(0, -m[0].length);
        removed = true;

        // re-sign
        signStr = JSON.stringify(q) + sdkobj.secretId + headers.algorithm + headers.timestamp;
        var hmac = crypto.createHmac('SHA256', sdkobj.secretKey);
        hmac.update(signStr);
        headers.signature = hmac.digest('hex');

    } finally {
        if (!removed) {
            q.pay_ext = q.pay_ext.slice(0, -m[0].length);
        }
    }

    // proxy now
    try {
        let url = m[1];
        let { data } = await axios.post(url, q,
            {
                headers: {
                    "Content-Type": "application/json",
                    secretId: headers.secretid,
                    algorithm: headers.algorithm,
                    signature: headers.signature,
                    timestamp: headers.timestamp,
                },
            }
        );
        res.json(data).end();
    } catch {
        res.status(500).end();
    }

    return true;
}


// ============================================================================
// chuxin
// ============================================================================

router.post('/bill/chuxin', _A_(async (req, res) => {
    if (await check_chuxin_proxy(req, res)) return;

    let q = req.body;

    let ec = await bill_flow({
        info: () => extract_info(q.payExt),

        sig: (sdkobj) => {
            let sign = q.sign;
            let cols = ["game", "orderId", "amount", "uid", "zone", "goodsId", "payTime", "payChannel", "payExt"];

            let sign2 = md5(`${cols.map(v => q[v]).join("")}#${sdkobj.bill_key}`);
            return sign == sign2;
        },

        doc: {
            orderid: q.orderId,
            amount: tonumber(q.amount),
            ccy: "CNY",

            chuxin_uid: q.uid,
            chuxin_goodsId: q.goodsId,
            chuxin_payTime: q.payTime,
            chuxin_payChannel: q.payChannel,
        },
    });

    // resp
    res.json({
        errno: [1000, -1008, -1008, 1001, -1][ec],
        errmsg: ['', 'einfo', 'esdk', 'esig', 'eorder'][ec],
        data: {
            orderId: q.orderId,
            amount: q.amount,
            game: q.game,
            zone: q.zone,
            uid: q.uid,
        },
    }).end();
}));


// ============================================================================
// chuxinea
// ============================================================================

router.post('/bill/chuxinea', _A_(async (req, res) => {
    if (await check_chuxin_proxy(req, res)) return;

    let q = req.body;

    let ec = await bill_flow({
        info: () => extract_info(q.payExt),

        sig: (sdkobj) => {
            let sign = q.sign;
            let cols = ["game", "orderId", "amount", "uid", "zone", "goodsId", "payTime", "payChannel", "payExt"];

            let sign2 = md5(`${cols.map(v => q[v]).join("")}#${sdkobj.bill_key}`);
            return sign == sign2;
        },

        doc: {
            orderid: q.orderId,
            amount: tonumber(q.amount),
            ccy: "USD",

            chuxin_uid: q.uid,
            chuxin_goodsId: q.goodsId,
            chuxin_payTime: q.payTime,
            chuxin_payChannel: q.payChannel,
        },
    });

    // resp
    res.json({
        errno: [1000, -1008, -1008, 1001, -1][ec],
        errmsg: ['', 'einfo', 'esdk', 'esig', 'eorder'][ec],
        data: {
            orderId: q.orderId,
            amount: q.amount,
            game: q.game,
            zone: q.zone,
            uid: q.uid,
        },
    }).end();
}));

// ============================================================================
// chuxin-h5
// ============================================================================

router.post('/bill/chuxinh5', _A_(async (req, res) => {
    if (await check_chuxinh5_proxy(req, res)) return;

    let q = req.body;

    let ec = await bill_flow({
        info: () => extract_info(q.pay_ext),

        sig: (sdkobj) => {
            let headers = req.headers;
            let signStr = JSON.stringify(q) + sdkobj.secretId + headers.algorithm + headers.timestamp;

            var hmac = crypto.createHmac('SHA256', sdkobj.secretKey);
            hmac.update(signStr);
            sign = hmac.digest('hex');
            return sign == headers.signature;
        },

        doc: {
            orderid: q.order_id,
            amount: tonumber(q.amount),
            ccy: "CNY",

            chuxin_uid: q.uid,
            chuxin_goodsId: q.goods_id,
            chuxin_payTime: q.pay_time,
            chuxin_payChannel: q.pay_channel,
            chuxin_sand_box: q.sand_box,
            chuxin_third_order_id: q.third_order_id,
        },
    });

    // resp
    res.json({
        code: [1000, -1008, -1008, 1001, -1][ec],
        msg: ['', 'einfo', 'esdk', 'esig', 'eorder'][ec],
    }).end();
}));



// ============================================================================

module.exports = router;
