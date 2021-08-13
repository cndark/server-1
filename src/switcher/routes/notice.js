
var express = require('express');
var router = express.Router();

var cluster = require('cluster');

var config = require('../../config.json');
var dbpool = require('../lib/dbpool');

// ============================================================================

router.get('/', (req, res) => {
    let q = req.query;

    let e = notice[q.lang]
    let now = new Date().unix();

    res.json(e && now >= e.d_start && now <= e.d_end ? e : {}).end();
});

router.get('/update', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({ op: 'notice.update' });

    res.end();
});

// ============================================================================

var notice = {
    // 'lang': {lang, d_start, d_end, title, content},
};

async function load_notice() {
    try {
        let db = dbpool.get("center");

        let docs = await db.collection('notice').find().toArray();

        notice = {};

        docs.forEach(doc => {
            doc.d_start = Math.floor(doc.d_start.getTime() / 1000);
            doc.d_end = Math.floor(doc.d_end.getTime() / 1000);

            notice[doc._id] = doc;
        });
    } catch (e) {
        console.error('load notice failed:', e);
    }
}

cluster.worker.on('message', msg => {
    if (msg.op == 'notice.update') {
        load_notice().catch(console.error);
    }
});

// load notice on startup
load_notice().catch(console.error);

// ============================================================================

module.exports = router;
