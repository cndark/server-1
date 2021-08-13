
var express = require('express');
var router = express.Router();

var dbpool = require('../../../lib/dbpool');
var gtab = require('../../../models/gtab');
var stats = require('../stats');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    var q = stats.make_q(req.query);

    // render func
    var render = (tab) => {
        tab = tab || {};
        tab.name = '付费排行榜';
        stats.render(req, res, tab);
    }

    // ----------------------

    var pipeline = [
        {
            $match: {
                bts: { $gte: q.d_start, $lt: q.d_end },
                amt: { $gt: 0 },
            }
        },
        {
            $group: {
                _id: "$uid",
                name: { $last: "$name" },
                amt: { $sum: "$amt" },
                ccy: { $last: "$ccy" },
                n: { $sum: "$n" },
                sdk: { $last: "$sdk" },
                svr: { $last: "$svr" },
                cts: { $last: "$cts" },
                bts: { $max: "$bts" },
            }
        },
        { $sort: { amt: -1 } },
        { $limit: 500 }
    ];

    // ---------------------

    stats.filter(pipeline, q);

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('bill').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ['玩家Id', '玩家名字', '金额', '币种', '充值次数', 'Sdk', '服务器', '玩家创建时间', '最后充值时间'];

        // body
        tab.body = docs.map(row => [
            row._id,
            row.name,
            row.amt + ` (${gtab.toCNY(row.amt).toFixed(2)} CNY)`,
            row.ccy,
            row.n,
            row.sdk,
            row.svr,
            row.cts.toString(),
            row.bts.toString(),
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
