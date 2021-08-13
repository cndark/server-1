
var md5 = require('md5');

var config = require('../../config.json');
var dbpool = require('../lib/dbpool');
var priv   = require('./priv');

// ============================================================================

const sys = {
    priv_pages:    priv.pages,
    priv_specials: priv.specials,
};

const empty = {user: '', rank: Number.MAX_SAFE_INTEGER, priv: {}, sys: sys};

var cache = {
    // 'user': {user, rank, priv, sys},
};

// ============================================================================

async function auth(req, user, pwd) {
    try {
        let db = dbpool.get('share');

        // check user and pwd
        let doc = await db.collection('adminuser').findOne({_id: user});
        if (!doc || doc.pwd != md5(pwd + config.admin.pwdfill)) throw 'failed';

        // create cache if not yet
        let d = cache[user];
        if (!d) {
            d = {
                user: user,
                rank: doc.rank,
                priv: priv.str2obj(doc.priv),
                sys:  sys,
            }
            cache[user] = d;
        }

        // ok
        req.session.user = user;

        return true;
    } catch {
        return false;
    }
}

function data(req) {
    return cache[req.session.user] || empty;
}

function destroy(req) {
    delete req.session.user;
    req.session.destroy();
}

function update_priv(user, pstr) {
    let d = cache[user];
    if (!d) return;

    d.priv = priv.str2obj(pstr);
}

function delete_cache(user) {
    delete cache[user];
}

// ============================================================================

module.exports = {
    auth:         auth,
    data:         data,
    destroy:      destroy,
    update_priv:  update_priv,
    delete_cache: delete_cache,
}
