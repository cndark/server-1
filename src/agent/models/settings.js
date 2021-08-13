
var config  = require('../../config.json');
var dbpool  = require('../lib/dbpool');

// ============================================================================

var settings = {
    // aid: {
    //     closereg_hwater, // number
    //     closereg_limit,  // number
    //     opennew_mode,    // manual|full|itv
    //     opennew_itv,     // in hours
    // },
};

// if new area is added, no settings for the new-area by default.
//  when configured, the new settings will be reloaded.
async function load() {
    let db = dbpool.get('share');

    let docs = await db.collection('settings').find().toArray();

    docs.forEach(doc => {
        settings[doc._id] = doc;
    });
}

function find(aid) {
    return settings[aid];
}

// ============================================================================

module.exports = {
    load: load,
    find: find,
}
