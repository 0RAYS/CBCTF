import test from 'node:test';
import assert from 'node:assert/strict';
import {
  createGuideConfig,
  createGuideService,
  createGuideUser,
  createGuideGroup,
  createGuideWriteFile,
  normalizeGuideConfigForMode,
  getGuideServiceTargets,
} from '../src/components/features/Admin/challenges/editor/guideModel.js';
import { parseComposeYamlToGuideConfig } from '../src/components/features/Admin/challenges/editor/composeParser.js';
import { buildGuidedComposeYaml } from '../src/components/features/Admin/challenges/editor/composeSerializer.js';
import {
  validateGuidedCompose,
  validateRawCompose,
} from '../src/components/features/Admin/challenges/editor/composeValidation.js';
import { buildChallengePayload } from '../src/components/features/Admin/challenges/editor/challengePayload.js';
import {
  buildNetworkTopology,
  defaultNetworkPolicy,
  normalizeNetworkPolicy,
} from '../src/components/features/Admin/challenges/editor/networkPolicy.js';

const service = (overrides = {}) => ({
  ...createGuideService(),
  name: 'web',
  image: 'nginx:alpine',
  ports: [{ published: 'http', target: '80', protocol: 'tcp' }],
  ...overrides,
});
const sharedConfig = () => ({ services: [service()], networks: [] });
const vpcConfig = () => ({
  services: [service({ networks: [{ name: 'lan', ipv4Address: '10.1.0.2', macAddress: '' }] })],
  networks: [{ name: 'lan', subnet: '10.1.0.0/24', gateway: '10.1.0.1' }],
});
const errorKey = (key) => `admin.challengeModal.composeGuide.validation.${key}`;

test('guide factories do not share nested state; mode normalization only disables VMs', () => {
  const first = createGuideService();
  first.userData.users.push(createGuideUser());
  assert.deepEqual(createGuideService().userData.users, []);
  assert.deepEqual(createGuideConfig(), { services: [], networks: [] });
  const config = { services: [service({ kubeVirt: true, bootloader: 'efi' })], networks: [] };
  const normalized = normalizeGuideConfigForMode(config);
  assert.equal(normalized.services[0].kubeVirt, false);
  assert.equal(normalized.services[0].bootloader, 'efi');
  assert.equal(config.services[0].kubeVirt, true);
  const vpc = vpcConfig();
  assert.equal(normalizeGuideConfigForMode(vpc), vpc);
});

test('shared Compose serialization retains field ordering and x-volumes flag defaults', () => {
  const config = sharedConfig();
  Object.assign(config.services[0], {
    containerName: 'web-container',
    cpus: '0.5',
    memLimit: '128m',
    workingDir: '/app',
    command: 'sh\n-c\nrun',
    environment: [{ key: 'FLAG', value: 'leet{example}' }],
    volumes: [{ target: '/flag', content: '' }],
  });
  const yaml = buildGuidedComposeYaml(config);
  assert.equal(
    yaml,
    [
      'services:',
      '  web:',
      '    container_name: web-container',
      '    image: nginx:alpine',
      '    cpus: 0.5',
      '    mem_limit: 128m',
      '    ports:',
      '      - mode: ingress',
      '        target: 80',
      '        published: http',
      '        protocol: tcp',
      '    environment:',
      '      - "FLAG=leet{example}"',
      '    x-volumes:',
      '      - path: "/flag"',
      '        content: "uuid{}"',
      '    working_dir: /app',
      '    command:',
      '      - "sh"',
      '      - "-c"',
      '      - "run"',
    ].join('\n')
  );
  const parsed = parseComposeYamlToGuideConfig(yaml);
  assert.equal(parsed.ok, true);
  assert.equal(parsed.config.services[0].volumes[0].content, 'uuid{}');
  assert.deepEqual(validateRawCompose(yaml).list, []);
});

