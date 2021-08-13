
var tab_sdk = require('../../gamedata/data/sdk.json');

// ============================================================================

module.exports = {

    sdk: (function () {
        var m = {};
        tab_sdk.forEach(v => {
            m[v.name] = v;
        });
        return m;
    })(),

};
