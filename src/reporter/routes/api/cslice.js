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

    try {
        let db = dbpool.get('stats');

        await db.collection('cslice').updateOne(
            {
                deviceid: q.deviceid || '',
                sdk: q.sdk || '',
            },
            {
                $set: {
                    step: q.step || '',
                    ts: new Date(),
                }
            },
            { upsert: true },
        );

        res.end("0");

    } catch {
        res.end("1");
    }
}));

// ============================================================================

module.exports = router;