test('VM Compose round trips x-kubevirt, x-boot, x-cloudinit and static network addresses', () => {
  const config = vpcConfig();
  const vm = config.services[0];
  Object.assign(vm, { kubeVirt: true, memLimit: '1g', bootloader: 'efi', secureBoot: true });
  vm.networks[0].macAddress = '02:00:00:00:00:01';
  vm.userData.users = [
    {
      ...createGuideUser(),
      name: 'ctf',
      gecos: 'CTF User',
      groups: ['ctf'],
      sudo: ['ALL=(ALL) NOPASSWD:ALL'],
      shell: '/bin/sh',
      homedir: '/home/ctf',
      lockPasswd: true,
      passwd: 'hash',
      plainTextPasswd: 'secret',
      sshAuthorizedKeys: ['ssh-ed25519 user-key'],
      noCreateHome: true,
      system: true,
    },
  ];
  vm.userData.groups = [{ ...createGuideGroup(), name: 'ctf', members: ['ctf', 'root'] }];
  vm.userData.sshAuthorizedKeys = ['ssh-ed25519 root-key'];
  vm.userData.writeFiles = [
    {
      ...createGuideWriteFile(),
      path: '/etc/example',
      content: '# retained in block\nhello: world\n\nlast line',
      owner: 'root:root',
      permissions: '0644',
      encoding: 'text/plain',
      append: true,
      defer: true,
    },
  ];
  const yaml = buildGuidedComposeYaml(config);
  for (const extension of ['x-kubevirt: true', 'x-boot:', 'x-cloudinit:']) assert.ok(yaml.includes(extension));
  for (const excluded of ['ports:', 'environment:', 'x-volumes:', 'working_dir:', 'command:'])
    assert.ok(!yaml.includes(excluded));
  const parsed = parseComposeYamlToGuideConfig(yaml);
  assert.equal(parsed.ok, true);
  assert.deepEqual(parsed.config.services[0].userData, vm.userData);
  assert.deepEqual(parsed.config.services[0].networks, vm.networks);
  assert.equal(parsed.config.services[0].secureBoot, true);
  assert.deepEqual(validateRawCompose(yaml).list, []);
  assert.equal(buildGuidedComposeYaml(parsed.config), yaml);
});

test('serializer omits incomplete list entries and VM-only fields in shared mode', () => {
  const config = sharedConfig();
  Object.assign(config.services[0], {
    name: '',
    kubeVirt: true,
    bootloader: 'efi',
    secureBoot: true,
    ports: [
      { published: '', target: '', protocol: '' },
      { published: '', target: '80', protocol: '' },
    ],
    volumes: [{ target: '', content: 'uuid{}' }],
    environment: [{ key: '', value: 'ignored' }],
  });
  const yaml = buildGuidedComposeYaml(config);
  assert.ok(yaml.includes('  service1:'));
  assert.ok(yaml.includes('protocol: tcp'));
  assert.ok(!yaml.includes('x-kubevirt'));
  assert.ok(!yaml.includes('x-boot'));
  assert.ok(!yaml.includes('x-cloudinit'));
  assert.ok(!yaml.includes('environment:'));
  assert.ok(!yaml.includes('x-volumes:'));
  assert.equal(buildGuidedComposeYaml(createGuideConfig()), 'services:');
});

test('parser accepts supported handwritten short ports, mapping environment, comments and CRLF', () => {
  const yaml = [
    '# handwritten',
    'services:',
    '  web:',
    '    image: nginx:alpine # image',
    '    ports:',
    '      - "http:80/tcp"',
    "      - '53/udp'",
    '    environment:',
    '      FLAG: "static{example}"',
  ].join('\r\n');
  const parsed = parseComposeYamlToGuideConfig(yaml);
  assert.equal(parsed.ok, true);
  assert.deepEqual(parsed.config.services[0].ports, [
    { published: 'http', target: '80', protocol: 'tcp' },
    { published: '', target: '53', protocol: 'udp' },
  ]);
  assert.deepEqual(parsed.config.services[0].environment, [{ key: 'FLAG', value: 'static{example}' }]);
  assert.deepEqual(validateRawCompose(yaml).list, []);
});

