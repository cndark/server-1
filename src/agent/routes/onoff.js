
var express = require('express');
var router  = express.Router();

var area      = require('../models/area');
var cmd       = require('../models/cmd');
var svrstatus = require('../models/svrstatus');

// ============================================================================

router.get('/status', (req, res) => {
    let q = req.query;

    try {
        let lst = svrstatus.onoff[q.area];
        if (!lst) throw 'area NOT found';

        res.json(lst).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
});

router.post('/cmd', _A_(async (req, res) => {
    let q = req.body;

    try {
        // op
        if (q.op != 'start' && q.op != 'stop') throw 'invalid op';

        // area
        let aid = tonumber(q.area);
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // check available servers
        if (!a.game_id_min) throw 'no available servers';

        // range
        if (!q.range) throw 'no range';

        let fromid, toid;

        if (q.range == 'all') {
            fromid = a.game_id_min;
            toid   = a.game_id_max;
        } else {
            let arr = q.range.split('-');
            if (arr.length == 1) {
                fromid = tonumber(arr[0]);
                toid   = tonumber(arr[0]);
            } else {
                fromid = tonumber(arr[0]);
                toid   = tonumber(arr[1]);
            }
            if (fromid < a.game_id_min) fromid = a.game_id_min;
            if (toid   > a.game_id_max) toid   = a.game_id_max;

            if (fromid > toid) throw 'invalid range';
        }

        // check status
        let lst = svrstatus.onoff[aid];
        if (!lst) throw 'service not available';

        for (let i = fromid; i <= toid; i++) {
            let e = lst[i];
            if (!e) continue;

            if (q.op == 'start' && e.st != 'off' || q.op == 'stop' && e.st != 'on') {
                throw `conflict: game${i} is ${e.st}`;
            }
        }

        // update status
        for (let i = fromid; i <= toid; i++) {
            let e = lst[i];
            if (!e) continue;

            if (q.op == 'start') {
                e.st = 'starting';
            } else if (q.op == 'stop') {
                e.st = 'stopping';
            }
        }

        // execute
        cmd.game_onoff(a, q.op, fromid, toid).catch(console.error); // do NOT wait

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
