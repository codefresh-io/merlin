declare let _ : any
declare let GetAvailablePort: any

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
    program: string;
    constructor(options : { name: string; description: string; workDir: string ; exec: string[]; program: string; }) {
        this.name = options.name;
        this.description = options.description;
        this.workDir = options.workDir;
        this.exec = options.exec;
        this.program = options.program;
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
};