for (const [name, yaml, expected, line] of [
  ['tabs', 'services:\n\tweb:', 'noTabs', 2],
  ['top-level field', 'version: 3\nservices:', 'unsupportedTopLevel', 1],
  ['service scope', 'services:\n    image: nginx', 'serviceScope', 2],
  ['service field', 'services:\n  web:\n    restart: always', 'unsupportedServiceField', 3],
  ['port mode', 'services:\n  web:\n    ports:\n      - mode: host', 'portModeInvalid', 4],
  ['long port mode', 'services:\n  web:\n    ports:\n      - target: 80\n        mode: host', 'portModeInvalid', 5],
  ['indentation', 'services:\n  web:\n     image: nginx', 'unparseableLine', 3],
  [
    'network external',
    'services:\n  web:\n    image: nginx\nnetworks:\n  lan:\n    external: maybe',
    'networkExternalInvalid',
    6,
  ],
]) {
  test(`parser retains ${name} error and source line`, () => {
    const parsed = parseComposeYamlToGuideConfig(yaml);
    assert.equal(parsed.ok, false);
    assert.ok(parsed.errors.includes(`${errorKey(expected)} ${line}`));
    const raw = validateRawCompose(yaml);
    assert.equal(raw.config, null);
    assert.deepEqual(raw.fields, {});
    assert.deepEqual(raw.list, parsed.errors);
  });
}

test('missing services remains a parse error distinct from guide validation', () => {
  assert.deepEqual(parseComposeYamlToGuideConfig('').errors, [errorKey('servicesMissing')]);
  assert.deepEqual(validateGuidedCompose(createGuideConfig()).fields, { services: [errorKey('serviceRequired')] });
});

test('validation retains indexed service and list field errors', () => {
  const config = sharedConfig();
  Object.assign(config.services[0], {
    name: '',
    image: '',
    cpus: 'bad',
    memLimit: 'bad',
    ports: [
      { published: 'bad name', target: 'x', protocol: 'sctp' },
      { published: '', target: 'x', protocol: '' },
    ],
    environment: [
      { key: '1BAD', value: '' },
      { key: '', value: '' },
    ],
    volumes: [
      { target: '/flag', content: '' },
      { target: '/flag', content: '' },
      { target: '', content: '' },
    ],
  });
  const { fields } = validateGuidedCompose(config);
  for (const [path, key] of [
    ['name', 'serviceNameRequired'],
    ['image', 'imageRequired'],
    ['cpus', 'cpusNumber'],
    ['memLimit', 'memLimitFormat'],
    ['ports.0.target', 'portTargetNumber'],
    ['ports.1.target', 'portTargetUnique'],
    ['ports.0.published', 'portNameFormat'],
    ['ports.0.protocol', 'portProtocolInvalid'],
    ['environment.0.key', 'envKeyFormat'],
    ['environment.1.key', 'envKeyRequired'],
    ['volumes.1.target', 'volumeTargetUnique'],
    ['volumes.2.target', 'volumeTargetRequired'],
  ])
    assert.ok(fields[`service.0.${path}`].includes(errorKey(key)), path);
});

test('validation enforces required open ports and duplicate service/container names', () => {
  const config = {
    services: [service({ ports: [], containerName: 'same' }), service({ ports: [], containerName: 'same' })],
    networks: [],
  };
  const { fields } = validateGuidedCompose(config);
  assert.deepEqual(fields.ports, [errorKey('portRequired')]);
  assert.deepEqual(fields['service.1.name'], [errorKey('serviceNameUnique')]);
  assert.deepEqual(fields['service.1.containerName'], [errorKey('containerNameUnique')]);
  config.services[0].ports.push({ target: '', published: '', protocol: '' });
  assert.deepEqual(validateGuidedCompose(config).fields['service.0.ports.0.target'], [errorKey('portTargetRequired')]);
});

