
// Code generated by go generate; DO NOT EDIT.
// using data from templates/js
package js

func TemplatesMap() map[string]string {
    templatesMap := make(map[string]string)

templatesMap["classes.js"] = `var Environment = /** @class */ (function () {
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
` 

templatesMap["classes.ts"] = `declare let _ : any

class Environment {
    version: string;
    operators: Operator[];
    components: Component[];

    constructor() {
        this.operators = [];
        this.components = [];
    }
    
    setVersion(v: string): Environment {
        this.version = v;
        return this;
    }

    addOperator(operator: Operator): Environment {
        this.operators.push(operator);
        return this;
    }
    addComponent(component: Component): Environment {
        this.components.push(component);
        return this;
    }

    build(): string {
        return JSON.stringify(this);
    }
};

class Component {
    name: string;
    spec: any;

    constructor(name: string, spec: any) {
        this.name = name;
        this.spec = spec;
    }
};

class Param {
    name: string;
    description: string;
    required: boolean;
    envVar: string;
    interactive: boolean;
    constructor(opt: any) {
        opt = opt || {}
        this.name = opt.name;
        this.description = opt.description;
        this.required = opt.required;
        this.envVar = opt.envVar;
        this.interactive = opt.interactive;
    }
};

class Operator {
    name: string;
    params: Param[];
    description: string;
    scope: string;
    constructor(options: { name: string; description: string; params: Param[]; scope: string; }) {
        this.name = options.name;
        this.description = options.description;
        this.params = options.params;
        this.scope = options.scope || 'environment';
    }
};

class Command {
    name: string;
    description: string;
    workDir: string;
    exec: string[];
    constructor(options : { name: string; description: string; workDir: string; exec: string[]; }) {
        this.name = options.name;
        this.description = options.description;
        this.workDir = options.workDir;
        this.exec = options.exec;
    }
};

class CommandSet {
    commands: Command[];
    constructor() {
        this.commands = [];
    }

    addCommand(cmd: Command): CommandSet {
        this.commands.push(cmd);
        return this;
    }

    build() : string {
        return JSON.stringify(this.commands);
    }
};` 

    return  templatesMap
}
