"use strict";
/*
 * @Description:
 * @Autor: chenjie
 * @Date: 2021-04-15 10:33:14
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-06-28 21:28:43
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Configs = void 0;
var CsvLoader_1 = require("./CsvLoader");
var fs = require("fs");
var Configs = {
    loadCsv: function () {
        var paths = ["./bat/calcbattle/csvNormal/", "./bat/calcbattle/csvWlevel/", "./bat/calcbattle/csvTower/"];
        this.worldLevelConf = {};
        this.towerConf = {};
        for (var m = 0; m < paths.length; m++) {
            var path = paths[m];
            var fileNames = fs.readdirSync(path);
            for (var i = 0, j = fileNames.length; i < j; i++) {
                var name_1 = fileNames[i];
                var content = fs.readFileSync(path + name_1);
                var rs = CsvLoader_1.csv_open(content.toString());
                var tab = CsvLoader_1.rs_make_tab(rs);
                if (path == "./bat/calcbattle/csvWlevel/") { //关卡分段
                    for (var key in tab) {
                        var val = tab[key];
                        this.worldLevelConf[key] = val;
                    }
                }
                else if (path == "./bat/calcbattle/csvTower/") { //爬塔分段
                    for (var key in tab) {
                        var val = tab[key];
                        this.towerConf[key] = val;
                    }
                }
                else {
                    var info = name_1.split(".");
                    this[info[0] + "Conf"] = tab;
                }
            }
        }
    }
};
exports.Configs = Configs;