test('VPC validation reserves gateways, enforces unique IPs, subnet membership and VM requirements', () => {
  const config = vpcConfig();
  const vm = config.services[0];
  Object.assign(vm, { kubeVirt: true, memLimit: '', bootloader: 'invalid' });
  vm.networks[0].ipv4Address = '10.1.0.1';
  let fields = validateGuidedCompose(config).fields;
  assert.deepEqual(fields['service.0.networks.0.ipv4Address'], [errorKey('ipUnique')]);
  assert.deepEqual(fields['service.0.networks.0.macAddress'], [errorKey('macAddressRequired')]);
  assert.deepEqual(fields['service.0.memLimit'], [errorKey('memLimitRequired')]);
  assert.deepEqual(fields['service.0.bootloader'], [errorKey('bootloaderInvalid')]);
  vm.networks[0].ipv4Address = '10.2.0.2';
  vm.networks[0].macAddress = 'invalid';
  fields = validateGuidedCompose(config).fields;
  assert.deepEqual(fields['service.0.networks.0.ipv4Address'], [errorKey('ipOutOfSubnet')]);
  assert.deepEqual(fields['service.0.networks.0.macAddress'], [errorKey('macAddressInvalid')]);
  config.services.push(service({ name: 'other', networks: [{ ...vm.networks[0] }] }));
  assert.ok(validateGuidedCompose(config).fields['service.1.networks.0.ipv4Address'].includes(errorKey('ipUnique')));
});

test('network definition and selection errors retain their original paths', () => {
  const config = vpcConfig();
  config.networks.push({ ...config.networks[0] });
  let fields = validateGuidedCompose(config).fields;
  assert.deepEqual(fields['network.1.name'], [errorKey('networkDefinitionNameUnique')]);
  assert.deepEqual(fields['network.1.subnet'], [errorKey('subnetUnique')]);
  assert.deepEqual(fields['network.1.gateway'], [errorKey('gatewayUnique')]);
  config.networks = [{ name: '', subnet: '', gateway: '' }];
  fields = validateGuidedCompose(config).fields;
  assert.deepEqual(fields['network.0.name'], [errorKey('networkDefinitionNameRequired')]);
  assert.deepEqual(fields['network.0.subnet'], [errorKey('subnetRequired')]);
  assert.deepEqual(fields['network.0.gateway'], [errorKey('gatewayRequired')]);
  assert.deepEqual(fields['service.0.networks'], [errorKey('serviceNetworksWithoutDefinitions')]);
  config.networks = [{ name: 'lan', subnet: 'bad', gateway: 'bad' }];
  fields = validateGuidedCompose(config).fields;
  assert.deepEqual(fields['network.0.subnet'], [errorKey('subnetInvalid')]);
  assert.deepEqual(fields['network.0.gateway'], [errorKey('gatewayInvalid'), errorKey('gatewayUnique')]);
  config.networks[0] = { name: 'lan', subnet: '10.1.0.0/24', gateway: '10.2.0.1' };
  assert.deepEqual(validateGuidedCompose(config).fields['network.0.gateway'], [errorKey('gatewayOutOfSubnet')]);
  config.services[0].networks = [];
  assert.deepEqual(validateGuidedCompose(config).fields['service.0.networks'], [errorKey('serviceNetworkRequired')]);
});

test('policy normalization accepts persisted uppercase rules without mutating input', () => {
  const policy = { service: 'web', from: [{ CIDR: '10.1.0.0/24', Except: ['10.1.0.2'] }], to: null };
  const original = structuredClone(policy);
  assert.deepEqual(normalizeNetworkPolicy(policy), {
    service: 'web',
    from: [{ cidr: '10.1.0.0/24', except: ['10.1.0.2'] }],
    to: [],
  });
  assert.deepEqual(policy, original);
  assert.equal(defaultNetworkPolicy({ service: 'web' }).service, 'web');
  assert.equal(defaultNetworkPolicy().to[0].except.length, 4);
});

