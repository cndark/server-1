
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

const C_filters = [
    ["--",       "--"],
    ["mname:s",  "合服名"],
    ["status:n", "状态"],
    ["flag:s",   "标志"],
];

const C_rpp = 30;

const C_status = [
    ['--', '--'],
    [1,    '维护 (1)'],
    [2,    '繁忙 (2)'],
    [3,    '火爆 (3)'],
];

const C_flag = [
    ['--',   '--'],
    ['',     '无'],
    ['new',  '新服 (new)'],
    ['good', '推荐 (good)'],
    ['test', '内部测试 (test)'],
];

// ============================================================================

function make_q(q) {
    // area
    q.area = tonumber(q.area);

    // from ~ to
    if (q.fromid) q.fromid = tonumber(q.fromid);
    if (q.toid)   q.toid   = tonumber(q.toid);

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
    let sess = session.data(req);
    let q    = make_q(req.query);

    // render func
    let render = (tab) => {
        tab = tab || {};

        res.render('svr/status', {
            sess: sess,

            src_status: C_status,
            src_flag:   C_flag,

            form: [
                {type: "select", name: "area",   src: area.idnames(),   value: q.area,    desc: "区域"},
                {type: "input",  name: "fromid", placeholder: "开始Id", value: q.fromid},
                {type: "input",  name: "toid",   placeholder: "结束Id", value: q.toid},
                {type: "select", name: "fkey",   src: C_filters,        value: q.fkey[0], icon: "filter"},
                {type: "input",  name: "fval",   placeholder: "值",     value: q.fval[0], style: "width: 120px"},
            ],

            tab: {
                name:   '列表查询',
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
    if (q.fromid) { cond.id = cond.id || {}; cond.id['$gte'] = q.fromid;}
    if (q.toid)   { cond.id = cond.id || {}; cond.id['$lte'] = q.toid;}

    let skip = (q.page - 1) * C_rpp;

    try {
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        let docs = await dbpool.exec(a.conf.common.db_center, async db => {
            return await db.collection('svrlist')
                .find(cond)
                .sort({id: 1})
                .skip(skip)
                .limit(C_rpp)
                .toArray();
        });

        let tab = {};

        tab.header = ["原服名", "合服名", "显示名", "开服时间", "状态", "标志", "关闭注册"];

        tab.body = docs.map(row => [
            row._id,
            row.mname,
            row.text,
            row.ts.toString(),
            C_status.find(v => v[0] == row.status)[1],
            C_flag.find(v => v[0] == row.flag)[1],
            row.closereg ? 'Y' : '-',
        ]);

        render(tab);
    } catch {
        render();
    }
}));

router.post('/', _A_(async (req, res) => {
    let q = req.body;

    // check params
    if (!q.area || !q.target) {
        res.json({err: 'error params'}).end();
        return;
    }

    // do it
    try {
        // area
        let aid = tonumber(q.area);
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // target
        let cond = {};
        let single = false;

        if (q.target == 'all') {
            // do nothing
        } else {
            let arr = q.target.split('-');
            if (arr.length == 1) {
                single = true;
                cond.id = tonumber(arr[0]);
            } else {
                cond.id = {$gte: tonumber(arr[0]), $lte: tonumber(arr[1])};
            }
        }

        // set-doc
        let doc = {};

        // text
        if (q.text) {
            if (!single) throw 'text NOT allowed for multi-records';
            doc.text = q.text;
        }

        // status
        if (q.status != '--') {
            q.status = tonumber(q.status);
            if (C_status.find(v => v[0] == q.status)) {
                doc.status = q.status;
            }
        }

        // flag
        if (q.flag != '--' && C_flag.find(v => v[0] == q.flag)) {
            doc.flag = q.flag;
        }

        // closereg
        if (q.closereg != '--') {
            q.closereg = tonumber(q.closereg);
            doc.closereg = q.closereg != 0;
        }

        // check doc
        if (Object.keys(doc).length == 0) throw 'no update specified';

        // update
        await dbpool.exec(a.conf.common.db_center, async db => {
            if (single) {
                await db.collection('svrlist').updateOne(cond, {$set: doc});
            } else {
                await db.collection('svrlist').updateMany(cond, {$set: doc});
            }
        });

        // notify switcher
        await axios.get(
            `http://${a.conf.switcher.ip}:${a.conf.switcher.port}/server/list/update?token=${a.conf.switcher.token}`,
        );

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
