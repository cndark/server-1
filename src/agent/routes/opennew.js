
var express = require('express');
var router  = express.Router();

var area      = require('../models/area');
var settings  = require('../models/settings');
var cmd       = require('../models/cmd');
var svrstatus = require('../models/svrstatus');

// ============================================================================

router.get('/status', _A_(async (req, res) => {
    let q = req.query;

    try {
        // area
        let aid = tonumber(q.area);
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // get opening status
        let opening = svrstatus.get_opening(aid);

        // get last status
        let last = await svrstatus.get_last(aid);

        // execute
        let str = await cmd.game_next_wait_id(a, true);
        let wait_ids = str.split(/\s+/).map(v => tonumber(v)).filter(v => v > 0);

        // ok
        let r = {opening, wait_ids, last: {}};
        if (last.id) {
            r.last.id = last.id;
            r.last.ts = last.ts.toString();
        }

        res.json(r).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

router.post('/cmd', _A_(async (req, res) => {
    let q = req.body;

    try {
        // area
        let aid = tonumber(q.area);
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // id
        let id = tonumber(q.id);
        if (a.conf.games[`game${id}`]) throw `game${id} is already open`;

        // get opening status
        let opening = svrstatus.get_opening(aid);
        if (opening) throw `game${opening} is opening`;

        // check settings
        let sobj = settings.find(aid);
        if (sobj && sobj.opennew_mode != 'manual') throw `manually open new is disabled`;

        // check id/next-id consistency
        let wait_id = await cmd.game_next_wait_id(a);
        wait_id = tonumber(wait_id);
        if (id != wait_id) throw `next id should be ${wait_id}`;

        // execute
        (async () => {
            try {
                svrstatus.set_opening(aid, id);
                await cmd.game_opennew(a, id);
            } finally {
                svrstatus.set_opening(aid, 0);
                await svrstatus.fetch_last(aid);
            }
        })().catch(console.error); // do NOT wait

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================
// open new by itv

async function check_itv_opennew(aid) {
    // area
    let a = area.find(aid);
    if (!a) return;

    // check settings
    let sobj = settings.find(aid);
    if (!sobj) return;
    if (sobj.opennew_mode != 'itv' || sobj.opennew_itv <= 0) return;

    // get opening
    let opening = svrstatus.get_opening(aid);
    if (opening) return;

    // get last
    let last = await svrstatus.get_last(aid);
    if (!last.ts) return;

    // check itv
    if (Date.now() - last.ts.getTime() < sobj.opennew_itv * 3600 * 1000) return;

    // open new
    (async () => {
        try {
            svrstatus.set_opening(aid, last.id + 1);
            await cmd.game_opennew(a, last.id + 1);
        } finally {
            svrstatus.set_opening(aid, 0);
            await svrstatus.fetch_last(aid);
        }
    })().catch(console.error); // no wait
}

async function check_itv_opennew_all() {
    for (let aid of area.ids()) {
        await check_itv_opennew(aid);
    }
}

setInterval(() => {
    check_itv_opennew_all().catch(console.error);
}, 15 * 1000);

// ============================================================================

module.exports = router;
