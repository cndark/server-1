
var express = require('express');
var router  = express.Router();

var fx      = require('./fx');

// ============================================================================

router.get('/', (req, res) => {
    fx.render(req, res, {
        name:  '充值项',

        db:    'bill',
        coll:  'order',

        p0: [
            {$match: {
                status: 'ok',
            }},
            {$project: {
                sdk:     1,
                svr:     1,
                ts:      "$create_ts",
                prod_id: 1,
                csext:   1,
                userid:  1,
                amount:  {$divide: ["$amount", 100]},
            }},
        ],
        filter0: null,

        p1: null,
        filter1: null,

        cols: [
            {type: 'group', key: 'prod_id',  text: '充值项', dict: 'bill'},
            {type: 'group', key: 'csext',    text: '子项',   dict: 'bill'},
            {type: 'stat',  key: 'plrn',     text: '充值人数', op: {$addToSet: "$userid"}, tosize: true},
            {type: 'stat',  key: 'n',        text: '充值次数', op: {$sum: 1}},
            {type: 'stat',  key: 'amt',      text: '充值金额', op: {$sum: "$amount"}},
        ],

        hide: [],
    });
});

// ============================================================================

module.exports = router;
