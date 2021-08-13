
var express = require('express');
var router = express.Router();

var cluster = require('cluster');

var config = require('../../config.json');
var wblist = require('../models/wblist');

// ============================================================================

router.post('/check', function (req, res) {
    let m = req.ip.match(/[\d.]+$/);
    let ip = m ? m[0] : req.ip;

    res.json({
        ip:     wblist.type_ip(ip),
        device: wblist.type_device(req.body.deviceid),
    }).end();
});

router.get('/update', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({op: 'wblist.update'});

    res.end();
});

// ============================================================================

cluster.worker.on('message', msg => {
    if (msg.op == 'wblist.update') {
        wblist.load().catch(console.error);
    }
});

wblist.load().catch(console.error);

// ============================================================================

module.exports = router;
