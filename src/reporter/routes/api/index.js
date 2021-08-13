
var express = require('express');
var router = express.Router();

// ============================================================================

router.use('/install', require('./install')); // 激活统计
router.use('/ccs', require('./ccs'));     // 客户端 callstack
router.use('/cslice', require('./cslice'));  // 客户端 分段下载统计
router.use('/cver', require('./cver'));    // 客户端 版本统计
router.use('/cstep', require('./cstep'));   // 客户端 启动步骤
router.use('/feedback', require('./feedback'));// 客户端 问题反馈
router.use('/glog', require('./glog'));// 客户端 日志上报
router.use('/userpush', require('./userpush')); // wx-h5 用户消息推送

// ============================================================================

module.exports = router;
