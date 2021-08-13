
var express = require('express');
var router = express.Router();

var dbpool = require('../../lib/dbpool');
var token = require('../../models/token');

// ============================================================================

router.post('/', _A_(async (req, res) => {
    let q = req.body;
    if (!token.check_token(q)) {
        res.end("token err");
        return
    }
    delete q.token;

    if (q.content && q.content.length > 4096) {
        res.end("1")
        return
    }

    try {
        let db = dbpool.get('stats');

        await db.collection('feedback').insertOne({
            devid: q.deviceid || '',
            sdk: q.sdk || '',
            uid: q.uid || '',
            tp: q.tp || '',
            content: q.content || '',

            ts: new Date(),
        });

        res.end("0");

    } catch {
        res.end("1");
    }
}));

// ============================================================================

module.exports = router;
