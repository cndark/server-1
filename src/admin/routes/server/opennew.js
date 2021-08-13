
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var config  = require('../../../config.json');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('svr/opennew', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
});

router.post('/refresh', _A_(async (req, res) => {
    let q = req.body;

    try {
        // area
        let aid = tonumber(q.area);
        let a = area.find(aid);
        if (!a) throw 'area NOT found';

        // get info
        let {data} = await axios.get(
            `http://${config.agent.ip}:${config.agent.port}/opennew/status?area=${aid}`,
        );
        if (data.err) throw data.err;

        // result
        res.json(data).end();
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
        if (id == 0) throw 'error params';

        // send to agent
        let {data} = await axios.post(
            `http://${config.agent.ip}:${config.agent.port}/opennew/cmd`,
            {
                area: q.area,
                id:   q.id,
            },
        );

        res.json(data).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
