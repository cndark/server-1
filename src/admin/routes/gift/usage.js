
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var session = require('../../models/session');

// ============================================================================

const C_filters = [
    ["--",          "--"],
    ["grpid:n",     "批次号"],
    ["area:n",      "区域Id"],
    ["code:s",      "礼包码"],
    ["userid:s",    "玩家Id"],
    ["svr:n",       "服Id"],
];

const C_rpp = 30;

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

    // page
    q.page = tonumber(q.page);
    if (q.page < 1) q.page = 1;

    // ok
    return q;
}

// ============================================================================

router.get('/', _A_(async (req, res) => {
    var sess = session.data(req);
    var q    = make_q(req.query);

    // render func
    var render = (tab) => {
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
            ],

            tab: {
                name:   '礼包码使用情况',
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
    var skip = (q.page - 1) * C_rpp;

    try {
        let db = dbpool.get('share');

        let docs = await db.collection('giftuse')
            .find(cond)
            .project({_id: 0})
            .sort({ts: -1})
            .skip(skip)
            .limit(C_rpp)
            .toArray();

        var tab = {};

        tab.header = ["批次号", "区域Id", "礼包码", "玩家Id", "服Id", "使用时间"];

        tab.body = docs.map(row => [
            row.grpid,
            row.area,
            row.code,
            row.userid,
            row.svr,
            row.ts.toString(),
        ]);

        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
