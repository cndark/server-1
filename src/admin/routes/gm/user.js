
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var session = require('../../models/session');

// ============================================================================

const C_filters = [
    ["--",      "--"],
    ["_id:s",   "玩家Id"],
    ["name:s",  "玩家名称"],
    ["area:n",  "区域Id"],
    ["svr:s",   "服名称"],
    ["sdk:s",   "sdk"],
    ["devid:s", "设备Id"],
];

const C_rpp = 22;

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

        res.render('gm/user', {
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
                name:   '玩家信息',
                header: tab.header,
                body:   tab.body,
            },

            paged: {
                page:  q.page,
                rec_n: tab.body ? tab.body.length : 0,
                rpp:   C_rpp,
            },

            scripts: [
                '/scripts/gm/user.js',
            ],
        });
    }

    // do it
    let cond = utils.filter({}, q.filters);
    var skip = (q.page - 1) * C_rpp;

    try {
        let db = dbpool.get('share');

        let docs = await db.collection('userinfo')
            .find(cond)
            .skip(skip)
            .limit(C_rpp)
            .toArray();

        var tab = {};

        tab.header = [
            "玩家Id", "名字", "等级", "Vip",
            "区域", "当前服", "Sdk",
            "设备型号", "设备Id", "操作系统", "系统版本",
            "创角IP",
            "创建时间",
            "封号截止",
        ];

        tab.body = docs.map(row => [
            row._id, row.name, row.lv, row.vip,
            row.area, row.svr, row.sdk,
            row.model, row.devid, row.os, row.osver,
            row.ip,
            row.cts.toString(),
            row.ban_ts.getTime() > Date.now() ? row.ban_ts.toString() : '-',
        ]);

        if (sess.priv['gm.tool.ban']) {
            tab.header.push('封号操作');
            tab.body.forEach(row => row.push(''));
        }

        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
