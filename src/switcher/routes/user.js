
var express = require('express');
var router = express.Router();

var dbpool = require('../lib/dbpool');

// ============================================================================

router.get('/chars', _A_(async (req, res) => {
    let q = req.query;

    if (!q.auth_id || !q.sdk) {
        res.status(404).end();
        return;
    }

    try {
        let db = dbpool.get("share");

        let docs = await db.collection('userinfo')
            .find({ authid: q.auth_id, sdk: q.sdk })
            .project({ _id: 0, svr0: 1, name: 1, lv: 1, vip: 1 })
            .toArray();

        res.json(docs).end();
    } catch {
        res.json([]).end();
    }
}));

router.get('/lastsvr', _A_(async (req, res) => {
    let q = req.query;

    if (!q.auth_id || !q.sdk) {
        res.status(404).end();
        return;
    }

    try {
        let db = dbpool.get("center");

        let doc = await db.collection('acctinfo').findOne({ _id: `${q.sdk}-${q.auth_id}` });

        let r = { svr0: doc ? doc.lastsvr : '' };

        res.json(r).end();
    } catch {
        res.status(500).end();
    }
}));

// ============================================================================

module.exports = router;
