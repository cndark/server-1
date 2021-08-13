
var express = require('express');
var router = express.Router();

var cluster = require('cluster');

var config = require('../../config.json');
var dbpool = require('../lib/dbpool');

var version = config.common.version;

// ============================================================================

router.post('/gate_add', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({ op: 'server.gate_add', body: req.body });

    res.end();
});

router.post('/gate_remove', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({ op: 'server.gate_remove', body: req.body });

    res.end();
});

router.post('/game_add', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({ op: 'server.game_add', body: req.body });

    res.end();
});

router.post('/game_remove', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;
    }

    cluster.worker.send({ op: 'server.game_remove', body: req.body });

    res.end();
});

router.get('/list/update', (req, res) => {
    // token is needed
    if (req.query.token != config.switcher.token) {
        res.status(404).end();
        return;

    }

    cluster.worker.send({ op: 'server.list_update' });

    res.end();
});

router.get('/list/gates', (req, res) => {
    res.json(servers.gates).end();
});

router.get('/list/svrready', (req, res) => {
    if (!req.query.version) {
        res.json({ code: 1 }).end();
        return
    }

    // check version
    let arr1 = req.query.version.split('.');
    let arr2 = version.split('.');

    if (arr1.length < 2 ||
        arr2.length < 2 ||
        arr1[0] != arr2[0] ||
        arr1[1] != arr2[1]) {

        res.json({ code: 1 }).end();
        return
    }

    if (Object.keys(servers.gates).length == 0) {
        res.json({ code: 2 }).end();
        return
    }

    res.json({ code: 0 }).end();
});

router.get('/list/games', (req, res) => {
    res.json(servers.games).end();
});

router.get('/list/single_game', (req, res) => {
    let name = req.query.name;

    try {
        if (!name) throw 'error params';

        let obj = servers.games[name];
        if (!obj) throw 'game not found';

        res.json(obj).end();
    } catch (e) {
        res.json({ err: typeof e == 'string' ? e : e.message }).end();
    }
});

// ============================================================================

var servers = {
    gates: {/*name: {ip, port, wsport, load}*/ },
    games: {/*name: {id, text, status, flag, closereg, on}*/ },
};

var merge_index = {
    // mname: [gameobj, ...]
}

var ttl = {
    gates: {/* name:  ts */ },
    games: {/* mname: ts */ },
}

// ============================================================================

function gate_add(q) {
    let name = q.name;
    let ip = q.ip;
    let port = Number(q.port);
    let wsport = Number(q.wsport);
    let load = Number(q.load);

    servers.gates[name] = { ip: ip, port: port, wsport: wsport, load: load };
    ttl.gates[name] = Date.now();
}

function gate_remove(names) {
    names.forEach(name => {
        delete servers.gates[name];
        delete ttl.gates[name];
    });
}

function game_add(q) {
    let name = q.name;

    let arr = merge_index[name];
    if (arr) {
        arr.forEach(svr => svr.on = true);
    }

    ttl.games[name] = Date.now();
}

function game_remove(names) {
    names.forEach(name => {
        let arr = merge_index[name];
        if (arr) {
            arr.forEach(svr => svr.on = false);
        }

        delete ttl.games[name];
    });
}

async function svrlist_load() {
    try {
        let db = dbpool.get("center");

        let docs = await db.collection('svrlist').find().sort({ id: 1 }).toArray();

        servers.games = {};
        merge_index = {};

        docs.forEach(r => {
            let id = r._id.match(/\d+$/)[0];
            let name = r._id;

            let svr = { id: id, text: r.text, status: r.status, flag: r.flag, closereg: r.closereg, on: false };

            // add svr
            servers.games[name] = svr;

            // add to merge index
            let e = merge_index[r.mname];
            if (!e) {
                e = [];
                merge_index[r.mname] = e;
            }
            e.push(svr);

            // check if svr is already on
            if (ttl.games[r.mname])
                svr.on = true;
        });
    } catch (e) {
        console.error("load svrlist failed:", e);
    }
}

// ============================================================================

cluster.worker.on('message', msg => {
    let q = msg.body;

    switch (msg.op) {
        case 'server.gate_add':
            gate_add(q);
            break;

        case 'server.gate_remove':
            gate_remove([q.name]);
            break;

        case 'server.game_add':
            game_add(q);
            break;

        case 'server.game_remove':
            game_remove([q.name]);
            break;

        case 'server.list_update':
            svrlist_load().catch(console.error);
            break;

        case 'reload':
            delete require.cache[require.resolve('../../config.json')];
            let c = require('../../config.json');
            version = c.common.version;
            break;

    }
});

// load server list on startup
svrlist_load().catch(console.error);

// clear timeout gates & games
setInterval(function () {
    let now = Date.now();
    let to_del;

    // check gates
    to_del = [];
    for (let name in ttl.gates) {
        if (now - ttl.gates[name] > 15000) {
            to_del.push(name);
        }
    }
    gate_remove(to_del);

    // check games
    to_del = [];
    for (let name in ttl.games) {
        if (now - ttl.games[name] > 15000) {
            to_del.push(name);
        }
    }
    game_remove(to_del);

}, 15000);

// ============================================================================

module.exports = router;