for (const direction of ['from', 'to', 'both']) {
  test(`topology treats JSON null ${direction} directions as empty arrays`, () => {
    const config = vpcConfig();
    config.services.push(service({ name: 'db', networks: [{ name: 'lan', ipv4Address: '10.1.0.3' }] }));
    const policies = ['web', 'db'].map((name) => ({
      service: name,
      from: direction === 'to' ? [{ cidr: '10.1.0.0/24', except: null }] : null,
      to: direction === 'from' ? [{ cidr: '10.1.0.0/24', except: null }] : null,
    }));
    const original = structuredClone(policies);
    const normalized = policies.map(normalizeNetworkPolicy);
    const topology = buildNetworkTopology(config, policies);
    assert.deepEqual(topology, buildNetworkTopology(config, normalized));
    assert.equal(topology.connections.length, 2);
    assert.ok(topology.connections.every((connection) => connection.allowed === (direction === 'from')));
    assert.deepEqual(policies, original);
  });
}

test('topology accepts a nil policy collection serialized as JSON null', () => {
  const config = vpcConfig();
  config.services.push(service({ name: 'db', networks: [{ name: 'lan', ipv4Address: '10.1.0.3' }] }));
  assert.deepEqual(buildNetworkTopology(config, null), buildNetworkTopology(config, []));
});

test('topology combines outbound and inbound policy, exceptions and target names', () => {
  const config = vpcConfig();
  config.services[0].containerName = 'frontend';
  config.services.push(service({ name: 'db', networks: [{ name: 'lan', ipv4Address: '10.1.0.3' }] }));
  assert.deepEqual(getGuideServiceTargets(config), [
    { label: 'frontend', service: 'frontend' },
    { label: 'db', service: 'db' },
  ]);
  let policies = [{ service: 'frontend', to: [{ CIDR: '0.0.0.0/0', Except: [] }] }];
  let topology = buildNetworkTopology(config, policies);
  assert.equal(topology.connections[0].allowed, true);
  assert.equal(topology.connections[1].reasonKey, 'outboundBlocked');
  assert.deepEqual(topology.connections[0].networks, ['lan']);
  policies.push({ service: 'db', from: [{ cidr: '10.1.0.0/24', except: ['10.1.0.2/32'] }] });
  assert.equal(buildNetworkTopology(config, policies).connections[0].reasonKey, 'inboundBlocked');
  policies = [{ service: 'frontend', to: [{ cidr: '0.0.0.0/0', except: ['10.0.0.0/8'] }] }];
  assert.equal(buildNetworkTopology(config, policies).connections[0].reasonKey, 'outboundBlocked');
  topology = buildNetworkTopology(sharedConfig());
  assert.equal(topology.nodes[0].x, 50);
  assert.equal(topology.nodes[0].y, 50);
  assert.deepEqual(topology.connections, []);
});

for (const type of ['static', 'dynamic']) {
  test(`${type} create/update payload preserves flag representations and IDs`, () => {
    const challenge = {
      id: 7,
      name: 'Example',
      description: 'Description',
      category: 'web',
      type,
      flags: ['static{one}', { id: 12, value: 'leet{two}' }, { id: 0, value: '' }],
      generator_image: 'generator:v1',
      docker_compose: 'must not submit',
      network_policies: [],
    };
    const original = structuredClone(challenge);
    const created = buildChallengePayload(challenge, 'add');
    const updated = buildChallengePayload(challenge, 'edit');
    assert.deepEqual(created.flags, ['static{one}', 'leet{two}', '']);
    assert.deepEqual(updated.flags, [
      { id: 0, value: 'static{one}' },
      { id: 12, value: 'leet{two}' },
      { id: 0, value: '' },
    ]);
    assert.equal(Object.hasOwn(created, 'docker_compose'), false);
    assert.equal(Object.hasOwn(updated, 'network_policies'), false);
    assert.equal(Object.hasOwn(created, 'id'), false);
    assert.equal(created.generator_image, type === 'dynamic' ? 'generator:v1' : undefined);
    assert.deepEqual(challenge, original);
  });
}

