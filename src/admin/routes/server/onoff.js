
var express = require('express');
var router  = express.Router();

var axios = require('axios');

var config  = require('../../../config.json');
var utils   = require('../../lib/utils');
var area    = require('../../models/area');
var session = require('../../models/session');

// ============================================================================

router.get('/', (req, res) => {
    res.render('svr/onoff', {
        sess:  session.data(req),
        areas: area.idnames(),
    });
});

router.post('/cmd', _A_(async (req, res) => {
    let q  = req.body;
    let op = req.query.op;

    try {
        // op
        if (op != 'start' && op != 'stop') throw 'invalid op';

        // area
        q.area = tonumber(q.area);
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // target
        if (!q.target) throw 'no target';

        // send to agent
        let {data} = await axios.post(
            `http://${config.agent.ip}:${config.agent.port}/onoff/cmd`,
            {
                op:    op,
                area:  q.area,
                range: q.target,
            },
        );

        res.json(data).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

router.post('/refresh', _A_(async (req, res) => {
    let q = req.body;

    try {
        // get conf
        q.area = tonumber(q.area);
        let a = area.find(q.area);
        if (!a) throw 'area NOT found';

        // parse target
        if (!q.target) throw 'no target';
        let ids = utils.parse_svrstr(a.conf, q.target);
        if (ids.length == 0) throw 'no target found';

        // get status
        let {data} = await axios.get(
            `http://${config.agent.ip}:${config.agent.port}/onoff/status?area=${q.area}`,
        );
        if (data.err) throw data.err;

        let started  = [];
        let stopped  = [];
        let starting = [];
        let stopping = [];

        for (let id of ids) {
            let obj = data[id];
            if (obj) {
                switch (obj.st) {
                    case 'on':
                        started.push(id);
                        break;
                    case 'off':
                        stopped.push(id);
                        break;
                    case 'starting':
                        starting.push(id);
                        break;
                    case 'stopping':
                        stopping.push(id);
                        break;
                }
            }
        }

        res.json({started, stopped, starting, stopping}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
