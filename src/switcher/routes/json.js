
var express = require('express');
var router  = express.Router();

var cluster = require('cluster');
var path    = require('path');
var process = require('process');

// ============================================================================

router.get('/', (req, res) => {
    let q = req.query;

    if (!q.f) {
        res.status(404).end();
        return;
    }

    try {
        let fn = path.join(process.env['WORK_DIR'], `update/conf/${q.f}`);
        let obj = require(fn);
        if (q.type == 'json') {
            res.json(obj).end();
        } else {
            res.end(Buffer.from(JSON.stringify(obj)).toString('base64'));
        }
    } catch (e) {
        res.status(404).end();
    }
});

// ============================================================================

cluster.worker.on('message', msg => {
    if (msg.op == 'reload') {
        let keys = [];
        for (k in require.cache) {
            if (k.match("update/conf/.+\.json$")) {
                keys.push(k);
            }
        }

        keys.forEach(v => {
            delete require.cache[v];
        });
    }
});

// ============================================================================

module.exports = router;