test('pods payload omits unchanged Compose, preserves handwritten text exactly and never submits flags', () => {
  const yaml = '# retain comments\r\nservices:\r\n  web:\r\n    image: nginx\r\n';
  const challenge = {
    id: 7,
    name: 'Example',
    description: '',
    category: 'web',
    type: 'pods',
    flags: ['ignored'],
    docker_compose: yaml,
    network_policies: [{ service: 'web', from: [{ CIDR: '10.1.0.0/24', Except: [] }], to: [] }],
  };
  const original = [{ id: 7, value: yaml }];
  const created = buildChallengePayload(challenge, 'add', original);
  const updated = buildChallengePayload(challenge, 'edit', original);
  assert.equal(created.docker_compose, yaml);
  assert.equal(Object.hasOwn(created, 'flags'), false);
  assert.equal(Object.hasOwn(updated, 'docker_compose'), false);
  assert.deepEqual(updated.network_policies, [{ from: [{ cidr: '10.1.0.0/24', except: [] }], to: [] }]);
  assert.equal(
    buildChallengePayload({ ...challenge, docker_compose: `${yaml}\n` }, 'edit', original).docker_compose,
    `${yaml}\n`
  );
  assert.equal(buildChallengePayload({ ...challenge, docker_compose: '' }, 'edit', original).docker_compose, '');
  assert.equal(Object.hasOwn(buildChallengePayload(challenge, 'edit', []), 'docker_compose'), false);
  assert.equal(buildChallengePayload({ ...challenge, docker_compose: '' }, 'add').docker_compose, '');
});

test('VPC payload includes policy service even when Compose is omitted from update', () => {
  const yaml = buildGuidedComposeYaml(vpcConfig());
  const challenge = {
    id: 1,
    type: 'pods',
    docker_compose: yaml,
    network_policies: [
      { service: 'web', from: [], to: [{ CIDR: '0.0.0.0/0', Except: ['10.0.0.0/8'] }] },
      { from: null, to: null },
    ],
  };
  const payload = buildChallengePayload(challenge, 'edit', [{ id: 1, value: yaml }]);
  assert.equal(Object.hasOwn(payload, 'docker_compose'), false);
  assert.deepEqual(payload.network_policies, [
    { service: 'web', from: [], to: [{ cidr: '0.0.0.0/0', except: ['10.0.0.0/8'] }] },
    { service: '', from: [], to: [] },
  ]);
});

test('network selection validates missing, duplicate and undefined names and invalid IPv4 values', () => {
  const config = vpcConfig();
  config.services[0].networks = [
    { name: '', ipv4Address: '', macAddress: '' },
    { name: 'missing', ipv4Address: '999.1.1.1', macAddress: '' },
    { name: 'lan', ipv4Address: '10.1.0.2', macAddress: '' },
    { name: 'lan', ipv4Address: '10.1.0.3', macAddress: '' },
  ];
  const { fields } = validateGuidedCompose(config);
  assert.deepEqual(fields['service.0.networks.0.name'], [errorKey('networkNameRequired')]);
  assert.deepEqual(fields['service.0.networks.0.ipv4Address'], [errorKey('ipRequired')]);
  assert.deepEqual(fields['service.0.networks.1.name'], [errorKey('networkUndefined')]);
  assert.deepEqual(fields['service.0.networks.1.ipv4Address'], [errorKey('ipInvalid')]);
  assert.deepEqual(fields['service.0.networks.3.name'], [errorKey('networkDuplicateSelect')]);
});

test('serializer preserves existing quote and backslash escaping for list values', () => {
  const config = sharedConfig();
  config.services[0].environment = [{ key: 'VALUE', value: 'a"b\\c' }];
  const yaml = buildGuidedComposeYaml(config);
  assert.ok(yaml.includes('      - "VALUE=a\\"b\\\\c"'));
  // The legacy parser strips surrounding quotes, but does not decode YAML escapes.
  assert.equal(parseComposeYamlToGuideConfig(yaml).config.services[0].environment[0].value, 'a\\"b\\\\c');
});
