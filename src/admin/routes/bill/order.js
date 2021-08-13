
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var session = require('../../models/session');
var gtab = require('../../models/gtab');

// ============================================================================

const C_filters = [
    ["--",      "--"],
    ["orderid:s",  "订单号"],
    ["area:n",     "区域"],
    ["svr:s",      "服务器名"],
    ["sdk:re",     "Sdk"],
    ["userid:s",   "玩家Id"],
    ["prod_id:n",  "产品Id"],
    ["status:s",   "状态"],
];

const C_rpp = 30;

// ============================================================================

const order_status = {
    payed:  "已支付-未发货",
    ok:     "完成-已发货",
    e_cfg:  "产品未配置",
    e_rate: "汇率未配置",
    e_amt:  "金额非法",
    e_user: "用户非法",
};

function status_text(st) {
    var txt = order_status[st];
    return txt ? txt : st;
}

// ============================================================================

function make_q(q) {
    // filter
    if (!q.fkey) {
        q.fkey = [];
        q.fval = [];
    } else if (!(q.fkey instanceof Array)) {
        q.fkey = [q.fkey];
        q.fval = [q.fval];
    }

    q.filters = [];
    for (let i = 0; i < q.fkey.length; i++) {
        if (q.fkey[i] == "--" || q.fval[i] == "") continue;
        q.filters.push([q.fkey[i], q.fval[i]]);
    }

    // dr
    if (!q.dr) {
        let now = new Date();
        q.dr = `${now.addDay(-15).toDateString()} ~ ${now.toDateString()}`;
    }

    let arr = q.dr.split(" ~ ");
    q.d_start = Date.fromString(arr[0]).startOfDay();
    q.d_end   = Date.fromString(arr[1]).endOfDay();

    // page
    q.page = tonumber(q.page);
    if (q.page < 1) q.page = 1;

    // ok
    return q;
}

// ============================================================================

router.get('/', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = make_q(req.query);

    // render func
    let render = (tab) => {
        tab = tab || {};

        res.render('lib/table/paged', {
            sess: sess,

            form: [
                {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[0], icon:  "filter"},
                {type: "input",     name: "fval", placeholder: "值", value: q.fval[0], style: "width: 120px"},
                {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[1], icon:  "filter"},
                {type: "input",     name: "fval", placeholder: "值", value: q.fval[1], style: "width: 120px"},
                {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[2], icon:  "filter"},
                {type: "input",     name: "fval", placeholder: "值", value: q.fval[2], style: "width: 120px"},
                {type: "daterange", name: "dr",                      value: q.dr,      icon:  "calendar-alt"},
            ],

            tab: {
                name:   '订单列表',
                header: tab.header,
                body:   tab.body,
            },

            paged: {
                page:  q.page,
                rec_n: tab.body ? tab.body.length : 0,
                rpp:   C_rpp,
            },
        });
    }

    // do it
    let cond = utils.filter({}, q.filters);
    cond.create_ts = {$gte: q.d_start, $lt: q.d_end};
    let skip = (q.page - 1) * C_rpp;

    try {
        let db = dbpool.get("bill");

        let docs = await db.collection('order')
            .find(cond)
            .sort({create_ts: -1})
            .skip(skip)
            .limit(C_rpp)
            .toArray();

        // tab
        var tab = {};

        // header
        tab.header = ["订单号", "区域", "Sdk", "玩家Id", "产品Id", "cs透传", "金额", "币种", "折扣", "创建时间", "同步时间", "状态"];

        // body
        tab.body = docs.map(row => [
            row.orderid,
            row.area,
            row.sdk,
            row.userid,
            gtab.dict_bill(row.prod_id),
            gtab.dict_bill(row.csext),
            row.amount,
            row.ccy,
            row.discount,
            row.create_ts.toString(),
            row.sync_ts.toString(),
            status_text(row.status),
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
