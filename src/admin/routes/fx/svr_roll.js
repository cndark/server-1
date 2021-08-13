
var express = require('express');
var router  = express.Router();

var fx      = require('./fx');

// ============================================================================

router.get('/', (req, res) => {
    fx.render(req, res, {
        name:  '滚服占比',

        db:    'share',
        coll:  'userinfo',

        p0: null,
        filter0: null,

        p1: [
            {$group: {
                _id:   {sdk: "$sdk", authid: "$authid"},
                first: {$min:  "$svr0"},
                all:   {$push: "$svr0"},
            }},
            {$unwind: "$all"},
            {$project: {
                svr:   "$all",
                roll:  {$cond: [{$eq: ["$first", "$all"]}, 0, 1]},
            }},
        ],
        filter1: null,

        cols: [
            {type: 'group', key: 'svr',       text: '区服'},
            {type: 'stat',  key: 'totalcnt',  text: '总人数',   op: {$sum: 1}},
            {type: 'stat',  key: 'rollcnt',   text: '滚服人数', op: {$sum: "$roll"}},
            {type: 'calc',  key: 'rollratio', text: '滚服占比', formula: '($rollcnt / $totalcnt * 100).toFixed(1) + "%"'},
        ],

        hide: ['fkey', 'dr'],
    });
});

// ============================================================================

module.exports = router;
