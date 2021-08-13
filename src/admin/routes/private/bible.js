
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    res.render('private/bible', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
}));

router.post('/dyncode/set', _A_(async (req, res) => {
    let q = req.body;

    let aid  = tonumber(q.area);
    let code = q.code;

    try {
        // area
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // code
        if (!code) throw 'error params';
        if (code.length < 32) throw 'code too short';

        // set
        let {data} = await axios.post(
            `http://${a.conf.auth.ip}:${a.conf.auth.port}/sdk/dyncode/update`,
            {code},
        );
        if (data.err) throw err;

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
