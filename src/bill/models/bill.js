
var axios = require('axios');

var dbpool = require('../lib/dbpool');
var gtab = require('./gtab');

// ============================================================================

function check_amount(q) {
    var conf_prod = gtab.bill_product[q.prod_id];
    if (!conf_prod) return "e_cfg";

    var conf_rate = gtab.bill_rate[q.ccy.toUpperCase()];
    if (!conf_rate) return "e_rate";

    var rate = conf_rate[conf_prod.ccy];
    if (!rate) return "e_rate";

    if (q.amount * rate + 3 < conf_prod.price * q.discount) return "e_amt";

    return "";
}

function get_gservice_addr(svr) {
    var config = require('../../config.json');

    var gs = config.games[svr];
    if (!gs) return null;

    return `http://${gs.svc}/bill/service`;
}

async function update_status(id, st) {
    let db = dbpool.get("bill");
    await db.collection('order').updateOne({ _id: id }, { $set: { status: st } });
}

// ============================================================================

async function give_items(q) {
    try {
        // update sync_ts
        let db = dbpool.get("bill");
        await db.collection('order').updateOne({ _id: q._id }, { $set: { sync_ts: new Date() } });

        // check amount
        let status = check_amount(q);
        if (status != "") {
            await update_status(q._id, status);
            throw `invalid amount: ${status}`;
        }

        // find which server the player is in
        db = dbpool.get("share");
        let doc = await db.collection('userinfo').findOne({ _id: q.userid }, { projection: { "svr": 1 } });
        if (!doc) {
            await update_status(q._id, 'e_user');
            throw `invalid userid: ${q.userid}`;
        }

        // give items. forward to gs
        var url = get_gservice_addr(doc.svr);
        if (!url) throw `${doc.svr} NOT found`;

        let res2 = await axios.post(url, {
            key: "give_items",
            order_key: q._id,
            userid: q.userid,
            prod_id: q.prod_id,
            csext: q.csext,
            amount: q.amount,
            ccy: q.ccy,
            orderid: q.orderid,
            cp_orderid: q.cp_orderid,
        });
        let body2 = res2.data;

        // update final status
        await update_status(q._id, body2);

    } catch (e) {
        console.error("give_items():", e);
    }
}

async function arrange_incomplete_orders() {
    try {
        let db = dbpool.get("bill");

        let docs = await db.collection('order')
            .find({ status: 'payed' })
            .limit(300)
            .toArray();

        let i = 0;
        let now = Date.now();
        docs.forEach(row => {
            if (i >= 200) return;
            if (now - row.sync_ts.getTime() < 5 * 60 * 1000) return;

            setTimeout(function () {
                give_items(row);
            }, i * 300);

            i++;
        });

    } catch (e) {
        console.log("arrange_incomplete_orders():", e);
    }
}

// ============================================================================

module.exports = {
    give_items: give_items,
    arrange_incomplete_orders: arrange_incomplete_orders,
}
