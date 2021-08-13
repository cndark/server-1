
var path  = require('path');

var dbpool = require('../lib/dbpool');

// ============================================================================

async function game_onoff(a, op, fromid, toid) {
    let cmd = `${path.join(a.dir, './deploy/game_cmd.sh')} ${op} ${fromid} ${toid} <<< Yes`;
    return await shell_exec(cmd);
}

async function game_next_wait_id(a, all) {
    let cmd = `${path.join(a.dir, './deploy/next_wait_id.sh')} ${all ? 'all' : ''}`;
    return await shell_exec(cmd);
}

async function game_opennew(a, id) {
    let cmd = `${path.join(a.dir, './deploy/open_new.sh')} ${id} <<< Yes`;
    return await shell_exec(cmd);
}

async function game_reg(a, op, name) {
    await dbpool.exec(a.conf.common.db_center, async db => {
        await db.collection('svrlist').updateOne(
            {_id: name},
            {$set: {closereg: op == 'close'}},
        )
    });
}

// ============================================================================

module.exports = {
    game_onoff:        game_onoff,
    game_next_wait_id: game_next_wait_id,
    game_opennew:      game_opennew,
    game_reg:          game_reg,
}
