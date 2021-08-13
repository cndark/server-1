
var express = require('express');
var router = express.Router();

var fx = require('./fx');

// ============================================================================

var A = [
    ["1", "英雄召唤"],
    ["2", "英雄合成"],
    ["3", "英雄升级"],
    ["4", "推图"],
    ["5", "在线奖励"],
]

// ============================================================================

router.get('/', (req, res) => {
    let q = req.query;

    fx.render(req, res, {
        name: '引导流失',

        db: 'stats',
        coll: 'tutorial',

        p0: [
            {
                $project: {
                    _id: 0,
                    area: 1,
                    svr: 1,
                    sdk: 1,
                    tp: `$tut.${q.tut_tp}`,
                }
            },
            {
                $match: {
                    tp: { $gte: 0 },
                }
            },
        ],
        filter0: null,

        p1: null,
        filter1: null,

        form: [
            { type: "select", name: "tut_tp", src: A, value: q.tut_tp, desc: "引导类型" },
        ],

        cols: [
            { type: 'group', key: "tp", text: '步数' },
            { type: 'stat', key: 'n', text: '人数', op: { $sum: 1 } },
        ],

        hide: ['dr'],
    });
});

// ============================================================================

module.exports = router;
