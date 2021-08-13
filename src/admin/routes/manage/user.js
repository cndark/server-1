
var express = require('express');
var router  = express.Router();

var md5 = require('md5');

var config  = require('../../../config.json');
var dbpool  = require('../../lib/dbpool');
var edit    = require('../../lib/edit');
var session = require('../../models/session');

// ============================================================================

const PWD_MASK = '••••••';

// ============================================================================

router.all('/', function (req, res) {
    let sess = session.data(req);

    // edit show
    edit.show(req, res, {
        view:    "manage/user",
        db:      "share",
        coll:    "adminuser",
        rpp:     13,
        scripts: ['/scripts/manage/user.js'],

        on_load: (docs) => {
            docs.forEach(doc => {
                doc.pwd = PWD_MASK;
                delete doc.priv;
            });
        },

        on_before_save: (docs) => {
            let defers = [];

            docs.forEach(doc => {
                let op = doc.__edit_op__;

                // check field types
                if (typeof doc._id != 'string') {
                    throw 'dont cheat';
                }
                if  (op != 'del') {
                    if (typeof doc.rank != 'number' || typeof doc.pwd != 'string') {
                        throw 'dont cheat';
                    }
                }

                // check rank
                if (op == 'add') {
                    if (doc.rank <= sess.rank) throw 'no right to operate high-rank';
                } else if (op == 'up') {
                    if (doc._id != sess.user) {
                        if (doc.rank <= sess.rank) throw 'no right to operate high-rank';
                        doc.__edit_cond__ = {rank: {$gt: sess.rank}};
                    }
                    // can't modify rank
                    delete doc.rank;
                } else if (op == 'del') {
                    doc.__edit_cond__ = {rank: {$gt: sess.rank}};
                    defers.push(() => session.delete_cache(doc._id));
                }

                // check pwd
                if (op != 'del') {
                    if (!doc.pwd) throw 'password needed';

                    if (doc.pwd.indexOf(PWD_MASK[0]) >= 0) {
                        delete doc.pwd;
                    } else {
                        doc.pwd = md5(doc.pwd + config.admin.pwdfill);
                        defers.push(() => session.delete_cache(doc._id));
                    }
                }

                // no priv
                delete doc.priv;
            });

            // execute defers
            defers.forEach(f => f());
        },

        cols: [
            {col: "_id",   name: "用户名",         dtype: "string", ctrl: "input", width: '160px'},
            {col: "pwd",   name: "密码",           dtype: "string", ctrl: "input", width: '160px'},
            {col: "rank",  name: "级别 (越小越高)", dtype: "number", ctrl: "input", width: '160px', immutable: true},
            {col: "memo",  name: "备注",           dtype: "string", ctrl: "input", width: '500px', optional: true},
            {col: "grant", name: "授权",           dtype: "string", ctrl: "html",  src: ''},
        ],

        filters: [
            "_id:s",
        ],
    });
});

router.post('/priv_get', _A_(async (req, res) => {
    let q = req.body;

    try {
        // check
        let _id = q._id;
        if (!_id) throw 'error params';

        let db = dbpool.get('share');

        let doc = await db.collection('adminuser').findOne({_id: _id});
        if (!doc) throw 'not found';

        let parr = doc.priv || '';
        parr = parr.split(',').map(v => v.trim()).filter(v => v != '');

        res.json(parr).end();
    } catch {
        res.json({err: 'failed'}).end();
    }
}));

router.post('/priv_set', _A_(async (req, res) => {
    let sess = session.data(req);
    let q    = req.body;

    try {
        // check params
        let _id  = q._id;
        let pstr = q.pstr || '';
        if (!_id) throw 'error params';

        // check priv
        let parr = pstr.split(',').map(v => v.trim()).filter(v => v != '');
        for (let key of parr) {
            if (!sess.priv[key]) throw 'no right to set some privs';
        }

        // set
        let db = dbpool.get('share');

        let r = await db.collection('adminuser').updateOne(
            {_id: _id, rank: {$gt: sess.rank}},
            {$set: {priv: pstr}},
        );

        // update priv
        if (r.modifiedCount == 1) {
            session.update_priv(_id, pstr);
        }

        res.json({}).end();
    } catch(e) {
        res.json({err: typeof e == 'string' ? e : e.message}).end();
    }
}));

// ============================================================================

module.exports = router;
