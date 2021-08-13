
var mongodb = require('mongodb');

// ============================================================================

var pool = {};

// ============================================================================

async function connect(cnnstr, size) {
    return await mongodb.connect(cnnstr, {
        poolSize:           size || 1,
        socketTimeoutMS:    0,
        useUnifiedTopology: true,
    });
}

async function init(pools) {
    /*
        pools: [
            [key, cnnstr, size],
            ...
        ]
    */

    await Promise.all(pools.map(async p => {
        let key    = p[0];
        let cnnstr = p[1];
        let size   = p[2];

        pool[key] = await connect(cnnstr, size);
    }));
}

function get(key) {
    let c = pool[key];
    return c ? c.db() : null;
}

async function exec(cnnstr, f) {
    let c;
    try {
        c = await connect(cnnstr);
        return await f(c.db());
    } finally {
        if (c) c.close();
    }
}

async function create_index(arr) {
    /*
        [
            [coll, name, cols, uk],
            ...
        ]
    */

    await Promise.limit(5, arr, async e => {
        let coll = e[0];
        let name = e[1];
        let cols = e[2];
        let uk   = e[3];

        await coll.createIndex(cols, {name: name, unique: uk});
    });
}

// ============================================================================

module.exports = {
    connect:      connect,
    init:         init,
    get:          get,
    exec:         exec,
    create_index: create_index,
};
