
var express = require('express');
var router  = express.Router();

var child_process = require('child_process');
var path          = require('path');

var dbpool  = require('../../lib/dbpool');
var utils   = require('../../lib/utils');
var session = require('../../models/session');
var area    = require('../../models/area');

// ============================================================================

const C_filters = [
    ["--",          "--"],
    ["title:re",    "标题"],
    ["apply:s",     "申请人"],
    ["audit:s",     "审核人"],
    ["status:s",    "状态"],
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

    // dr
    if (!q.dr) {
        let now = new Date();
        q.dr = `${now.addDay(-15).toDateString()} ~ ${now.addDay(90).toDateString()}`;
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

        res.render('gm/tmail', {
            sess:  sess,
            areas: area.idnames(),

            form: [
                `<button id='btn_add' class='btn btn-primary' type='button'><i class='fa fa-plus'></i></button>`,
                {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[0], icon: "filter"},
                {type: "input",     name: "fval", placeholder: "值", value: q.fval[0], style: "width: 120px"},
                {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[1], icon: "filter"},
                {type: "input",     name: "fval", placeholder: "值", value: q.fval[1], style: "width: 120px"},
                {type: "daterange", name: "dr",                      value: q.dr,      icon: "calendar-alt"},
            ],

            tab: {
                name:   '定时邮件',
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
    cond.send_ts = {$gte: q.d_start, $lt: q.d_end};

    let skip = (q.page - 1) * C_rpp;

    try {
        let db = dbpool.get("share");

        let docs = await db.collection('tmail')
            .find(cond)
            .project({text: 0, res_k: 0, res_v: 0})
            .sort({send_ts: -1})
            .skip(skip)
            .limit(C_rpp)
            .toArray();

            let tab = {};

        tab.header = ["发送时间", "目标", "标题", "状态", "申请人", "审核人", "详细", "删除"];

        tab.body = docs.map(row => [
            row._id,
            row.send_now ? '立即发送' : row.send_ts.toString(),
            row.target.substring(0, 50),
            row.title,
            row.status,
            row.apply,
            row.audit,
            "",
            "",
        ]);

        render(tab);

    } catch {
        render();
    }
}));

router.post('/detail', _A_(async (req, res) => {
    let q = req.body;

    // check args
    if (!q._id) {
        res.json({err: "invalid args"}).end();
        return;
    }

    // query
    try {
        let db = dbpool.get("share");

        let doc = await db.collection('tmail').findOne({_id: q._id});
        if (!doc) throw 'not found';

        doc.send_ts = doc.send_now ? '' : doc.send_ts.toString();

        res.json(doc).end();
    } catch {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

router.post('/update', _A_(async (req, res) => {
    let sess  = session.data(req);
    let q     = req.body;
    let isnew = req.query.isnew == "true";

    // check args
    if (!q._id || !q.area || !q.title || !q.text || q.area != '-1' && !q.target) {
        res.json({err: "invalid args"}).end();
        return;
    }

    // area & target
    q.area = tonumber(q.area);
    if (q.area == -1) q.target = '';

    // if send_ts is specified, it MUST NOT be a history time
    if (q.send_ts) {
        q.send_ts = Date.fromString(q.send_ts);
        if (q.send_ts.getTime() <= Date.now()) {
            res.json({err: "send time MUST be future time"}).end();
            return;
        }
    }

    // check res
    if (!(q.res_k instanceof Array) || !(q.res_v instanceof Array)) {
        res.json({err: "invalid rewards"}).end();
        return;
    }

    // update
    try {
        let db = dbpool.get("share");

        let now = new Date();
        let doc = {
            _id:       q._id,
            send_ts:   q.send_ts || now,
            send_now:  q.send_ts ? false : true,
            area:      q.area,
            target:    q.target,
            title:     q.title,
            text:      q.text,
            res_k:     q.res_k,
            res_v:     q.res_v,
            apply:     sess.user,
            audit:     '-',
            status:    'wait',
        };

        if (isnew) {
            // insert
            try {
                await db.collection('tmail').insertOne(doc);
            } catch {
                throw 'failed';
            }
        } else {
            // update
            try {
                let r = await db.collection('tmail').replaceOne(
                    {_id: q._id, status: 'wait'},
                    doc,
                );
                if (r.matchedCount == 0) throw 'not found or not allowed';
            } catch {
                throw 'failed';
            }
        }

        res.json({}).end();
    } catch(e) {
        res.json({err: e}).end();
    }
}));

router.post('/delete', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    // check args
    if (!q._id) {
        res.json({err: "invalid args"}).end();
        return;
    }

    // delete
    try {
        let db    = dbpool.get("share");
        let audit = sess.priv['gm.mail.audit'];
        let now   = new Date();

        let doc_or = [{status: 'wait'}];
        if (audit) {
            doc_or.push({status: 'audit', send_ts: {$gt: now.add('M', 3)}});
        }

        let r = await db.collection('tmail').deleteOne({_id: q._id, $or: doc_or});
        if (r.deletedCount == 0) throw 'not found or not allowed';

        // update sched
        sched_update('del', {_id: q._id});

        res.json({}).end();
    } catch(e) {
        res.json({err: e.message}).end();
    }
}));

router.post('/audit', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    // check priv
    if (!sess.priv['gm.mail.audit']) {
        res.json({err: "no priv"}).end();
        return;
    }

    // check args
    if (!q._id) {
        res.json({err: "invalid args"}).end();
        return;
    }

    // audit
    try {
        let db = dbpool.get("share");

        let r = await db.collection('tmail').findOneAndUpdate(
            {_id: q._id, status: 'wait'},
            {$set: {status: 'audit', audit: sess.user}},
        );
        if (!r.value) throw 'not found or not allowed';

        // update sched
        sched_update('add', {
            _id:      q._id,
            send_ts:  r.value.send_ts,
            send_now: r.value.send_now,
        });

        res.json({}).end();
    } catch(e) {
        res.json({err: e.message}).end();
    }
}));

// ============================================================================

var sched_mails = [];
var sched_h     = {};

// ============================================================================

async function sched_load() {
    try {
        let db = dbpool.get("share");

        let docs = await db.collection('tmail')
            .find({status: 'audit'})
            .project({send_ts: 1, send_now: 1})
            .toArray();

        sched_mails = docs;

        sched_sort();
        sched_start();
    } catch(e) {
        console.error('loading tmail failed:', e.message)
    }
}

function sched_sort() {
    sched_mails.sort((a, b) => b.send_ts.unix() - a.send_ts.unix());
}

function sched_start() {
    let L = sched_mails.length;
    if (L == 0) return;

    let obj = sched_mails[L - 1];

    setRunAt(sched_h, obj.send_ts, () => {
        sched_mails.pop();

        child_process.spawn(path.join(__dirname, './tmail_sender.js'), [obj._id], {
            detached: true,
            stdio:    'inherit',
        }).unref();

        sched_start();
    });
}

function sched_stop() {
    clearRunAt(sched_h);
}

function sched_update(op, obj) {
    sched_stop();

    if (op == 'add') {
        sched_mails.push(obj);
    } else if (op == 'del') {
        let L = sched_mails.length;
        for (let i = 0; i < L; i++) {
            if (sched_mails[i]._id == obj._id) {
                sched_mails[i] = sched_mails[L - 1];
                sched_mails.pop();
                break;
            }
        }
    }

    sched_sort();
    sched_start();
}

// ============================================================================

// initial load
sched_load().catch(console.error);

// ============================================================================

module.exports = router;
