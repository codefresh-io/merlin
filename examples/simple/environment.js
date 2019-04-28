function _createStandardNodejsComponent(name) {
    var port = {
        name: 'port',
        envVar: 'PORT',
        default: 9000
    };
    return new Component(name, {
        ports: [port]
    });
}

function _createNonStandardComponent(name, ports) {
    var nonStandardPorts = _.map(ports, function (p, i) {
        return {
            name: 'port-' + i,
            envVar: p.envVar,
            default: p.default
        }
    });
    return new Component(name, {
        ports: nonStandardPorts
    });
}

function $create(config) {
    return new CommandSet()
        .addCommand(
            new Command({
                name: 'git',
                description: 'Checkout to master',
                workDir: config['cf-helm-path'],
                exec: [
                    'git',
                    'checkout',
                    'dynamic'
                ]
            })
        )
        .addCommand(
            new Command({
                name: 'git',
                description: 'Create branch',
                workDir: config['cf-helm-path'],
                exec: [
                    'git',
                    'checkout',
                    'dynamic'
                ]
            })
        )
        .addCommand(
            new Command({
                name: 'git',
                description: 'push branch to upstream',
                workDir: config['cf-helm-path'],
                exec: [
                    'git',
                    'checkout',
                    '-b',
                    'dynamic' + '-' + config.name
                ]
            })
        )
        .addCommand(
            new Command({
                name: 'wait',
                description: 'Wait for environment',
                exec: [
                    'sh',
                    '-c',
                    'codefresh logs -f $(codefresh get build --pipeline-name --branch dynamic-' + config.name + ' | awk \'NR >1\' | awk \'{ print $1}\')',
                ]
            })
        )
        .build();
}

function $connect(config, component) {
    var env = [];
    env.push('MERLIN_COMPONENT=' + component.name)

    var exec = [];
    _.chain(exec)
        .push('telepresence')
        .push('--context')
        .push(config.kubernetes.context)
        .push('--swap-deployment')
        .push(config.name + '-' + component.name)
        .push('--namespace')
        .push(config.name)
        .value()
    _.map(component.spec.ports, function (p) {
        var port = GetAvailablePort()
        exec.push('--expose')
        exec.push(port + ':' + p.default)
        env.push(component.name + '_' + p.name + '=' + port);
    })

    if (config.run) {
        exec.push('--run')
        exec.push(config.run)
    }


    return JSON.stringify([{
        name: 'run-telepresence',
        exec: exec,
        env: env,
        detached: true
    }]);
}

function $start(config, component) {
    var env = _.map(component.spec.ports, function (p) {
        return p.envVar + '=$' + component.name + '_' + p.name;
    })
    env.push('FORMAT_LOGS_TO_ELK=false')
    return JSON.stringify([{
            name: 'ensure-tools',
            description: 'Ensure telepresence & kubectl && jq && codefresh exist',
            env: env,
            exec: [
                'node',
                '--inspect=' + GetAvailablePort(),
                'server/index.js'
            ]
        },

    ]);
}

function build() {
    var env = new Environment();
    return env
        .addComponent(_createNonStandardComponent('cfapi', [{
            envVar: 'INTERNAL_PORT',
            default: 40000.
        }, {
            envVar: 'PORT',
            default: 80
        }]))
        .addComponent(_createNonStandardComponent('runtime-environment-manager', [{
            envVar: 'PORT',
            default: 80
        }]))
        .addComponent(_createStandardNodejsComponent('pipeline-manager'))
        .addOperator(new Operator({
            name: 'create',
            description: 'Create an environment',
            params: [{
                name: 'cf-helm-path',
                description: 'Path to cf-helm local repo',
                required: true,
                envVar: 'CF_HELM_PATH',
                interactive: true,
            }]
        }))
        .addOperator(new Operator({
            name: 'connect',
            description: 'Establish connection to the cluster and open ssh tunnel',
            scope: 'component'
        }))
        .addOperator(new Operator({
            name: 'start',
            description: 'Start application locally (assuming the connection is been established)',
            scope: 'component'
        }))
        .build();
}