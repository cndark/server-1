
var express = require('express');
var router  = express.Router();

var dbpool    = require('../lib/dbpool');
var area      = require('../models/area');
var settings  = require('../models/settings');
var cmd       = require('../models/cmd');
var svrstatus = require('../models/svrstatus');

// ============================================================================

// svrlists cache
var svrlists = {
    // aid: {
    //     name: {_id, id, tri_hwater, tri_closereg},
    // },
};

// ============================================================================

router.post('/', _A_(async (req, res) => {
    let q = req.body;

    let aid  = tonumber(q.area);
    let name = q.name;
    let n    = tonumber(q.n);

    try {
        // area
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // get cache. load cache if not found
        let lst = svrlists[aid];
        if (!lst) {
            lst = {};
            svrlists[aid] = lst;
        }

        let gs = lst[name];
        if (!gs) {
            let doc = await dbpool.exec(a.conf.common.db_center, async db => {
                return await db.collection('svrlist').findOne(
                    {_id: name},
                    {projection: {id: 1, tri_hwater: 1, tri_closereg: 1}},
                );
            });
            if (!doc) throw 'svr NOT found';

            gs = doc;
            lst[name] = gs;
        }

        // get setting object
        let sobj = settings.find(aid);
        if (!sobj) {
            res.json({}).end();
            return;
        }

        // check trigger hwater
        if (!gs.tri_hwater && sobj.closereg_hwater > 0 && n > sobj.closereg_hwater) {
            // update
            gs.tri_hwater = true;
            await dbpool.exec(a.conf.common.db_center, async db => {
                await db.collection('svrlist').updateOne(
                    {_id: name},
                    {$set: {tri_hwater: true}},
                );
            });

            // trigger it
            if (sobj.opennew_mode == 'full') {
                (async () => {
                    let opening = svrstatus.get_opening(aid);
                    if (opening) return;

                    try {
                        svrstatus.set_opening(aid, gs.id + 1);
                        await cmd.game_opennew(a, gs.id + 1);
                    } finally {
                        svrstatus.set_opening(aid, 0);
                        await svrstatus.fetch_last(aid);
                    }
                })().catch(console.error); // do NOT wait
            }
        }

        // check trigger closereg
        if (!gs.tri_closereg && sobj.closereg_limit > 0 && n > sobj.closereg_limit) {
            // update
            gs.tri_closereg = true;
            await dbpool.exec(a.conf.common.db_center, async db => {
                await db.collection('svrlist').updateOne(
                    {_id: name},
                    {$set: {tri_closereg: true}},
                );
            });

            // trigger it
            cmd.game_reg(a, 'close', name).catch(console.error); // do NOT wait
        }

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
