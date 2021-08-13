
var express = require('express');
var router = express.Router();

var fx = require('./fx');

// ============================================================================

router.get('/', (req, res) => {
    fx.render(req, res, {
        name: '首冲',

        db: 'bill',
        coll: 'order',

        p0: [
            {
                $match: {
                    status: 'ok',
                    first: true,
                }
            },
            {
                $project: {
                    sdk: 1,
                    svr: 1,
                    ts: "$create_ts",
                    prod_id: 1,
                    csext: 1,
                    userid: 1,
                    lv: 1,
                }
            },
        ],
        filter0: null,

        p1: null,
        filter1: null,

        cols: [
            { type: 'group', key: 'prod_id', text: '充值项', dict: 'bill', selectable: true },
            { type: 'group', key: 'csext', text: '子项', dict: 'bill', selectable: true },
            { type: 'group', key: 'lv', text: '充值等级', selectable: true },
            { type: 'stat', key: 'plrn', text: '充值人数', op: { $addToSet: "$userid" }, tosize: true },
        ],

        hide: [],
    });
});

// ============================================================================

module.exports = router;
