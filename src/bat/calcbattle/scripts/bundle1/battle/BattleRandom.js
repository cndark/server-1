"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BattleRandom = void 0;
var BattleRandom = /** @class */ (function () {
    function BattleRandom(seed) {
        this._seed = 0;
        this._A = 0;
        this._M = 0;
        if (!seed || seed <= 0) {
            seed = 1;
        }
        this._seed = seed;
        this._A = 16807;
        this._M = 1999999999;
        this.random();
        this.random();
        this.random();
        this.random();
    }
    BattleRandom.prototype.random = function (start, end) {
        this._seed = Math.floor((this._seed * this._A) % this._M + 0.5);
        var ret = 0;
        if (typeof (start) == "number") {
            if (typeof (end) == "number") {
                ret = Math.floor(this._seed / this._M * (end - start + 1)) + start;
            }
            else {
                ret = Math.floor(this._seed / this._M * start) + 1;
            }
        }
        else {
            ret = this._seed / this._M;
        }
        return ret;
    };
    return BattleRandom;
}());
exports.BattleRandom = BattleRandom;
