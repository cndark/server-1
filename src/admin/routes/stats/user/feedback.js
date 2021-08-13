
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../../lib/dbpool');
var gtab    = require('../../../models/gtab');
var stats   = require('../stats');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    var q = stats.make_q(req.query);

    // render func
    var render = (tab) => {
        tab = tab || {};
        tab.name = '玩家反馈';
        stats.render(req, res, tab);
    }

    // ----------------------

    var pipeline = [
        {$match: {
            ts: {$gte: q.d_start, $lt: q.d_end},
        }},
     ];

    // ---------------------

    stats.filter(pipeline, q);

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('feedback').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ['设备id', '玩家id', '类型', '内容', '时间'];

        // body
        tab.body = docs.map(row => [
            row.devid,
            row.uid,
            row.tp,
            row.content,
            row.ts.toString(),
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
