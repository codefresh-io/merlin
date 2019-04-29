var Environment = /** @class */ (function () {
    function Environment() {
        this.operators = [];
        this.components = [];
    }
    Environment.prototype.setVersion = function (v) {
        this.version = v;
        return this;
    };
    Environment.prototype.addOperator = function (operator) {
        this.operators.push(operator);
        return this;
    };
    Environment.prototype.addComponent = function (component) {
        this.components.push(component);
        return this;
    };
    Environment.prototype.build = function () {
        return JSON.stringify(this);
    };
    return Environment;
}());
;
var Component = /** @class */ (function () {
    function Component(name, spec) {
        this.name = name;
        this.spec = spec;
    }
    return Component;
}());
;
var Param = /** @class */ (function () {
    function Param(opt) {
        opt = opt || {};
        this.name = opt.name;
        this.description = opt.description;
        this.required = opt.required;
        this.envVar = opt.envVar;
        this.interactive = opt.interactive;
    }
    return Param;
}());
;
var Operator = /** @class */ (function () {
    function Operator(options) {
        this.name = options.name;
        this.description = options.description;
        this.params = options.params;
        this.scope = options.scope || 'environment';
    }
    return Operator;
}());
;
var Command = /** @class */ (function () {
    function Command(options) {
        this.name = options.name;
        this.description = options.description;
        this.workDir = options.workDir;
        this.exec = options.exec;
        this.program = options.program;
    }
    return Command;
}());
;
var CommandSet = /** @class */ (function () {
    function CommandSet() {
        this.commands = [];
    }
    CommandSet.prototype.addCommand = function (cmd) {
        this.commands.push(cmd);
        return this;
    };
    CommandSet.prototype.build = function () {
        return JSON.stringify(this.commands);
    };
    return CommandSet;
}());
;
