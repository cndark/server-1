
var express = require('express');
var router  = express.Router();

var config  = require('../../../config.json');
var dbpool  = require('../../lib/dbpool');
var gsvc    = require('../../models/gsvc');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('gm/tool', {
        sess:    session.data(req),
        devmode: config.common.dev_mode,
    });
});

router.post('/', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    // check gm priv
    if (q.key && q.key.startsWith("w.")) {
        do {
            if (sess.priv['gm.w']) break;
            if (q.key == "w.ban" && sess.priv['gm.tool.ban']) break;

            res.json({err: "access denied"}).end();
            return;
        } while (0);
    }

    // check plrid
    if (!q.plrid && !q.plrname) {
        res.json({err: "invalid request"}).end();
        return;
    }

    // find player svr and run gsvc
    try {
        let db = dbpool.get('share');

        let doc = await db.collection('userinfo').findOne(
            {$or: [{_id: q.plrid}, {name: q.plrname}]},
            {projection: {name: 1, area: 1, svr: 1}},
        );
        if (!doc) throw 'player not found';

        q.plrid   = doc._id;
        q.plrname = doc.name;

        let r = await gsvc.run(sess.user, q, doc.area, doc.svr);
        res.json(r).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
