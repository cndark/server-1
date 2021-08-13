#!/usr/bin/env node

require('../../lib/ext');

var axios  = require('axios');

var config = require('../../../config.json');
var dbpool = require('../../lib/dbpool');
var utils  = require('../../lib/utils');

// ============================================================================

var args = process.argv.slice(2);
if (args.length < 1) {
    console.error('args: tmail-id');
    process.exit(1);
}

var tmail_id = args[0];

// ============================================================================

var db_client;
var db;

async function db_connect() {
    try {
        db_client = await dbpool.connect(config.common.db_share);
        db = db_client.db();
        return true;
    } catch(e) {
        db_close();
        console.error('connect to db-share failed:', e);
        return false;
    }
}

function db_close() {
    if (db_client) db_client.close();
}

// ============================================================================

var areas = {};

// ============================================================================

async function load_areas() {
    let docs = await db.collection('areas').find().toArray();
    docs.forEach(doc => {
        areas[doc._id] = doc.config;
    });
}

function process_tmail(doc) {
    /*
        arr format based on tp:
            user: [uid, ...]
            svr:  [{conf, ids}, ...]
    */

    let tp;
    let arr = [];

    if (doc.target.indexOf("u") >= 0) {
        tp = 'user';

        doc.target.split(",").forEach(v => {
            v = v.trim();
            if (!v) return;

            arr.push(v);
        });

    } else {
        tp = 'svr';

        if (doc.area == -1) {
            for (let i in areas) {
                let conf = areas[i];
                let ids = utils.parse_svrstr(conf, 'all');
                arr.push({conf, ids});
            }
        } else {
            let conf = areas[doc.area];
            if (!conf) throw `area not found: ${doc.area}`;

            let ids = utils.parse_svrstr(conf, doc.target);
            arr.push({conf, ids});
        }
    }

    return [tp, arr];
}

async function send_mail_svr(arr, mail) {
    for (let a of arr) {
        for (let id of a.ids) {
            let svr = `game${id}`;

            // get gs object
            let gs = a.conf.games[svr];
            if (!gs) {
                console.error('tmail: server NOT found:', tmail_id, a.conf.area, svr);
                continue;
            }

            // send
            try {
                await axios.post(`http://${gs.svc}/gm/service`, {
                    key:   "w.gmail",
                    title: mail.title,
                    text:  mail.text,
                    cond:  '',
                    res_k: mail.res_k,
                    res_v: mail.res_v,
                });
            } catch(e) {
                console.error('tmail: send error:', tmail_id, a.conf.area, svr, e.message);
            }
        }
    }
}

async function send_mail_user(arr, mail) {
    // lookup user svr
    let docs = await db.collection('userinfo').find({_id: {$in: arr}}).project({svr: 1}).toArray();

    // get conf
    let conf = areas[mail.area];
    if (!conf) {
        console.error('tmail: area NOT found:', mail.area);
        return;
    }

    for (let u of docs) {
        // get gs object
        let gs = conf.games[u.svr];
        if (!gs) {
            console.error('tmail: server NOT found:', tmail_id, conf.area, u.svr);
            continue;
        }

        // send
        try {
            await axios.post(`http://${gs.svc}/gm/service`, {
                key:   "w.pmail",
                plrid: u._id,
                title: mail.title,
                text:  mail.text,
                res_k: mail.res_k,
                res_v: mail.res_v,
            });
        } catch(e) {
            console.error('tmail: send error:', tmail_id, conf.area, u.svr, u._id, e.message);
        }
    }
}

async function update_status(st) {
    try {
        await db.collection('tmail').updateOne({_id: tmail_id}, {$set: {status: st}});
    } catch(e) {
        console.error('updating status failed:', e.message);
    }
}

// ============================================================================

(async function main() {
    // connect db
    if (! await db_connect()) return;

    try {
        // fetch tmail
        let r = await db.collection('tmail').findOneAndUpdate(
            {_id: tmail_id},
            {$set: {status: 'sending'}},
        );
        let doc = r.value;
        if (!doc) throw 'tmail not found';

        // load areas
        await load_areas();

        // preprocess tmail
        let [tp, arr] = process_tmail(doc);

        // send mails
        if (tp == 'svr') {
            await send_mail_svr(arr, doc);
        } else if (tp == 'user') {
            await send_mail_user(arr, doc);
        }

        // ok
        await update_status('ok');
        console.log('tmail successfully sent:', tmail_id);

    } catch(e) {
        await update_status('failed');
        console.error('tmail failed:', tmail_id, typeof e == 'string' ? e : e.message);
    }

    // close db
    db_close();

})().catch(console.error);
