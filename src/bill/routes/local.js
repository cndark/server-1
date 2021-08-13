
var express = require('express');
var router  = express.Router();

var dbpool = require('../lib/dbpool');

// ============================================================================

router.post('/', _A_(async (req, res) => {
    var q = req.body;

    var h = handlers[q.op];
    if (h) {
        let r = await h(q);
        res.json(r).end();
    } else {
        res.status(404).end();
    }
}));

// ============================================================================

var handlers = {
    "genorder": async q => {
        try {
            let db = dbpool.get('bill');

            let doc = {
                _id:      `${Date.now()}-${Math.floor(Math.random() * 100000)}`,
                prod_id:  q.prod_id,
                csext:    q.csext,
                ts:       new Date().unix(),
            };

            await db.collection('genorder').insertOne(doc);

            return {
                op:       "genorder",
                orderid:  doc._id,
                prod_id:  doc.prod_id,
                csext:    doc.csext,
                ts:       doc.ts,
            };
        } catch {
            return {err: 'failed'};
        }
    },
}

// ============================================================================

module.exports = router;
