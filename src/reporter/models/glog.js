
var config = require('../../config.json');

const { Kafka } = require('kafkajs');

// ============================================================================

var kfk = {};
var producer = {};

// ============================================================================

async function init() {
    if (config.common.kfk.urls.length == 0) {
        return
    }

    kfk = new Kafka({
        brokers: config.common.kfk.urls,
    });

    producer = kfk.producer();
    await producer.connect();
}

async function send(msg) {
    if (config.common.kfk.urls.length == 0) {
        return
    }

    if (producer && msg) {
        await producer.send({
            topic: "glog_benu",
            messages: [{
                value: msg,
            }],
        })
    }
}

// ============================================================================

module.exports = {
    init: init,
    send: send,
};