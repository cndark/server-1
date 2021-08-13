
var dbpool = require('../lib/dbpool');

// ============================================================================

var _ids;
var _idnames;
var _areas;

// ============================================================================

async function load() {
    try {
        let db = dbpool.get('share');

        let docs = await db.collection('areas').find().toArray();

        _ids     = [];
        _idnames = [];
        _areas   = {};

        docs.forEach(doc => {
            // ids & idnames
            _ids.push(doc._id);
            _idnames.push([doc._id, doc.name]);

            // area object
            let obj = {
                dir:  doc.dir,
                conf: doc.config,
            };
            _areas[doc._id] = obj;

            let ids = Object.keys(doc.config.games)
                .map(v => Number(v.match(/\d+$/)[0]))
                .sort((a, b) => a - b);

            if (ids.length > 0) {
                obj.game_id_min = ids[0];
                obj.game_id_max = ids[ids.length - 1];
            }
        });
    } catch(e) {
        console.error('loading areas failed:', e.message);
    }
}

function ids() {
    return _ids;
}

function idnames() {
    return _idnames;
}

function find(id) {
    return _areas[id];
}

// ============================================================================

module.exports = {
    load:    load,
    ids:     ids,
    idnames: idnames,
    find:    find,
}
