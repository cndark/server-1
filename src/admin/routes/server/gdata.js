
var express = require('express');
var router  = express.Router();

var gsvc    = require('../../models/gsvc');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('svr/gdata', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
});

router.post('/', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    try {
        // check
        if (!q.area || !q.target || !q.conf_k) throw 'error params';

        // add key
        q.key = 'w.conf';

        // run gsvc
        await gsvc.run(sess.user, q, tonumber(q.area), q.target);

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
