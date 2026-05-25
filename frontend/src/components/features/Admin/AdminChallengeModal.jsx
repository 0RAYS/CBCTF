import { motion } from 'motion/react';
import { IconX, IconPlus, IconTrash } from '@tabler/icons-react';
import Button from '../../common/Button';
import { lazy, Suspense, useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';

const Editor = lazy(() => import('../../../lib/monacoSetup').then(() => import('@monaco-editor/react')));

const createGuideWriteFile = () => ({
  path: '',
  content: '',
  owner: '',
  permissions: '',
  encoding: '',
  append: false,
  defer: false,
});

const createGuideUser = () => ({
  name: '',
  gecos: '',
  groups: [],
  sudo: [],
  shell: '',
  homedir: '',
  lockPasswd: false,
  passwd: '',
  plainTextPasswd: '',
  sshAuthorizedKeys: [],
  noCreateHome: false,
  system: false,
});

const createGuideGroup = () => ({
  name: '',
  members: [],
});

const createGuideCloudConfig = () => ({
  users: [],
  groups: [],
  writeFiles: [],
  sshAuthorizedKeys: [],
});

/**
 * 题目管理弹窗组件
 * @param {boolean} props.isOpen - 是否显示弹窗
 * @param {string} props.mode - 模式, 'add'或'edit'
 * @param {Object} props.challenge - 当前编辑的题目对象
 * @param {Array} props.categories - 分类列表
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Function} props.onSubmit - 提交表单回调
 * @param {Function} props.onChange - 表单字段变更回调
 * @param {Function} props.onAddFlag - 添加Flag回调
 * @param {Function} props.onRemoveFlag - 移除Flag回调
 * @param {Function} props.onFlagChange - Flag变更回调
 */

const createGuideService = () => ({
  name: '',
  containerName: '',
  image: '',
  cpus: '',
  memLimit: '',
  workingDir: '',
  command: '',
  kubeVirt: false,
  bootloader: '',
  secureBoot: false,
  userData: createGuideCloudConfig(),
  ports: [],
  environment: [],
  volumes: [],
  networks: [],
});

const createGuideNetwork = () => ({
  name: '',
  external: false,
  subnet: '',
  gateway: '',
});

const createGuideConfig = () => ({
  services: [],
  networks: [],
});

const defaultNetworkPolicy = (target = {}) => ({
  ...(target.service ? { service: target.service } : {}),
  from: [
    {
      cidr: '',
      except: [''],
    },
  ],
  to: [
    {
      cidr: '0.0.0.0/0',
      except: ['10.0.0.0/8', '172.16.0.0/12', '192.168.0.0/16', '100.64.0.0/10'],
    },
  ],
});

const getGuideServiceTargets = (config) =>
  (config.services || []).map((service, index) => {
    const name = service.containerName.trim() || service.name.trim() || `service${index + 1}`;
    return {
      label: name,
      service: name,
    };
  });

const hasVpcNetworks = (config) => (config.networks || []).some((network) => network.name.trim());

const ipToNumber = (ip) => {
  const parts = String(ip || '')
    .trim()
    .split('.')
    .map((part) => Number(part));
  if (parts.length !== 4 || parts.some((part) => !Number.isInteger(part) || part < 0 || part > 255)) return null;
  return parts.reduce((value, part) => value * 256 + part, 0);
};

const cidrContainsIp = (cidr, ip) => {
  const [baseIp, prefixText = '32'] = String(cidr || '')
    .trim()
    .split('/');
  const base = ipToNumber(baseIp);
  const target = ipToNumber(ip);
  const prefix = Number(prefixText);
  if (base === null || target === null || !Number.isInteger(prefix) || prefix < 0 || prefix > 32) return false;
  if (prefix === 0) return true;
  const mask = (0xffffffff << (32 - prefix)) >>> 0;
  return (base & mask) === (target & mask);
};

const getPolicyBlocks = (blocks = []) =>
  blocks
    .map((block) => ({
      cidr: String(block?.cidr || block?.CIDR || '').trim(),
      except: block?.except || block?.Except || [],
    }))
    .filter((block) => block.cidr);

const blockAllowsIp = (block, ip) => {
  if (!cidrContainsIp(block.cidr, ip)) return false;
  return !(block.except || []).some((exceptCidr) => cidrContainsIp(exceptCidr, ip));
};

const blocksAllowAnyIp = (blocks, ips) => ips.some((ip) => blocks.some((block) => blockAllowsIp(block, ip)));

const normalizePolicyRule = (rule = {}) => ({
  cidr: rule.cidr || rule.CIDR || '',
  except: Array.isArray(rule.except) ? rule.except : Array.isArray(rule.Except) ? rule.Except : [],
});

const normalizePolicyRules = (rules = []) => (Array.isArray(rules) ? rules.map(normalizePolicyRule) : []);

const normalizeNetworkPolicy = (policy = {}) => ({
  ...policy,
  from: normalizePolicyRules(policy.from),
  to: normalizePolicyRules(policy.to),
});

const buildNetworkTopology = (config, policies = []) => {
  const services = config.services || [];
  const nodes = services.map((service, index) => {
    const name = service.containerName.trim() || service.name.trim() || `service${index + 1}`;
    return {
      id: name,
      label: name,
      image: service.image,
      networks: (service.networks || [])
        .filter((network) => network.name || network.ipv4Address)
        .map((network) => ({
          name: network.name || 'default',
          ip: network.ipv4Address || '-',
        })),
    };
  });

  const positionedNodes = nodes.map((node, index) => {
    if (nodes.length === 1) return { ...node, x: 50, y: 50 };
    const angle = (Math.PI * 2 * index) / nodes.length - Math.PI / 2;
    return {
      ...node,
      x: 50 + Math.cos(angle) * 34,
      y: 50 + Math.sin(angle) * 34,
    };
  });

  const policyByService = new Map(policies.map((policy) => [policy.service, policy]));
  const getNodeIps = (node) => node.networks.map((network) => network.ip).filter((ip) => ipToNumber(ip) !== null);
  const connections = [];

  positionedNodes.forEach((source) => {
    positionedNodes.forEach((target) => {
      if (source.id === target.id) return;

      const sourcePolicy = policyByService.get(source.id);
      const targetPolicy = policyByService.get(target.id);
      const targetIps = getNodeIps(target);
      const sourceIps = getNodeIps(source);
      const toBlocks = getPolicyBlocks(sourcePolicy?.to);
      const fromBlocks = getPolicyBlocks(targetPolicy?.from);
      const outboundAllowed = targetIps.length > 0 && blocksAllowAnyIp(toBlocks, targetIps);
      const inboundAllowed = fromBlocks.length === 0 || blocksAllowAnyIp(fromBlocks, sourceIps);
      const sharedNetworks = source.networks
        .map((network) => network.name)
        .filter((name) => target.networks.some((network) => network.name === name));

      connections.push({
        id: `${source.id}->${target.id}`,
        source,
        target,
        allowed: outboundAllowed && inboundAllowed,
        reasonKey: !outboundAllowed ? 'outboundBlocked' : !inboundAllowed ? 'inboundBlocked' : 'allowed',
        networks: sharedNetworks.length ? sharedNetworks : [],
      });
    });
  });

  return { nodes: positionedNodes, connections };
};

const normalizeGuideConfigForMode = (config) => {
  if (hasVpcNetworks(config)) return config;

  return {
    ...config,
    services: (config.services || []).map((service) => ({
      ...service,
      kubeVirt: false,
    })),
  };
};

const hasOpenPort = (config) =>
  (config.services || []).some(
    (service) => !service.kubeVirt && (service.ports || []).some((port) => port.target.trim())
  );

const emptyComposeYaml = 'services:';

const yamlQuote = (value) =>
  `"${String(value ?? '')
    .replace(/\\/g, '\\\\')
    .replace(/"/g, '\\"')}"`;

const appendYamlList = (lines, values, indent) => {
  values
    .map((value) => String(value || '').trim())
    .filter(Boolean)
    .forEach((value) => lines.push(`${indent}- ${yamlQuote(value)}`));
};

const yamlLineValues = (value) =>
  (Array.isArray(value)
    ? value
    : String(value || '')
        .replace(/\r\n/g, '\n')
        .split('\n')
  )
    .map((item) => item.trim())
    .filter(Boolean);

const appendYamlStringListField = (lines, key, value, indent) => {
  const values = yamlLineValues(value);
  if (values.length === 0) return;
  lines.push(`${indent}${key}:`);
  appendYamlList(lines, values, `${indent}  `);
};

const hasGuideUser = (user) =>
  user.name.trim() ||
  user.gecos.trim() ||
  yamlLineValues(user.groups).length > 0 ||
  yamlLineValues(user.sudo).length > 0 ||
  user.shell.trim() ||
  user.homedir.trim() ||
  user.lockPasswd ||
  user.passwd.trim() ||
  user.plainTextPasswd.trim() ||
  yamlLineValues(user.sshAuthorizedKeys).length > 0 ||
  user.noCreateHome ||
  user.system;

const hasGuideGroup = (group) => group.name.trim() || yamlLineValues(group.members).length > 0;

const isGuideCloudConfigEmpty = (cloudConfig = {}) =>
  !(cloudConfig.users || []).some(hasGuideUser) &&
  !(cloudConfig.groups || []).some(hasGuideGroup) &&
  !(cloudConfig.writeFiles || []).some(
    (file) =>
      file.path.trim() ||
      file.content.trim() ||
      file.owner.trim() ||
      file.permissions.trim() ||
      file.encoding.trim() ||
      file.append ||
      file.defer
  ) &&
  yamlLineValues(cloudConfig.sshAuthorizedKeys).length === 0;

const appendUserStringField = (lines, key, value) => {
  if (String(value || '').trim()) lines.push(`          ${key}: ${yamlQuote(String(value).trim())}`);
};

const appendUserStringListField = (lines, key, value) => {
  const values = yamlLineValues(value);
  if (values.length === 0) return;
  lines.push(`          ${key}:`);
  appendYamlList(lines, values, '            ');
};

const appendUsers = (lines, users = []) => {
  const values = users.filter(hasGuideUser);
  if (values.length === 0) return;
  lines.push('      users:');
  values.forEach((user) => {
    lines.push(`        - name: ${yamlQuote(user.name.trim())}`);
    appendUserStringField(lines, 'gecos', user.gecos);
    appendUserStringListField(lines, 'groups', user.groups);
    appendUserStringListField(lines, 'sudo', user.sudo);
    appendUserStringField(lines, 'shell', user.shell);
    appendUserStringField(lines, 'homedir', user.homedir);
    lines.push(`          lock_passwd: ${user.lockPasswd ? 'true' : 'false'}`);
    appendUserStringField(lines, 'passwd', user.passwd);
    appendUserStringField(lines, 'plain_text_passwd', user.plainTextPasswd);
    appendUserStringListField(lines, 'ssh_authorized_keys', user.sshAuthorizedKeys);
    lines.push(`          no_create_home: ${user.noCreateHome ? 'true' : 'false'}`);
    lines.push(`          system: ${user.system ? 'true' : 'false'}`);
  });
};

const appendGroups = (lines, groups = []) => {
  const values = groups.filter(hasGuideGroup);
  if (values.length === 0) return;
  lines.push('      groups:');
  values.forEach((group) => {
    lines.push(`        - name: ${yamlQuote(group.name.trim())}`);
    const members = yamlLineValues(group.members);
    if (members.length > 0) {
      lines.push('          members:');
      appendYamlList(lines, members, '            ');
    }
  });
};

const appendGuideCloudConfig = (lines, cloudConfig) => {
  if (isGuideCloudConfigEmpty(cloudConfig)) return;
  lines.push('    x-cloudinit:');
  appendUsers(lines, cloudConfig.users);
  appendGroups(lines, cloudConfig.groups);

  const writeFiles = (cloudConfig.writeFiles || []).filter(
    (file) =>
      file.path.trim() ||
      file.content.trim() ||
      file.owner.trim() ||
      file.permissions.trim() ||
      file.encoding.trim() ||
      file.append ||
      file.defer
  );
  if (writeFiles.length > 0) {
    lines.push('      write_files:');
    writeFiles.forEach((file) => {
      lines.push(`        - path: ${yamlQuote(file.path.trim())}`);
      if (file.content.trim()) {
        lines.push('          content: |');
        String(file.content || '')
          .replace(/\r\n/g, '\n')
          .split('\n')
          .forEach((item) => lines.push(`            ${item}`));
      }
      if (file.owner.trim()) lines.push(`          owner: ${yamlQuote(file.owner.trim())}`);
      if (file.permissions.trim()) lines.push(`          permissions: ${yamlQuote(file.permissions.trim())}`);
      if (file.encoding.trim()) lines.push(`          encoding: ${yamlQuote(file.encoding.trim())}`);
      if (file.append) lines.push('          append: true');
      if (file.defer) lines.push('          defer: true');
    });
  }

  appendYamlStringListField(lines, 'ssh_authorized_keys', cloudConfig.sshAuthorizedKeys, '      ');
};

const buildGuidedComposeYaml = (config) => {
  const lines = ['services:'];
  const networks = (config.networks || []).filter((network) => network.name.trim());
  const vpcMode = networks.length > 0;

  (config.services || []).forEach((service, index) => {
    const kubeVirt = vpcMode && service.kubeVirt;
    const serviceName = service.name.trim() || `service${index + 1}`;
    lines.push(`  ${serviceName}:`);
    if (service.containerName.trim()) lines.push(`    container_name: ${service.containerName.trim()}`);
    lines.push(`    image: ${service.image.trim()}`);
    if (service.cpus.trim()) lines.push(`    cpus: ${service.cpus.trim()}`);
    if (service.memLimit.trim()) lines.push(`    mem_limit: ${service.memLimit.trim()}`);

    const ports = kubeVirt ? [] : (service.ports || []).filter((port) => port.target.trim());
    if (ports.length > 0) {
      lines.push('    ports:');
      ports.forEach((port) => {
        const protocol = port.protocol || 'tcp';
        const published = port.published.trim();
        const target = port.target.trim();
        lines.push('      - mode: ingress', `        target: ${target}`);
        if (published) lines.push(`        published: ${published}`);
        lines.push(`        protocol: ${protocol}`);
      });
    }

    const envs = kubeVirt ? [] : (service.environment || []).filter((env) => env.key.trim());
    if (envs.length > 0) {
      lines.push('    environment:');
      envs.forEach((env) => lines.push(`      - ${yamlQuote(`${env.key.trim()}=${env.value}`)}`));
    }

    const volumes = kubeVirt ? [] : (service.volumes || []).filter((volume) => volume.target.trim());
    if (volumes.length > 0) {
      lines.push('    x-volumes:');
      volumes.forEach((volume) => {
        lines.push(
          '      - path: ' + yamlQuote(volume.target.trim()),
          `        content: ${yamlQuote(volume.content || 'uuid{}')}`
        );
      });
    }

    if (kubeVirt) {
      lines.push('    x-kubevirt: true');
    }

    if (kubeVirt && (service.bootloader.trim() || service.secureBoot)) {
      lines.push('    x-boot:');
      if (service.bootloader.trim()) lines.push(`      bootloader: ${service.bootloader.trim()}`);
      lines.push(`      secure_boot: ${service.secureBoot ? 'true' : 'false'}`);
    }

    if (kubeVirt) appendGuideCloudConfig(lines, service.userData || createGuideCloudConfig());

    if (!kubeVirt && service.workingDir.trim()) lines.push(`    working_dir: ${service.workingDir.trim()}`);

    const command = String(service.command || '').split('\n');
    if (!kubeVirt && command.some((item) => item.trim())) {
      lines.push('    command:');
      appendYamlList(lines, command, '      ');
    }

    if (networks.length > 0) {
      lines.push('    networks:');
      (service.networks || [])
        .filter((network) => network.name.trim())
        .forEach((network) => {
          lines.push(`      ${network.name.trim()}:`);
          lines.push(`        ipv4_address: ${network.ipv4Address.trim()}`);
          if (network.macAddress?.trim()) lines.push(`        mac_address: ${yamlQuote(network.macAddress.trim())}`);
        });
    }
  });

  if (networks.length > 0) {
    lines.push('', 'networks:');
    networks.forEach((network) => {
      lines.push(`  ${network.name.trim()}:`);
      if (network.external) lines.push('    external: true');
      lines.push(
        '    ipam:',
        '      config:',
        `        - subnet: ${network.subnet.trim()}`,
        `          gateway: ${network.gateway.trim()}`
      );
    });
  }

  return lines.join('\n');
};

const ip4Segment = '(\\d|[1-9]\\d|1\\d\\d|2([0-4]\\d|5[0-5]))';
const ip4Regex = new RegExp(`^(${ip4Segment}\\.){3}${ip4Segment}$`);
const cidrRegex = new RegExp(`^(${ip4Segment}\\.){3}${ip4Segment}/(\\d|[1-2]\\d|3[0-2])$`);
const macAddressRegex = /^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$/;

const validateIp4 = (ip) => ip4Regex.test(String(ip || '').trim());

const validateCidr = (cidr) => cidrRegex.test(String(cidr || '').trim());

const ip4ToInt = (ip) => ip.split('.').reduce((sum, part) => (sum << 8) + parseInt(part, 10), 0) >>> 0;

const isIp4InCidr = (ip) => (cidr) => {
  const trimmedIp = String(ip || '').trim();
  const trimmedCidr = String(cidr || '').trim();
  if (!validateIp4(trimmedIp) || !validateCidr(trimmedCidr)) return false;

  const [range, bits] = trimmedCidr.split('/');
  const mask = ~(2 ** (32 - Number(bits)) - 1);
  return (ip4ToInt(trimmedIp) & mask) === (ip4ToInt(range) & mask);
};

const stripYamlValue = (value = '') => {
  const trimmed = String(value).trim();
  if ((trimmed.startsWith('"') && trimmed.endsWith('"')) || (trimmed.startsWith("'") && trimmed.endsWith("'"))) {
    return trimmed.slice(1, -1);
  }
  return trimmed;
};

const parseKeyValueLine = (line) => {
  const index = line.indexOf(':');
  if (index === -1) return null;
  return {
    key: line.slice(0, index).trim(),
    value: stripYamlValue(line.slice(index + 1)),
  };
};

const appendCloudConfigLineValue = (target, field, value) => {
  if (Array.isArray(target[field])) {
    target[field] = [...target[field], stripYamlValue(value)];
    return;
  }
  target[field] = [
    ...String(target[field] || '')
      .split('\n')
      .filter(Boolean),
    stripYamlValue(value),
  ].join('\n');
};

const userFieldMap = {
  lock_passwd: 'lockPasswd',
  plain_text_passwd: 'plainTextPasswd',
  ssh_authorized_keys: 'sshAuthorizedKeys',
  no_create_home: 'noCreateHome',
};

const setGuideUserField = (user, key, value) => {
  const field = userFieldMap[key] || key;
  if (['groups', 'sudo', 'sshAuthorizedKeys'].includes(field)) {
    if (value) user[field] = [value];
    return field;
  }
  if (['lockPasswd', 'noCreateHome', 'system'].includes(field)) user[field] = value === 'true';
  else user[field] = value;
  return field;
};

const parsePortString = (value) => {
  const [main, protocol = 'tcp'] = stripYamlValue(value).split('/');
  const parts = main.split(':');
  if (parts.length === 1) return { published: '', target: parts[0] || '', protocol };
  return { published: parts[0] || '', target: parts[1] || '', protocol };
};

const parseComposeYamlToGuideConfig = (yaml, t = (key, values) => `${key}${values?.line ? ` ${values.line}` : ''}`) => {
  const lines = String(yaml || '')
    .replace(/\r\n/g, '\n')
    .split('\n')
    .map((raw, index) => ({ raw, index: index + 1, text: raw }));
  const errors = [];
  const services = [];
  const networks = [];
  let section = '';
  let service = null;
  let network = null;
  let serviceList = '';
  let serviceNetwork = null;
  let currentPort = null;
  let currentVolume = null;
  let currentUser = null;
  let currentGroup = null;
  let currentWriteFile = null;
  let cloudConfigList = '';
  let cloudConfigNestedList = '';
  let blockScalar = null;

  for (const line of lines) {
    if (/\t/.test(line.raw))
      errors.push(t('admin.challengeModal.composeGuide.validation.noTabs', { line: line.index }));
    const indent = line.text.match(/^ */)[0].length;
    const rawText = line.text.trim();

    if (blockScalar && indent >= blockScalar.indent) {
      blockScalar.lines.push(line.raw.slice(blockScalar.indent));
      blockScalar.target[blockScalar.field] = blockScalar.lines.join('\n');
      continue;
    }
    if (blockScalar && indent < blockScalar.indent) blockScalar = null;
    if (!rawText || rawText.startsWith('#')) continue;
    const text = line.text.replace(/#.*$/, '').trim();
    if (!text) continue;

    if (indent === 0) {
      section = text.replace(/:$/, '');
      service = null;
      network = null;
      serviceList = '';
      serviceNetwork = null;
      currentPort = null;
      currentVolume = null;
      currentUser = null;
      currentGroup = null;
      currentWriteFile = null;
      cloudConfigList = '';
      cloudConfigNestedList = '';
      blockScalar = null;
      if (!['services', 'networks'].includes(section)) {
        errors.push(
          t('admin.challengeModal.composeGuide.validation.unsupportedTopLevel', { line: line.index, field: section })
        );
      }
      continue;
    }

    if (section === 'services') {
      if (indent === 2 && text.endsWith(':')) {
        service = {
          name: text.slice(0, -1),
          containerName: '',
          image: '',
          cpus: '',
          memLimit: '',
          workingDir: '',
          command: '',
          kubeVirt: false,
          bootloader: '',
          secureBoot: false,
          userData: createGuideCloudConfig(),
          ports: [],
          environment: [],
          volumes: [],
          networks: [],
        };
        services.push(service);
        serviceList = '';
        serviceNetwork = null;
        currentPort = null;
        currentVolume = null;
        currentUser = null;
        currentGroup = null;
        currentWriteFile = null;
        cloudConfigList = '';
        cloudConfigNestedList = '';
        blockScalar = null;
        continue;
      }
      if (!service) {
        errors.push(t('admin.challengeModal.composeGuide.validation.serviceScope', { line: line.index }));
        continue;
      }
      if (indent === 4) {
        const pair = parseKeyValueLine(text);
        if (!pair) continue;
        serviceNetwork = null;
        currentPort = null;
        currentVolume = null;
        currentUser = null;
        currentGroup = null;
        currentWriteFile = null;
        cloudConfigList = '';
        cloudConfigNestedList = '';
        blockScalar = null;
        if (
          ['ports', 'environment', 'x-volumes', 'x-boot', 'x-cloudinit', 'command', 'networks'].includes(pair.key) &&
          pair.value === ''
        ) {
          serviceList = pair.key;
          continue;
        }
        serviceList = '';
        if (pair.key === 'container_name') service.containerName = pair.value;
        else if (pair.key === 'image') service.image = pair.value;
        else if (pair.key === 'cpus') service.cpus = pair.value;
        else if (pair.key === 'mem_limit') service.memLimit = pair.value;
        else if (pair.key === 'working_dir') service.workingDir = pair.value;
        else if (pair.key === 'x-kubevirt') service.kubeVirt = pair.value === 'true';
        else
          errors.push(
            t('admin.challengeModal.composeGuide.validation.unsupportedServiceField', {
              line: line.index,
              field: pair.key,
            })
          );
        continue;
      }
      if (indent === 6 && serviceList === 'command' && text.startsWith('- ')) {
        service.command = [...service.command.split('\n').filter(Boolean), stripYamlValue(text.slice(2))].join('\n');
        continue;
      }
      if (indent === 6 && serviceList === 'environment' && text.startsWith('- ')) {
        const env = stripYamlValue(text.slice(2));
        const equalIndex = env.indexOf('=');
        service.environment.push({
          key: equalIndex === -1 ? env : env.slice(0, equalIndex),
          value: equalIndex === -1 ? '' : env.slice(equalIndex + 1),
        });
        continue;
      }
      if (indent === 6 && serviceList === 'environment') {
        const pair = parseKeyValueLine(text);
        if (pair) {
          service.environment.push({ key: pair.key, value: pair.value });
          continue;
        }
      }
      if (indent === 6 && serviceList === 'x-boot') {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'bootloader') service.bootloader = pair.value;
        else if (pair?.key === 'secure_boot') service.secureBoot = pair.value === 'true';
        continue;
      }
      if (indent === 6 && serviceList === 'x-cloudinit') {
        const pair = parseKeyValueLine(text);
        if (!pair) continue;
        currentUser = null;
        currentGroup = null;
        currentWriteFile = null;
        cloudConfigNestedList = '';
        if (['users', 'groups', 'ssh_authorized_keys', 'write_files'].includes(pair.key)) {
          cloudConfigList = pair.key;
          if (pair.key === 'ssh_authorized_keys' && pair.value) service.userData.sshAuthorizedKeys = [pair.value];
        }
        continue;
      }
      if (indent === 8 && serviceList === 'x-cloudinit') {
        if (cloudConfigList === 'ssh_authorized_keys' && text.startsWith('- ')) {
          appendCloudConfigLineValue(service.userData, 'sshAuthorizedKeys', text.slice(2));
          continue;
        }
        if (cloudConfigList === 'users' && text.startsWith('- ')) {
          const pair = parseKeyValueLine(text.slice(2).trim());
          currentUser = createGuideUser();
          currentGroup = null;
          currentWriteFile = null;
          service.userData.users.push(currentUser);
          if (pair?.key) setGuideUserField(currentUser, pair.key, pair.value);
          cloudConfigNestedList = '';
          continue;
        }
        if (cloudConfigList === 'groups' && text.startsWith('- ')) {
          const pair = parseKeyValueLine(text.slice(2).trim());
          currentGroup = createGuideGroup();
          currentUser = null;
          currentWriteFile = null;
          service.userData.groups.push(currentGroup);
          if (pair?.key === 'name') currentGroup.name = pair.value;
          else if (pair?.key === 'members' && pair.value) currentGroup.members = [pair.value];
          cloudConfigNestedList = '';
          continue;
        }
        if (cloudConfigList === 'write_files' && text.startsWith('- ')) {
          const pair = parseKeyValueLine(text.slice(2).trim());
          currentWriteFile = createGuideWriteFile();
          service.userData.writeFiles.push(currentWriteFile);
          if (pair?.key === 'path') currentWriteFile.path = pair.value;
          else if (pair?.key === 'content') currentWriteFile.content = pair.value;
          continue;
        }
      }
      if (indent === 10 && serviceList === 'x-cloudinit' && cloudConfigList === 'users' && currentUser) {
        const pair = parseKeyValueLine(text);
        if (pair?.key) {
          cloudConfigNestedList = setGuideUserField(currentUser, pair.key, pair.value === '' ? '' : pair.value);
          continue;
        }
      }
      if (indent === 12 && serviceList === 'x-cloudinit' && cloudConfigList === 'users' && currentUser) {
        if (['groups', 'sudo', 'sshAuthorizedKeys'].includes(cloudConfigNestedList) && text.startsWith('- ')) {
          currentUser[cloudConfigNestedList] = [...currentUser[cloudConfigNestedList], stripYamlValue(text.slice(2))];
          continue;
        }
      }
      if (indent === 10 && serviceList === 'x-cloudinit' && cloudConfigList === 'groups' && currentGroup) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'name') currentGroup.name = pair.value;
        else if (pair?.key === 'members') {
          cloudConfigNestedList = 'members';
          if (pair.value) currentGroup.members = [pair.value];
        }
        continue;
      }
      if (indent === 12 && serviceList === 'x-cloudinit' && cloudConfigList === 'groups' && currentGroup) {
        if (cloudConfigNestedList === 'members' && text.startsWith('- ')) {
          currentGroup.members = [...currentGroup.members, stripYamlValue(text.slice(2))];
          continue;
        }
      }
      if (indent === 10 && serviceList === 'x-cloudinit' && cloudConfigList === 'write_files' && currentWriteFile) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'path') currentWriteFile.path = pair.value;
        else if (pair?.key === 'content') {
          currentWriteFile.content = pair.value === '|' ? '' : pair.value;
          if (pair.value === '|') blockScalar = { target: currentWriteFile, field: 'content', indent: 12, lines: [] };
        } else if (pair?.key === 'owner') currentWriteFile.owner = pair.value;
        else if (pair?.key === 'permissions') currentWriteFile.permissions = pair.value;
        else if (pair?.key === 'encoding') currentWriteFile.encoding = pair.value;
        else if (pair?.key === 'append') currentWriteFile.append = pair.value === 'true';
        else if (pair?.key === 'defer') currentWriteFile.defer = pair.value === 'true';
        continue;
      }
      if (indent === 6 && serviceList === 'x-volumes' && text.startsWith('- ')) {
        const pair = parseKeyValueLine(text.slice(2).trim());
        currentVolume = { target: '', content: '' };
        service.volumes.push(currentVolume);
        if (pair?.key === 'path') currentVolume.target = pair.value;
        else if (pair?.key === 'content') currentVolume.content = pair.value;
        continue;
      }
      if (indent === 8 && serviceList === 'x-volumes' && currentVolume) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'path') currentVolume.target = pair.value;
        else if (pair?.key === 'content') currentVolume.content = pair.value;
        continue;
      }
      if (indent === 6 && serviceList === 'ports' && text.startsWith('- ')) {
        const value = text.slice(2).trim();
        if (value.includes(':')) {
          const pair = parseKeyValueLine(value);
          if (pair && ['mode', 'target', 'published', 'protocol'].includes(pair.key)) {
            currentPort = { published: '', target: '', protocol: 'tcp' };
            service.ports.push(currentPort);
            if (pair.key === 'target') currentPort.target = pair.value;
            else if (pair.key === 'published') currentPort.published = pair.value;
            else if (pair.key === 'protocol') currentPort.protocol = pair.value;
            else if (pair.key === 'mode' && pair.value !== 'ingress') {
              errors.push(t('admin.challengeModal.composeGuide.validation.portModeInvalid', { line: line.index }));
            }
          } else {
            service.ports.push(parsePortString(value));
            currentPort = null;
          }
        } else {
          service.ports.push(parsePortString(value));
          currentPort = null;
        }
        continue;
      }
      if (indent === 8 && serviceList === 'ports' && currentPort) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'target') currentPort.target = pair.value;
        else if (pair?.key === 'published') currentPort.published = pair.value;
        else if (pair?.key === 'protocol') currentPort.protocol = pair.value;
        else if (pair?.key === 'mode' && pair.value !== 'ingress') {
          errors.push(t('admin.challengeModal.composeGuide.validation.portModeInvalid', { line: line.index }));
        }
        continue;
      }
      if (indent === 6 && serviceList === 'networks' && text.endsWith(':')) {
        serviceNetwork = { name: text.slice(0, -1), ipv4Address: '', macAddress: '' };
        service.networks.push(serviceNetwork);
        continue;
      }
      if (indent === 8 && serviceList === 'networks' && serviceNetwork) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'ipv4_address') serviceNetwork.ipv4Address = pair.value;
        else if (pair?.key === 'mac_address') serviceNetwork.macAddress = pair.value;
        continue;
      }
      errors.push(t('admin.challengeModal.composeGuide.validation.unparseableLine', { line: line.index }));
    } else if (section === 'networks') {
      if (indent === 2 && text.endsWith(':')) {
        network = { name: text.slice(0, -1), external: false, subnet: '', gateway: '' };
        networks.push(network);
      } else if (network && indent === 4) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'external') {
          if (!['true', 'false'].includes(pair.value))
            errors.push(t('admin.challengeModal.composeGuide.validation.networkExternalInvalid', { line: line.index }));
          network.external = pair.value === 'true';
        }
      } else if (network && indent >= 8) {
        const pair = parseKeyValueLine(text);
        const key = pair?.key.replace(/^-\s*/, '');
        if (key === 'subnet') network.subnet = pair.value;
        else if (key === 'gateway') network.gateway = pair.value;
      }
    }
  }

  if (services.length === 0) errors.push(t('admin.challengeModal.composeGuide.validation.servicesMissing'));
  return { ok: errors.length === 0, errors, config: { services, networks } };
};

const validateRawCompose = (yaml, t) => {
  const parsed = parseComposeYamlToGuideConfig(yaml, t);
  if (!parsed.ok) return { list: parsed.errors, fields: {}, config: null };
  const validation = validateGuidedCompose(parsed.config, t);
  return { ...validation, config: parsed.config };
};

const validateGuidedCompose = (config, t = (key) => key) => {
  const errors = [];
  const fields = {};
  const addError = (path, message) => {
    errors.push(message);
    fields[path] = [...(fields[path] || []), message];
  };
  const serviceNames = new Set();
  const serviceContainerNames = new Set();
  const networkDefinitionNames = new Set();
  const networkIps = new Set();
  const definedNetworks = (config.networks || []).filter((network) => network.name.trim());
  const networkNames = new Set(definedNetworks.map((network) => network.name.trim()));
  const networkCidrs = new Map(definedNetworks.map((network) => [network.name.trim(), network.subnet.trim()]));
  const networkAssignedIps = new Map();

  definedNetworks.forEach((network) => {
    const name = network.name.trim();
    const gateway = network.gateway.trim();
    if (!name || !validateIp4(gateway)) return;
    networkAssignedIps.set(name, new Set([gateway]));
  });

  if (!config.services?.length) addError('services', t('admin.challengeModal.composeGuide.validation.serviceRequired'));
  if (config.services?.some((service) => !service.kubeVirt) && !hasOpenPort(config)) {
    addError('ports', t('admin.challengeModal.composeGuide.validation.portRequired'));
  }
  (config.services || []).forEach((service, serviceIndex) => {
    const name = service.name.trim();
    const label = name || t('admin.challengeModal.composeGuide.serviceIndexed', { index: serviceIndex + 1 });
    if (!name)
      addError(
        `service.${serviceIndex}.name`,
        t('admin.challengeModal.composeGuide.validation.serviceNameRequired', { label })
      );
    if (name && serviceNames.has(name))
      addError(
        `service.${serviceIndex}.name`,
        t('admin.challengeModal.composeGuide.validation.serviceNameUnique', { label })
      );
    serviceNames.add(name);
    const containerName = service.containerName.trim();
    if (containerName && serviceContainerNames.has(containerName)) {
      addError(
        `service.${serviceIndex}.containerName`,
        t('admin.challengeModal.composeGuide.validation.containerNameUnique', { label })
      );
    }
    if (containerName) serviceContainerNames.add(containerName);
    if (!service.image.trim())
      addError(
        `service.${serviceIndex}.image`,
        t('admin.challengeModal.composeGuide.validation.imageRequired', { label })
      );
    if (service.cpus.trim() && !/^\d+(\.\d+)?$/.test(service.cpus.trim()))
      addError(`service.${serviceIndex}.cpus`, t('admin.challengeModal.composeGuide.validation.cpusNumber', { label }));
    if (service.memLimit.trim() && !/^\d+(\.\d+)?[bkmgBKMG]?$/.test(service.memLimit.trim())) {
      addError(
        `service.${serviceIndex}.memLimit`,
        t('admin.challengeModal.composeGuide.validation.memLimitFormat', { label })
      );
    }
    if (service.kubeVirt && !service.memLimit.trim()) {
      addError(
        `service.${serviceIndex}.memLimit`,
        t('admin.challengeModal.composeGuide.validation.memLimitRequired', { label })
      );
    }
    if (service.kubeVirt && service.bootloader.trim() && !['bios', 'efi'].includes(service.bootloader.trim())) {
      addError(
        `service.${serviceIndex}.bootloader`,
        t('admin.challengeModal.composeGuide.validation.bootloaderInvalid', { label })
      );
    }
    const servicePortTargets = new Set();
    (service.kubeVirt ? [] : service.ports || []).forEach((port, portIndex) => {
      if (!port.target.trim())
        addError(
          `service.${serviceIndex}.ports.${portIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.portTargetRequired', { label, index: portIndex + 1 })
        );
      if (port.target.trim() && !/^\d+$/.test(port.target.trim()))
        addError(
          `service.${serviceIndex}.ports.${portIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.portTargetNumber', { label, index: portIndex + 1 })
        );
      if (port.target.trim() && servicePortTargets.has(port.target.trim())) {
        addError(
          `service.${serviceIndex}.ports.${portIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.portTargetUnique', { label, index: portIndex + 1 })
        );
      }
      if (port.target.trim()) servicePortTargets.add(port.target.trim());
      if (port.published.trim() && !/^[A-Za-z0-9_-]+$/.test(port.published.trim())) {
        addError(
          `service.${serviceIndex}.ports.${portIndex}.published`,
          t('admin.challengeModal.composeGuide.validation.portNameFormat', { label, index: portIndex + 1 })
        );
      }
      if (port.protocol.trim() && !['tcp', 'udp'].includes(port.protocol.trim())) {
        addError(
          `service.${serviceIndex}.ports.${portIndex}.protocol`,
          t('admin.challengeModal.composeGuide.validation.portProtocolInvalid', { label, index: portIndex + 1 })
        );
      }
    });
    (service.kubeVirt ? [] : service.environment || []).forEach((env, envIndex) => {
      if (!env.key.trim())
        addError(
          `service.${serviceIndex}.environment.${envIndex}.key`,
          t('admin.challengeModal.composeGuide.validation.envKeyRequired', { label, index: envIndex + 1 })
        );
      if (env.key.trim() && !/^[A-Za-z_][A-Za-z0-9_]*$/.test(env.key.trim())) {
        addError(
          `service.${serviceIndex}.environment.${envIndex}.key`,
          t('admin.challengeModal.composeGuide.validation.envKeyFormat', { label, index: envIndex + 1 })
        );
      }
    });
    const serviceVolumeTargets = new Set();
    (service.kubeVirt ? [] : service.volumes || []).forEach((volume, volumeIndex) => {
      const target = volume.target.trim();
      if (!volume.target.trim())
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.volumeTargetRequired', { label, index: volumeIndex + 1 })
        );
      if (target && serviceVolumeTargets.has(target))
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.volumeTargetUnique', { label, index: volumeIndex + 1 })
        );
      if (target) serviceVolumeTargets.add(target);
    });
    if (networkNames.size === 0 && service.networks?.length > 0) {
      addError(
        `service.${serviceIndex}.networks`,
        t('admin.challengeModal.composeGuide.validation.serviceNetworksWithoutDefinitions', { label })
      );
    }
    if (networkNames.size > 0) {
      if (!service.networks?.length)
        addError(
          `service.${serviceIndex}.networks`,
          t('admin.challengeModal.composeGuide.validation.serviceNetworkRequired', { label })
        );
      const serviceNetworkNames = new Set();
      (service.networks || []).forEach((network, networkIndex) => {
        const networkName = network.name.trim();
        const ipv4Address = network.ipv4Address.trim();
        const networkLabel =
          networkName || t('admin.challengeModal.composeGuide.networkIndexed', { index: networkIndex + 1 });
        if (!networkName)
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.name`,
            t('admin.challengeModal.composeGuide.validation.networkNameRequired', { label, network: networkLabel })
          );
        if (networkName && serviceNetworkNames.has(networkName))
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.name`,
            t('admin.challengeModal.composeGuide.validation.networkDuplicateSelect', { label, network: networkName })
          );
        if (networkName) serviceNetworkNames.add(networkName);
        if (networkName && !networkNames.has(networkName))
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.name`,
            t('admin.challengeModal.composeGuide.validation.networkUndefined', { label, network: network.name })
          );
        if (!ipv4Address)
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.ipv4Address`,
            t('admin.challengeModal.composeGuide.validation.ipRequired', { label, network: networkLabel })
          );
        if (ipv4Address && !validateIp4(ipv4Address))
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.ipv4Address`,
            t('admin.challengeModal.composeGuide.validation.ipInvalid', { label, network: networkLabel })
          );
        if (networkName && validateIp4(ipv4Address)) {
          const assignedIps = networkAssignedIps.get(networkName) || new Set();
          if (assignedIps.has(ipv4Address)) {
            addError(
              `service.${serviceIndex}.networks.${networkIndex}.ipv4Address`,
              t('admin.challengeModal.composeGuide.validation.ipUnique', { label, network: networkLabel })
            );
          }
          assignedIps.add(ipv4Address);
          networkAssignedIps.set(networkName, assignedIps);
        }
        if (
          networkName &&
          ipv4Address &&
          networkCidrs.has(networkName) &&
          !isIp4InCidr(ipv4Address)(networkCidrs.get(networkName))
        ) {
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.ipv4Address`,
            t('admin.challengeModal.composeGuide.validation.ipOutOfSubnet', { label, network: networkLabel })
          );
        }
        if (service.kubeVirt && !network.macAddress?.trim()) {
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.macAddress`,
            t('admin.challengeModal.composeGuide.validation.macAddressRequired', { label, network: networkLabel })
          );
        }
        if (network.macAddress?.trim() && !macAddressRegex.test(network.macAddress.trim())) {
          addError(
            `service.${serviceIndex}.networks.${networkIndex}.macAddress`,
            t('admin.challengeModal.composeGuide.validation.macAddressInvalid', { label, network: networkLabel })
          );
        }
      });
    }
  });
  (config.networks || []).forEach((network, networkIndex) => {
    const label =
      network.name.trim() || t('admin.challengeModal.composeGuide.networkIndexed', { index: networkIndex + 1 });
    if (!network.name.trim())
      addError(
        `network.${networkIndex}.name`,
        t('admin.challengeModal.composeGuide.validation.networkDefinitionNameRequired', { label })
      );
    if (network.name.trim() && networkDefinitionNames.has(network.name.trim())) {
      addError(
        `network.${networkIndex}.name`,
        t('admin.challengeModal.composeGuide.validation.networkDefinitionNameUnique', { label })
      );
    }
    if (network.name.trim()) networkDefinitionNames.add(network.name.trim());
    if (!network.subnet.trim())
      addError(
        `network.${networkIndex}.subnet`,
        t('admin.challengeModal.composeGuide.validation.subnetRequired', { label })
      );
    if (network.subnet.trim() && !validateCidr(network.subnet))
      addError(
        `network.${networkIndex}.subnet`,
        t('admin.challengeModal.composeGuide.validation.subnetInvalid', { label })
      );
    if (network.subnet.trim() && networkIps.has(network.subnet.trim())) {
      addError(
        `network.${networkIndex}.subnet`,
        t('admin.challengeModal.composeGuide.validation.subnetUnique', { label })
      );
    }
    if (network.subnet.trim()) networkIps.add(network.subnet.trim());
    if (!network.gateway.trim())
      addError(
        `network.${networkIndex}.gateway`,
        t('admin.challengeModal.composeGuide.validation.gatewayRequired', { label })
      );
    if (network.gateway.trim() && !validateIp4(network.gateway))
      addError(
        `network.${networkIndex}.gateway`,
        t('admin.challengeModal.composeGuide.validation.gatewayInvalid', { label })
      );
    if (network.gateway.trim() && networkIps.has(network.gateway.trim())) {
      addError(
        `network.${networkIndex}.gateway`,
        t('admin.challengeModal.composeGuide.validation.gatewayUnique', { label })
      );
    }
    if (network.gateway.trim()) networkIps.add(network.gateway.trim());
    if (network.gateway.trim() && validateCidr(network.subnet) && !isIp4InCidr(network.gateway)(network.subnet)) {
      addError(
        `network.${networkIndex}.gateway`,
        t('admin.challengeModal.composeGuide.validation.gatewayOutOfSubnet', { label })
      );
    }
  });
  return { list: errors, fields };
};

function GuideListHeader({ title, addLabel, onAdd, disabled = false }) {
  return (
    <div className="flex justify-between items-center pt-2">
      <span className="text-xs font-mono text-neutral-400">{title}</span>
      <Button
        variant="ghost"
        size="sm"
        align="icon-left"
        icon={<IconPlus size={12} />}
        className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs !h-8"
        disabled={disabled}
        onClick={onAdd}
      >
        {addLabel}
      </Button>
    </div>
  );
}

function IconButton({ onClick }) {
  return (
    <Button
      variant="ghost"
      size="icon"
      className="!bg-transparent !text-red-400 hover:!text-red-300 !w-8 !h-10"
      onClick={onClick}
    >
      <IconTrash size={14} />
    </Button>
  );
}

function GuideField({ label, children }) {
  return (
    <label className="block space-y-1">
      <span className="text-[11px] font-mono text-neutral-500">{label}</span>
      {children}
    </label>
  );
}

function GuideErrors({ errors }) {
  if (!errors?.length) return null;
  return (
    <div className="mt-1 space-y-1">
      {errors.map((error, index) => (
        <div key={index} className="text-[11px] font-mono text-red-300">
          {error}
        </div>
      ))}
    </div>
  );
}

function NetworkTopologyPreview({ topology, t }) {
  if (topology.nodes.length === 0) {
    return (
      <div className="rounded-md border border-neutral-700 bg-black/20 p-4 text-center font-mono text-sm text-neutral-500">
        {t('admin.challengeModal.topology.empty')}
      </div>
    );
  }

  return (
    <div className="rounded-md border border-neutral-700 bg-black/20 p-3">
      <div className="mb-3 flex items-center justify-between gap-3">
        <div>
          <div className="text-sm font-mono text-neutral-100">{t('admin.challengeModal.topology.title')}</div>
          <div className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.topology.subtitle')}</div>
        </div>
        <div className="flex items-center gap-3 text-xs font-mono text-neutral-400">
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-geek-400" /> {t('admin.challengeModal.topology.allow')}
          </span>
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-red-400" /> {t('admin.challengeModal.topology.deny')}
          </span>
        </div>
      </div>

      <div className="relative h-[420px] overflow-hidden rounded-md border border-neutral-800 bg-black/20">
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 100 100" preserveAspectRatio="none">
          <defs>
            <marker
              id="topology-arrow-allow"
              viewBox="0 0 10 10"
              refX="8"
              refY="5"
              markerWidth="3"
              markerHeight="3"
              orient="auto-start-reverse"
            >
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#22c55e" />
            </marker>
            <marker
              id="topology-arrow-deny"
              viewBox="0 0 10 10"
              refX="8"
              refY="5"
              markerWidth="3"
              markerHeight="3"
              orient="auto-start-reverse"
            >
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#f87171" />
            </marker>
          </defs>
          {topology.connections.map((connection, index) => {
            const offset = (index % 2 === 0 ? 1 : -1) * 2.5;
            const midX = (connection.source.x + connection.target.x) / 2 + offset;
            const midY = (connection.source.y + connection.target.y) / 2 - offset;
            return (
              <path
                key={connection.id}
                d={`M ${connection.source.x} ${connection.source.y} Q ${midX} ${midY} ${connection.target.x} ${connection.target.y}`}
                fill="none"
                stroke={connection.allowed ? '#22c55e' : '#f87171'}
                strokeWidth="0.35"
                strokeDasharray={connection.allowed ? 'none' : '1.2 1.2'}
                markerEnd={`url(#${connection.allowed ? 'topology-arrow-allow' : 'topology-arrow-deny'})`}
                opacity={connection.allowed ? 0.75 : 0.55}
              />
            );
          })}
        </svg>

        {topology.nodes.map((node) => (
          <div
            key={node.id}
            className="absolute w-40 -translate-x-1/2 -translate-y-1/2 rounded-md border border-neutral-700 bg-neutral-950/95 p-2"
            style={{ left: `${node.x}%`, top: `${node.y}%` }}
          >
            <div className="truncate text-sm font-mono text-neutral-50" title={node.label}>
              {node.label}
            </div>
            {node.image ? (
              <div className="mt-0.5 truncate text-[10px] font-mono text-neutral-500" title={node.image}>
                {node.image}
              </div>
            ) : null}
            <div className="mt-2 space-y-1">
              {node.networks.length > 0 ? (
                node.networks.map((network, index) => (
                  <div
                    key={`${network.name}-${index}`}
                    className="rounded border border-neutral-700 bg-black/40 px-1.5 py-1"
                  >
                    <div className="truncate text-[10px] font-mono text-neutral-500" title={network.name}>
                      {network.name}
                    </div>
                    <div className="truncate text-xs font-mono text-geek-300" title={network.ip}>
                      {network.ip}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.topology.noIp')}</div>
              )}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-3 max-h-52 space-y-2 overflow-y-auto pr-1">
        {topology.connections.map((connection) => (
          <div
            key={connection.id}
            className={`rounded border px-2 py-1.5 font-mono text-xs ${
              connection.allowed
                ? 'border-geek-400/25 bg-geek-400/10 text-geek-200'
                : 'border-red-400/25 bg-red-400/10 text-red-200'
            }`}
          >
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-neutral-200">{connection.source.label}</span>
              <span>{connection.allowed ? '->' : '-/->'}</span>
              <span className="text-neutral-200">{connection.target.label}</span>
              <span className="text-neutral-500">
                {connection.networks.length
                  ? connection.networks.join(', ')
                  : t('admin.challengeModal.topology.crossNetwork')}
              </span>
            </div>
            <div className="mt-1 text-neutral-400">
              {t(`admin.challengeModal.topology.reasons.${connection.reasonKey}`)}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function AdminChallengeModal({
  isOpen = false,
  mode = 'add',
  challenge = {
    name: '',
    description: '',
    category: '',
    type: 'static',
    generator_image: '',
    flags: [''],
    docker_compose: '',
    network_policies: [
      {
        from: [
          {
            cidr: '',
            except: [''],
          },
        ],
        to: [
          {
            cidr: '',
            except: [''],
          },
        ],
      },
    ],
  },
  categories = [],
  onClose,
  onSubmit,
  onChange,
  onAddFlag,
  onRemoveFlag,
  onFlagChange,
}) {
  const { t } = useTranslation();
  const ct = (key, options) => t(`admin.challengeModal.composeGuide.${key}`, options);

  const [guideConfig, setGuideConfig] = useState(createGuideConfig);
  const guideGeneratedComposeRef = useRef('');

  const guideValidation = validateGuidedCompose(guideConfig, t);
  const guideValidationErrors = guideValidation.list;
  const guideFieldErrors = guideValidation.fields;
  const rawValidation = validateRawCompose(challenge.docker_compose || emptyComposeYaml, t);
  const vpcMode = hasVpcNetworks(guideConfig);
  const policyTargets = getGuideServiceTargets(guideConfig);
  const networkTopology = buildNetworkTopology(guideConfig, challenge.network_policies || []);

  const syncGuideConfig = (nextConfig) => {
    const normalizedConfig = normalizeGuideConfigForMode(nextConfig);
    const dockerCompose = buildGuidedComposeYaml(normalizedConfig);
    guideGeneratedComposeRef.current = dockerCompose;
    setGuideConfig(normalizedConfig);
    onChange({ ...challenge, docker_compose: dockerCompose });
  };

  useEffect(() => {
    if (!isOpen || challenge.type !== 'pods') {
      guideGeneratedComposeRef.current = '';
      setGuideConfig(createGuideConfig());
      return;
    }
    const dockerCompose = challenge.docker_compose || '';
    if (!dockerCompose) {
      guideGeneratedComposeRef.current = '';
      setGuideConfig(createGuideConfig());
      return;
    }
    if (guideGeneratedComposeRef.current === dockerCompose) return;
    guideGeneratedComposeRef.current = '';
    const parsed = parseComposeYamlToGuideConfig(dockerCompose, t);
    if (parsed.ok) setGuideConfig(normalizeGuideConfigForMode(parsed.config));
  }, [isOpen, challenge.id, challenge.type, challenge.docker_compose, t]);

  // docker_compose 更新
  const updateDockerCompose = (value) => {
    const finalValue = value || '';
    guideGeneratedComposeRef.current = '';
    const parsed = parseComposeYamlToGuideConfig(finalValue, t);
    if (!parsed.ok) {
      onChange({ ...challenge, docker_compose: finalValue });
      return;
    }
    const normalizedConfig = normalizeGuideConfigForMode(parsed.config);
    setGuideConfig(normalizedConfig);
    onChange({ ...challenge, docker_compose: finalValue });
  };

  const addGuideService = () => {
    syncGuideConfig({
      ...guideConfig,
      services: [
        ...guideConfig.services,
        {
          ...createGuideService(),
          name: '',
          containerName: '',
        },
      ],
    });
  };

  const removeGuideService = (serviceIndex) => {
    syncGuideConfig({
      ...guideConfig,
      services: guideConfig.services.filter((_, index) => index !== serviceIndex),
    });
  };

  const updateGuideService = (serviceIndex, field, value) => {
    const services = guideConfig.services.map((service, index) =>
      index === serviceIndex ? { ...service, [field]: value } : service
    );
    syncGuideConfig({ ...guideConfig, services });
  };

  const updateGuideServiceList = (serviceIndex, field, itemIndex, itemField, value) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      return {
        ...service,
        [field]: service[field].map((item, currentIndex) =>
          currentIndex === itemIndex ? { ...item, [itemField]: value } : item
        ),
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const updateGuideCloudConfigList = (serviceIndex, field, itemIndex, value) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          [field]: userData[field].map((item, currentIndex) => (currentIndex === itemIndex ? value : item)),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const updateGuideCloudConfigObject = (serviceIndex, field, itemIndex, itemField, value) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          [field]: userData[field].map((item, currentIndex) =>
            currentIndex === itemIndex ? { ...item, [itemField]: value } : item
          ),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const updateGuideCloudConfigNestedList = (serviceIndex, field, itemIndex, listField, valueIndex, value) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          [field]: userData[field].map((item, currentIndex) =>
            currentIndex === itemIndex
              ? {
                  ...item,
                  [listField]: item[listField].map((currentValue, currentValueIndex) =>
                    currentValueIndex === valueIndex ? value : currentValue
                  ),
                }
              : item
          ),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideCloudConfigNestedListItem = (serviceIndex, field, itemIndex, listField) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          [field]: userData[field].map((item, currentIndex) =>
            currentIndex === itemIndex ? { ...item, [listField]: [...item[listField], ''] } : item
          ),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const removeGuideCloudConfigNestedListItem = (serviceIndex, field, itemIndex, listField, valueIndex) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          [field]: userData[field].map((item, currentIndex) =>
            currentIndex === itemIndex
              ? {
                  ...item,
                  [listField]: item[listField].filter((_, currentValueIndex) => currentValueIndex !== valueIndex),
                }
              : item
          ),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideCloudConfigObject = (serviceIndex, field, item) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return { ...service, userData: { ...userData, [field]: [...userData[field], item] } };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const removeGuideCloudConfigObject = (serviceIndex, field, itemIndex) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return { ...service, userData: { ...userData, [field]: userData[field].filter((_, i) => i !== itemIndex) } };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideCloudConfigListItem = (serviceIndex, field) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return { ...service, userData: { ...userData, [field]: [...userData[field], ''] } };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const removeGuideCloudConfigListItem = (serviceIndex, field, itemIndex) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return { ...service, userData: { ...userData, [field]: userData[field].filter((_, i) => i !== itemIndex) } };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const updateGuideWriteFile = (serviceIndex, fileIndex, field, value) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: {
          ...userData,
          writeFiles: userData.writeFiles.map((file, currentIndex) =>
            currentIndex === fileIndex ? { ...file, [field]: value } : file
          ),
        },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideWriteFile = (serviceIndex) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return { ...service, userData: { ...userData, writeFiles: [...userData.writeFiles, createGuideWriteFile()] } };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const removeGuideWriteFile = (serviceIndex, fileIndex) => {
    const services = guideConfig.services.map((service, index) => {
      if (index !== serviceIndex) return service;
      const userData = { ...createGuideCloudConfig(), ...service.userData };
      return {
        ...service,
        userData: { ...userData, writeFiles: userData.writeFiles.filter((_, i) => i !== fileIndex) },
      };
    });
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideServiceListItem = (serviceIndex, field, item) => {
    const services = guideConfig.services.map((service, index) =>
      index === serviceIndex ? { ...service, [field]: [...service[field], item] } : service
    );
    syncGuideConfig({ ...guideConfig, services });
  };

  const removeGuideServiceListItem = (serviceIndex, field, itemIndex) => {
    const services = guideConfig.services.map((service, index) =>
      index === serviceIndex
        ? { ...service, [field]: service[field].filter((_, currentIndex) => currentIndex !== itemIndex) }
        : service
    );
    syncGuideConfig({ ...guideConfig, services });
  };

  const addGuideNetwork = () => {
    syncGuideConfig({
      ...guideConfig,
      networks: [
        ...guideConfig.networks,
        {
          ...createGuideNetwork(),
          name: '',
          subnet: '',
          gateway: '',
        },
      ],
    });
  };

  const updateGuideNetwork = (networkIndex, field, value) => {
    syncGuideConfig({
      ...guideConfig,
      networks: guideConfig.networks.map((network, index) =>
        index === networkIndex ? { ...network, [field]: value } : network
      ),
    });
  };

  const removeGuideNetwork = (networkIndex) => {
    const removedName = guideConfig.networks[networkIndex]?.name;
    syncGuideConfig({
      networks: guideConfig.networks.filter((_, index) => index !== networkIndex),
      services: guideConfig.services.map((service) => ({
        ...service,
        networks: service.networks.filter((network) => network.name !== removedName),
      })),
    });
  };

  // 网络策略操作
  const addNetworkPolicy = () => {
    const newNetworkPolicies = [
      ...(challenge.network_policies || []),
      defaultNetworkPolicy(vpcMode ? policyTargets[0] : {}),
    ];
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  const removeNetworkPolicy = (policyIndex) => {
    const newNetworkPolicies = challenge.network_policies.filter((_, i) => i !== policyIndex);
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  const updatePolicyTarget = (policyIndex, value) => {
    const newNetworkPolicies = [...(challenge.network_policies || [])];
    newNetworkPolicies[policyIndex] = {
      ...newNetworkPolicies[policyIndex],
      service: value,
    };
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 from/to 规则
  const addPolicyRule = (policyIndex, ruleType) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);

    policy[ruleType] = [
      ...(policy[ruleType] || []),
      {
        cidr: '',
        except: [''],
      },
    ];

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 移除 from/to 规则
  const removePolicyRule = (policyIndex, ruleType, ruleIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);

    policy[ruleType] = (policy[ruleType] || []).filter((_, i) => i !== ruleIndex);

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 cidr 值
  const updatePolicyCidr = (policyIndex, ruleType, ruleIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);
    const rule = normalizePolicyRule(policy[ruleType]?.[ruleIndex]);

    rule.cidr = value;

    policy[ruleType] = [...(policy[ruleType] || [])];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 except 规则
  const addPolicyExcept = (policyIndex, ruleType, ruleIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);
    const rule = normalizePolicyRule(policy[ruleType]?.[ruleIndex]);

    rule.except = [...(rule.except || []), ''];

    policy[ruleType] = [...(policy[ruleType] || [])];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 移除 except 规则
  const removePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);
    const rule = normalizePolicyRule(policy[ruleType]?.[ruleIndex]);

    rule.except = (rule.except || []).filter((_, i) => i !== exceptIndex);

    policy[ruleType] = [...(policy[ruleType] || [])];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 except 值
  const updatePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = normalizeNetworkPolicy(newNetworkPolicies[policyIndex]);
    const rule = normalizePolicyRule(policy[ruleType]?.[ruleIndex]);

    rule.except = [...(rule.except || [])];
    rule.except[exceptIndex] = value;

    policy[ruleType] = [...(policy[ruleType] || [])];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 常用样式类
  const inputBaseClass =
    'w-full h-10 bg-black/20 border border-neutral-300/30 rounded-md px-4 text-neutral-50 focus:outline-none focus:border-geek-400';
  const selectClass = 'select-custom select-custom-md';
  const textareaClass =
    'w-full h-20 bg-black/20 border border-neutral-300/30 rounded-md px-4 py-2 text-neutral-50 focus:outline-none focus:border-geek-400 resize-none';
  const isEditMode = mode === 'edit';
  const modalClass = isEditMode
    ? 'w-[min(96vw,1800px)] h-[min(86vh,980px)] bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden flex flex-col'
    : 'w-full max-w-7xl bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden';
  const contentClass = isEditMode
    ? 'flex-1 min-h-0 overflow-y-auto overscroll-contain p-3 sm:p-4 lg:p-5 [&_*]:min-w-0'
    : 'p-6 max-h-[70vh] overflow-y-auto';

  const renderCloudConfigList = (service, serviceIndex, field, label, placeholder) => (
    <div className="space-y-2">
      <GuideListHeader
        title={label}
        addLabel={ct('actions.add')}
        onAdd={() => addGuideCloudConfigListItem(serviceIndex, field)}
      />
      {(service.userData?.[field] || []).map((value, itemIndex) => (
        <div key={itemIndex} className="grid grid-cols-[1fr_32px] gap-2">
          <input
            className={inputBaseClass}
            value={value}
            placeholder={placeholder}
            onChange={(e) => updateGuideCloudConfigList(serviceIndex, field, itemIndex, e.target.value)}
          />
          <IconButton onClick={() => removeGuideCloudConfigListItem(serviceIndex, field, itemIndex)} />
        </div>
      ))}
    </div>
  );

  const renderCloudConfigNestedList = (serviceIndex, field, itemIndex, item, listField, label, placeholder) => (
    <div className="space-y-2">
      <GuideListHeader
        title={label}
        addLabel={ct('actions.add')}
        onAdd={() => addGuideCloudConfigNestedListItem(serviceIndex, field, itemIndex, listField)}
      />
      {(item[listField] || []).map((value, valueIndex) => (
        <div key={valueIndex} className="grid grid-cols-[1fr_32px] gap-2">
          <input
            className={inputBaseClass}
            value={value}
            placeholder={placeholder}
            onChange={(e) =>
              updateGuideCloudConfigNestedList(serviceIndex, field, itemIndex, listField, valueIndex, e.target.value)
            }
          />
          <IconButton
            onClick={() => removeGuideCloudConfigNestedListItem(serviceIndex, field, itemIndex, listField, valueIndex)}
          />
        </div>
      ))}
    </div>
  );

  const podNoticeLines = [
    t('admin.challengeModal.podsNotice.flagFormat', {
      format: '`static{}`, `leet{}`, `uuid{}`',
    }),
    t('admin.challengeModal.podsNotice.flagPrefix'),
    t('admin.challengeModal.podsNotice.flagVolume'),
    '',
    t('admin.challengeModal.podsNotice.services'),
    '',
    t('admin.challengeModal.podsNotice.networks'),
    t('admin.challengeModal.podsNotice.networksWithIp'),
    t('admin.challengeModal.podsNotice.networksShared'),
    t('admin.challengeModal.podsNotice.portUnique'),
    '',
    t('admin.challengeModal.podsNotice.exampleIndent'),
  ];

  const yamlNoticeLines = [
    t('admin.challengeModal.yamlNotice.flagFormat', {
      format: '`static{}`, `leet{}`, `uuid{}`',
    }),
    t('admin.challengeModal.yamlNotice.flagPrefix'),
    t('admin.challengeModal.yamlNotice.flagVolume'),
    t('admin.challengeModal.yamlNotice.dockerParams'),
    t('admin.challengeModal.yamlNotice.containerUnique'),
    t('admin.challengeModal.yamlNotice.networkIpRequired'),
    t('admin.challengeModal.yamlNotice.exposeNetwork'),
    t('admin.challengeModal.yamlNotice.noNetworkMerge'),
    t('admin.challengeModal.yamlNotice.containerNetwork'),
  ];

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <motion.div
        className={modalClass}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* 标题栏 */}
        <div className="flex shrink-0 justify-between items-center p-3 sm:p-4 border-b border-neutral-700 bg-neutral-950/50">
          <h2 className="text-xl font-mono text-neutral-50">
            {mode === 'add'
              ? t('admin.challengeModal.title.add')
              : mode === 'edit'
                ? t('admin.challengeModal.title.edit')
                : t('admin.challengeModal.title.delete')}
          </h2>
          <Button
            variant="ghost"
            size="icon"
            className="!bg-transparent !text-neutral-400 hover:!text-neutral-300"
            onClick={onClose}
          >
            <IconX size={18} />
          </Button>
        </div>

        {/* 表单内容 */}
        <div className={contentClass}>
          {mode === 'delete' ? (
            <div className="text-center py-12">
              <div className="mb-8">
                <div className="mx-auto w-20 h-20 bg-red-400/20 rounded-full flex items-center justify-center mb-6">
                  <IconTrash size={40} className="text-red-400" />
                </div>
                <h3 className="text-2xl font-mono text-neutral-50 mb-4">{t('admin.challengeModal.delete.title')}</h3>
                <p className="text-neutral-400 font-mono text-lg mb-2">{t('admin.challengeModal.delete.prompt')}</p>
                <p className="text-red-400 font-mono text-sm">{t('admin.challengeModal.delete.warning')}</p>
              </div>
            </div>
          ) : (
            <div className={isEditMode ? 'space-y-3' : 'space-y-4'}>
              {/* 基本信息 */}
              <div className={isEditMode ? 'rounded-md border border-neutral-700/70 bg-black/10 p-3' : 'mb-4'}>
                <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.basic')}</h3>

                <div className="space-y-3">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-3 lg:gap-4">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.name')}
                      </label>
                      <input
                        type="text"
                        value={challenge.name}
                        onChange={(e) => onChange({ ...challenge, name: e.target.value })}
                        className={inputBaseClass}
                        placeholder={t('admin.challengeModal.placeholders.name')}
                        required
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.category')}
                      </label>
                      <input
                        type="text"
                        value={challenge.category}
                        onChange={(e) => onChange({ ...challenge, category: e.target.value })}
                        className={selectClass}
                        placeholder={t('admin.challengeModal.placeholders.category')}
                        list="category-options"
                      />
                      <datalist id="category-options">
                        {categories.map((category) => (
                          <option key={category} value={category} />
                        ))}
                      </datalist>
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.type')}
                      </label>
                      <select
                        value={challenge.type}
                        onChange={(e) => onChange({ ...challenge, type: e.target.value })}
                        className={selectClass}
                        required
                      >
                        <option value="static">{t('admin.challengeModal.types.static')}</option>
                        <option value="dynamic">{t('admin.challengeModal.types.dynamic')}</option>
                        <option value="pods">{t('admin.challengeModal.types.pods')}</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.description')}
                    </label>
                    <textarea
                      value={challenge.description}
                      onChange={(e) => onChange({ ...challenge, description: e.target.value })}
                      className={textareaClass}
                      placeholder={t('admin.challengeModal.placeholders.description')}
                    />
                  </div>
                </div>
              </div>

              {/* 如果是动态类型题目, 显示输入 generator */}
              {challenge.type === 'dynamic' && (
                <div className="border-t border-neutral-700 pt-3 lg:pt-4">
                  <h3 className="text-lg font-mono text-neutral-50 mb-3">
                    {t('admin.challengeModal.sections.generator')}
                  </h3>
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.generatorImage')}
                    </label>
                    <input
                      type="text"
                      value={challenge.generator_image}
                      onChange={(e) =>
                        onChange({
                          ...challenge,
                          generator_image: e.target.value,
                        })
                      }
                      className={inputBaseClass}
                      placeholder={t('admin.challengeModal.placeholders.generatorImage')}
                      required
                    />
                  </div>
                </div>
              )}

              {/* flag 设置 - 非pods类型才显示 */}
              {challenge.type !== 'pods' && (
                <div className="border-t border-neutral-700 pt-3 lg:pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.flags')}</h3>
                    <Button
                      variant="primary"
                      size="sm"
                      align="icon-left"
                      icon={<IconPlus size={16} />}
                      onClick={onAddFlag}
                    >
                      {t('admin.challengeModal.actions.addFlag')}
                    </Button>
                  </div>
                  <div className="space-y-3">
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {challenge.type === 'static' ? 'static{}' : 'leet{} / uuid{}'}
                    </label>
                    {challenge.flags.map((flag, index) => (
                      <div key={index} className="flex gap-2 items-center">
                        <input
                          type="text"
                          value={
                            typeof flag === 'string'
                              ? flag
                              : flag.value || (challenge.type === 'static' ? 'static{}' : 'uuid{}')
                          }
                          onChange={(e) => onFlagChange(index, e.target.value)}
                          className={inputBaseClass}
                          placeholder={t('admin.challengeModal.placeholders.flag', { index: index + 1 })}
                        />
                        {challenge.flags.length > 1 && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                            onClick={() => onRemoveFlag(index)}
                          >
                            <IconTrash size={18} />
                          </Button>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* pods类型的flag说明 */}
              {challenge.type === 'pods' && (
                <div className="border-t border-neutral-700 pt-3 lg:pt-4">
                  <h3 className="text-lg font-mono text-neutral-50 mb-3">
                    {t('admin.challengeModal.sections.podsNotice')}
                  </h3>
                  <div className="p-3 bg-geek-400/10 border border-geek-400/20 rounded-md">
                    <p className="text-sm text-geek-400 font-mono">
                      {podNoticeLines.map((line, index) =>
                        line ? (
                          <span key={index}>
                            {line}
                            <br />
                          </span>
                        ) : (
                          <br key={index} />
                        )
                      )}
                    </p>
                  </div>
                </div>
              )}

              {/* 容器设置 */}
              {challenge.type === 'pods' && (
                <div className="border-t border-neutral-700 pt-3 lg:pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">
                      {t('admin.challengeModal.sections.containers')}
                    </h3>
                  </div>

                  <div className="space-y-6">
                    <div className="border border-neutral-700 rounded-md p-3 lg:p-4 space-y-4 bg-black/10">
                      {/* YAML 配置 */}
                      <div>
                        <div className="flex flex-wrap justify-between items-center gap-3 mb-3">
                          <label className="block text-sm font-mono text-neutral-400">docker-compose.yaml</label>
                        </div>

                        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
                          <div className="space-y-4">
                            <div className="flex justify-between items-center">
                              <h4 className="text-sm font-mono text-neutral-200">{ct('sections.services')}</h4>
                              <Button
                                variant="ghost"
                                size="sm"
                                align="icon-left"
                                icon={<IconPlus size={14} />}
                                onClick={addGuideService}
                              >
                                {ct('actions.addService')}
                              </Button>
                            </div>
                            {guideConfig.services.map((service, serviceIndex) => (
                              <div
                                key={serviceIndex}
                                className="border border-neutral-700 rounded-md p-3 bg-black/20 space-y-3"
                              >
                                <div className="flex justify-between items-center">
                                  <span className="text-xs font-mono text-neutral-400">
                                    {ct('serviceIndexed', { index: serviceIndex + 1 })}
                                  </span>
                                  {guideConfig.services.length > 1 && (
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      className="!bg-transparent !text-red-400 hover:!text-red-300 !w-8 !h-8"
                                      onClick={() => removeGuideService(serviceIndex)}
                                    >
                                      <IconTrash size={14} />
                                    </Button>
                                  )}
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                  <GuideField label={ct('fields.serviceName')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.name}
                                      required
                                      placeholder={ct('placeholders.serviceName')}
                                      onChange={(e) => updateGuideService(serviceIndex, 'name', e.target.value)}
                                    />
                                    <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.name`]} />
                                  </GuideField>
                                  <GuideField label={ct('fields.containerName')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.containerName}
                                      placeholder={ct('placeholders.containerName')}
                                      onChange={(e) =>
                                        updateGuideService(serviceIndex, 'containerName', e.target.value)
                                      }
                                    />
                                    <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.containerName`]} />
                                  </GuideField>
                                  <GuideField label={ct('fields.image')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.image}
                                      required
                                      placeholder={ct('placeholders.image')}
                                      onChange={(e) => updateGuideService(serviceIndex, 'image', e.target.value)}
                                    />
                                    <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.image`]} />
                                  </GuideField>
                                  <GuideField label={ct('fields.cpus')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.cpus}
                                      placeholder={ct('placeholders.cpus')}
                                      onChange={(e) => updateGuideService(serviceIndex, 'cpus', e.target.value)}
                                    />
                                    <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.cpus`]} />
                                  </GuideField>
                                  <GuideField label={ct('fields.memLimit')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.memLimit}
                                      placeholder={ct('placeholders.memLimit')}
                                      onChange={(e) => updateGuideService(serviceIndex, 'memLimit', e.target.value)}
                                    />
                                    <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.memLimit`]} />
                                  </GuideField>
                                  {!service.kubeVirt ? (
                                    <>
                                      <GuideField label={ct('fields.workingDir')}>
                                        <input
                                          className={inputBaseClass}
                                          value={service.workingDir}
                                          placeholder={ct('placeholders.workingDir')}
                                          onChange={(e) =>
                                            updateGuideService(serviceIndex, 'workingDir', e.target.value)
                                          }
                                        />
                                      </GuideField>
                                      <GuideField label={ct('fields.command')}>
                                        <textarea
                                          className={textareaClass}
                                          value={service.command}
                                          placeholder={ct('placeholders.command')}
                                          onChange={(e) => updateGuideService(serviceIndex, 'command', e.target.value)}
                                        />
                                      </GuideField>
                                    </>
                                  ) : null}
                                  {vpcMode ? (
                                    <GuideField label={ct('fields.kubeVirt')}>
                                      <select
                                        className={selectClass}
                                        value={service.kubeVirt ? 'true' : 'false'}
                                        onChange={(e) =>
                                          updateGuideService(serviceIndex, 'kubeVirt', e.target.value === 'true')
                                        }
                                      >
                                        <option value="false">false</option>
                                        <option value="true">true</option>
                                      </select>
                                    </GuideField>
                                  ) : null}
                                  {service.kubeVirt ? (
                                    <>
                                      <GuideField label={ct('fields.bootloader')}>
                                        <select
                                          className={selectClass}
                                          value={service.bootloader}
                                          onChange={(e) =>
                                            updateGuideService(serviceIndex, 'bootloader', e.target.value)
                                          }
                                        >
                                          <option value="">{ct('placeholders.bootloader')}</option>
                                          <option value="bios">bios</option>
                                          <option value="efi">efi</option>
                                        </select>
                                        <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.bootloader`]} />
                                      </GuideField>
                                      <GuideField label={ct('fields.secureBoot')}>
                                        <select
                                          className={selectClass}
                                          value={service.secureBoot ? 'true' : 'false'}
                                          onChange={(e) =>
                                            updateGuideService(serviceIndex, 'secureBoot', e.target.value === 'true')
                                          }
                                        >
                                          <option value="false">false</option>
                                          <option value="true">true</option>
                                        </select>
                                      </GuideField>
                                    </>
                                  ) : null}
                                </div>

                                {service.kubeVirt ? (
                                  <>
                                    <GuideListHeader
                                      title={ct('sections.users')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() => addGuideCloudConfigObject(serviceIndex, 'users', createGuideUser())}
                                    />
                                    {(service.userData?.users || []).map((user, userIndex) => (
                                      <div
                                        key={userIndex}
                                        className="space-y-2 rounded border border-neutral-700/60 p-3"
                                      >
                                        <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_1fr_32px] gap-2">
                                          <GuideField label={ct('fields.name')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.name}
                                              placeholder={ct('placeholders.name')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'name',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.gecos')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.gecos}
                                              placeholder={ct('placeholders.gecos')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'gecos',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.shell')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.shell}
                                              placeholder={ct('placeholders.shell')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'shell',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <IconButton
                                            onClick={() =>
                                              removeGuideCloudConfigObject(serviceIndex, 'users', userIndex)
                                            }
                                          />
                                        </div>
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                                          <GuideField label={ct('fields.homeDir')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.homedir}
                                              placeholder={ct('placeholders.homeDir')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'homedir',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.passwd')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.passwd}
                                              placeholder={ct('placeholders.passwd')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'passwd',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.plainTextPasswd')}>
                                            <input
                                              className={inputBaseClass}
                                              value={user.plainTextPasswd}
                                              placeholder={ct('placeholders.plainTextPasswd')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'users',
                                                  userIndex,
                                                  'plainTextPasswd',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <div className="grid grid-cols-3 gap-2">
                                            <GuideField label={ct('fields.lockPasswd')}>
                                              <select
                                                className={selectClass}
                                                value={user.lockPasswd ? 'true' : 'false'}
                                                onChange={(e) =>
                                                  updateGuideCloudConfigObject(
                                                    serviceIndex,
                                                    'users',
                                                    userIndex,
                                                    'lockPasswd',
                                                    e.target.value === 'true'
                                                  )
                                                }
                                              >
                                                <option value="false">false</option>
                                                <option value="true">true</option>
                                              </select>
                                            </GuideField>
                                            <GuideField label={ct('fields.noCreateHome')}>
                                              <select
                                                className={selectClass}
                                                value={user.noCreateHome ? 'true' : 'false'}
                                                onChange={(e) =>
                                                  updateGuideCloudConfigObject(
                                                    serviceIndex,
                                                    'users',
                                                    userIndex,
                                                    'noCreateHome',
                                                    e.target.value === 'true'
                                                  )
                                                }
                                              >
                                                <option value="false">false</option>
                                                <option value="true">true</option>
                                              </select>
                                            </GuideField>
                                            <GuideField label={ct('fields.system')}>
                                              <select
                                                className={selectClass}
                                                value={user.system ? 'true' : 'false'}
                                                onChange={(e) =>
                                                  updateGuideCloudConfigObject(
                                                    serviceIndex,
                                                    'users',
                                                    userIndex,
                                                    'system',
                                                    e.target.value === 'true'
                                                  )
                                                }
                                              >
                                                <option value="false">false</option>
                                                <option value="true">true</option>
                                              </select>
                                            </GuideField>
                                          </div>
                                        </div>
                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
                                          {renderCloudConfigNestedList(
                                            serviceIndex,
                                            'users',
                                            userIndex,
                                            user,
                                            'groups',
                                            ct('fields.groups'),
                                            ct('placeholders.groups')
                                          )}
                                          {renderCloudConfigNestedList(
                                            serviceIndex,
                                            'users',
                                            userIndex,
                                            user,
                                            'sudo',
                                            ct('fields.sudo'),
                                            ct('placeholders.sudo')
                                          )}
                                          {renderCloudConfigNestedList(
                                            serviceIndex,
                                            'users',
                                            userIndex,
                                            user,
                                            'sshAuthorizedKeys',
                                            ct('fields.sshAuthorizedKeys'),
                                            ct('placeholders.sshAuthorizedKeys')
                                          )}
                                        </div>
                                      </div>
                                    ))}

                                    <GuideListHeader
                                      title={ct('sections.groups')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() =>
                                        addGuideCloudConfigObject(serviceIndex, 'groups', createGuideGroup())
                                      }
                                    />
                                    {(service.userData?.groups || []).map((group, groupIndex) => (
                                      <div
                                        key={groupIndex}
                                        className="space-y-2 rounded border border-neutral-700/60 p-3"
                                      >
                                        <div className="grid grid-cols-[1fr_32px] gap-2">
                                          <GuideField label={ct('fields.name')}>
                                            <input
                                              className={inputBaseClass}
                                              value={group.name}
                                              placeholder={ct('placeholders.name')}
                                              onChange={(e) =>
                                                updateGuideCloudConfigObject(
                                                  serviceIndex,
                                                  'groups',
                                                  groupIndex,
                                                  'name',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <IconButton
                                            onClick={() =>
                                              removeGuideCloudConfigObject(serviceIndex, 'groups', groupIndex)
                                            }
                                          />
                                        </div>
                                        {renderCloudConfigNestedList(
                                          serviceIndex,
                                          'groups',
                                          groupIndex,
                                          group,
                                          'members',
                                          ct('fields.members'),
                                          ct('placeholders.members')
                                        )}
                                      </div>
                                    ))}

                                    {renderCloudConfigList(
                                      service,
                                      serviceIndex,
                                      'sshAuthorizedKeys',
                                      ct('fields.sshAuthorizedKeys'),
                                      ct('placeholders.sshAuthorizedKeys')
                                    )}
                                  </>
                                ) : null}

                                {service.kubeVirt ? (
                                  <>
                                    <GuideListHeader
                                      title={ct('sections.writeFiles')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() => addGuideWriteFile(serviceIndex)}
                                    />
                                    {(service.userData?.writeFiles || []).map((file, fileIndex) => (
                                      <div
                                        key={fileIndex}
                                        className="space-y-2 rounded border border-neutral-700/60 p-3"
                                      >
                                        <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_1fr_32px] gap-2">
                                          <GuideField label={ct('fields.path')}>
                                            <input
                                              className={inputBaseClass}
                                              value={file.path}
                                              placeholder={ct('placeholders.path')}
                                              onChange={(e) =>
                                                updateGuideWriteFile(serviceIndex, fileIndex, 'path', e.target.value)
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.owner')}>
                                            <input
                                              className={inputBaseClass}
                                              value={file.owner}
                                              placeholder={ct('placeholders.owner')}
                                              onChange={(e) =>
                                                updateGuideWriteFile(serviceIndex, fileIndex, 'owner', e.target.value)
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.permissions')}>
                                            <input
                                              className={inputBaseClass}
                                              value={file.permissions}
                                              placeholder={ct('placeholders.permissions')}
                                              onChange={(e) =>
                                                updateGuideWriteFile(
                                                  serviceIndex,
                                                  fileIndex,
                                                  'permissions',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <IconButton onClick={() => removeGuideWriteFile(serviceIndex, fileIndex)} />
                                        </div>
                                        <div className="grid grid-cols-1 md:grid-cols-[1fr_120px_120px] gap-2">
                                          <GuideField label={ct('fields.encoding')}>
                                            <input
                                              className={inputBaseClass}
                                              value={file.encoding}
                                              placeholder={ct('placeholders.encoding')}
                                              onChange={(e) =>
                                                updateGuideWriteFile(
                                                  serviceIndex,
                                                  fileIndex,
                                                  'encoding',
                                                  e.target.value
                                                )
                                              }
                                            />
                                          </GuideField>
                                          <GuideField label={ct('fields.append')}>
                                            <select
                                              className={selectClass}
                                              value={file.append ? 'true' : 'false'}
                                              onChange={(e) =>
                                                updateGuideWriteFile(
                                                  serviceIndex,
                                                  fileIndex,
                                                  'append',
                                                  e.target.value === 'true'
                                                )
                                              }
                                            >
                                              <option value="false">false</option>
                                              <option value="true">true</option>
                                            </select>
                                          </GuideField>
                                          <GuideField label={ct('fields.defer')}>
                                            <select
                                              className={selectClass}
                                              value={file.defer ? 'true' : 'false'}
                                              onChange={(e) =>
                                                updateGuideWriteFile(
                                                  serviceIndex,
                                                  fileIndex,
                                                  'defer',
                                                  e.target.value === 'true'
                                                )
                                              }
                                            >
                                              <option value="false">false</option>
                                              <option value="true">true</option>
                                            </select>
                                          </GuideField>
                                        </div>
                                        <GuideField label={ct('fields.content')}>
                                          <textarea
                                            className={textareaClass}
                                            value={file.content}
                                            placeholder={ct('placeholders.content')}
                                            onChange={(e) =>
                                              updateGuideWriteFile(serviceIndex, fileIndex, 'content', e.target.value)
                                            }
                                          />
                                        </GuideField>
                                      </div>
                                    ))}
                                  </>
                                ) : null}

                                {!service.kubeVirt ? (
                                  <>
                                    <GuideListHeader
                                      title={ct('sections.ports')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() =>
                                        addGuideServiceListItem(serviceIndex, 'ports', {
                                          published: '',
                                          target: '',
                                          protocol: '',
                                        })
                                      }
                                    />
                                    {service.ports.map((port, portIndex) => (
                                      <div key={portIndex} className="grid grid-cols-[1fr_1fr_90px_32px] gap-2">
                                        <GuideField label={ct('fields.name')}>
                                          <input
                                            className={inputBaseClass}
                                            value={port.published}
                                            placeholder={ct('placeholders.name')}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'ports',
                                                portIndex,
                                                'published',
                                                e.target.value
                                              )
                                            }
                                          />
                                          <GuideErrors
                                            errors={
                                              guideFieldErrors[`service.${serviceIndex}.ports.${portIndex}.published`]
                                            }
                                          />
                                        </GuideField>
                                        <GuideField label={ct('fields.target')}>
                                          <input
                                            className={inputBaseClass}
                                            value={port.target}
                                            required
                                            placeholder={ct('placeholders.target')}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'ports',
                                                portIndex,
                                                'target',
                                                e.target.value
                                              )
                                            }
                                          />
                                          <GuideErrors
                                            errors={
                                              guideFieldErrors[`service.${serviceIndex}.ports.${portIndex}.target`]
                                            }
                                          />
                                        </GuideField>
                                        <GuideField label={ct('fields.protocol')}>
                                          <select
                                            className={selectClass}
                                            value={port.protocol}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'ports',
                                                portIndex,
                                                'protocol',
                                                e.target.value
                                              )
                                            }
                                          >
                                            <option value="">{ct('fields.protocol')}</option>
                                            <option value="tcp">tcp</option>
                                            <option value="udp">udp</option>
                                          </select>
                                        </GuideField>
                                        <IconButton
                                          onClick={() => removeGuideServiceListItem(serviceIndex, 'ports', portIndex)}
                                        />
                                      </div>
                                    ))}
                                  </>
                                ) : null}

                                {!service.kubeVirt ? (
                                  <>
                                    <GuideListHeader
                                      title={ct('sections.environment')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() =>
                                        addGuideServiceListItem(serviceIndex, 'environment', { key: '', value: '' })
                                      }
                                    />
                                    {service.environment.map((env, envIndex) => (
                                      <div key={envIndex} className="grid grid-cols-[1fr_1fr_32px] gap-2">
                                        <GuideField label={ct('fields.key')}>
                                          <input
                                            className={inputBaseClass}
                                            value={env.key}
                                            required
                                            placeholder={ct('placeholders.key')}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'environment',
                                                envIndex,
                                                'key',
                                                e.target.value
                                              )
                                            }
                                          />
                                          <GuideErrors
                                            errors={
                                              guideFieldErrors[`service.${serviceIndex}.environment.${envIndex}.key`]
                                            }
                                          />
                                        </GuideField>
                                        <GuideField label={ct('fields.value')}>
                                          <input
                                            className={inputBaseClass}
                                            value={env.value}
                                            placeholder={ct('placeholders.value')}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'environment',
                                                envIndex,
                                                'value',
                                                e.target.value
                                              )
                                            }
                                          />
                                        </GuideField>
                                        <IconButton
                                          onClick={() =>
                                            removeGuideServiceListItem(serviceIndex, 'environment', envIndex)
                                          }
                                        />
                                      </div>
                                    ))}
                                  </>
                                ) : null}

                                {!service.kubeVirt ? (
                                  <>
                                    <GuideListHeader
                                      title={ct('sections.fileFlags')}
                                      addLabel={ct('actions.add')}
                                      onAdd={() =>
                                        addGuideServiceListItem(serviceIndex, 'volumes', { target: '', content: '' })
                                      }
                                    />
                                    {service.volumes.map((volume, volumeIndex) => (
                                      <div key={volumeIndex} className="grid grid-cols-[1fr_1fr_32px] gap-2">
                                        <GuideField label={ct('fields.target')}>
                                          <input
                                            className={inputBaseClass}
                                            value={volume.target}
                                            required
                                            placeholder={ct('placeholders.flagTarget')}
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'volumes',
                                                volumeIndex,
                                                'target',
                                                e.target.value
                                              )
                                            }
                                          />
                                          <GuideErrors
                                            errors={
                                              guideFieldErrors[`service.${serviceIndex}.volumes.${volumeIndex}.target`]
                                            }
                                          />
                                        </GuideField>
                                        <GuideField label={ct('fields.content')}>
                                          <input
                                            className={inputBaseClass}
                                            value={volume.content || ''}
                                            placeholder="uuid{}"
                                            onChange={(e) =>
                                              updateGuideServiceList(
                                                serviceIndex,
                                                'volumes',
                                                volumeIndex,
                                                'content',
                                                e.target.value
                                              )
                                            }
                                          />
                                        </GuideField>
                                        <IconButton
                                          onClick={() =>
                                            removeGuideServiceListItem(serviceIndex, 'volumes', volumeIndex)
                                          }
                                        />
                                      </div>
                                    ))}
                                  </>
                                ) : null}

                                <GuideListHeader
                                  title={ct('sections.networks')}
                                  addLabel={ct('actions.add')}
                                  disabled={guideConfig.networks.length === 0}
                                  onAdd={() =>
                                    addGuideServiceListItem(serviceIndex, 'networks', {
                                      name: '',
                                      ipv4Address: '',
                                      macAddress: '',
                                    })
                                  }
                                />
                                <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.networks`]} />
                                {service.networks.map((network, networkIndex) => (
                                  <div key={networkIndex} className="grid grid-cols-[1fr_1fr_1fr_32px] gap-2">
                                    <GuideField label={ct('fields.network')}>
                                      <select
                                        className={selectClass}
                                        value={network.name}
                                        required
                                        onChange={(e) =>
                                          updateGuideServiceList(
                                            serviceIndex,
                                            'networks',
                                            networkIndex,
                                            'name',
                                            e.target.value
                                          )
                                        }
                                      >
                                        <option value="">{ct('placeholders.network')}</option>
                                        {guideConfig.networks.map((item, index) => (
                                          <option key={index} value={item.name}>
                                            {item.name}
                                          </option>
                                        ))}
                                      </select>
                                      <GuideErrors
                                        errors={
                                          guideFieldErrors[`service.${serviceIndex}.networks.${networkIndex}.name`]
                                        }
                                      />
                                    </GuideField>
                                    <GuideField label={ct('fields.ipv4Address')}>
                                      <input
                                        className={inputBaseClass}
                                        value={network.ipv4Address}
                                        required
                                        placeholder={ct('placeholders.ipv4Address')}
                                        onChange={(e) =>
                                          updateGuideServiceList(
                                            serviceIndex,
                                            'networks',
                                            networkIndex,
                                            'ipv4Address',
                                            e.target.value
                                          )
                                        }
                                      />
                                      <GuideErrors
                                        errors={
                                          guideFieldErrors[
                                            `service.${serviceIndex}.networks.${networkIndex}.ipv4Address`
                                          ]
                                        }
                                      />
                                    </GuideField>
                                    <GuideField label={ct('fields.macAddress')}>
                                      <input
                                        className={inputBaseClass}
                                        value={network.macAddress || ''}
                                        required={service.kubeVirt}
                                        placeholder={ct('placeholders.macAddress')}
                                        onChange={(e) =>
                                          updateGuideServiceList(
                                            serviceIndex,
                                            'networks',
                                            networkIndex,
                                            'macAddress',
                                            e.target.value
                                          )
                                        }
                                      />
                                      <GuideErrors
                                        errors={
                                          guideFieldErrors[
                                            `service.${serviceIndex}.networks.${networkIndex}.macAddress`
                                          ]
                                        }
                                      />
                                    </GuideField>
                                    <IconButton
                                      onClick={() => removeGuideServiceListItem(serviceIndex, 'networks', networkIndex)}
                                    />
                                  </div>
                                ))}
                              </div>
                            ))}

                            <div className="flex justify-between items-center">
                              <h4 className="text-sm font-mono text-neutral-200">{ct('sections.networks')}</h4>
                              <Button
                                variant="ghost"
                                size="sm"
                                align="icon-left"
                                icon={<IconPlus size={14} />}
                                onClick={addGuideNetwork}
                              >
                                {ct('actions.addNetwork')}
                              </Button>
                            </div>
                            {guideConfig.networks.map((network, networkIndex) => (
                              <div key={networkIndex} className="grid grid-cols-[2fr_3fr_2fr_112px_32px] gap-2">
                                <GuideField label={ct('fields.nameRequired')}>
                                  <input
                                    className={inputBaseClass}
                                    value={network.name}
                                    required
                                    placeholder={ct('placeholders.nameRequired')}
                                    onChange={(e) => updateGuideNetwork(networkIndex, 'name', e.target.value)}
                                  />
                                  <GuideErrors errors={guideFieldErrors[`network.${networkIndex}.name`]} />
                                </GuideField>
                                <GuideField label={ct('fields.subnet')}>
                                  <input
                                    className={inputBaseClass}
                                    value={network.subnet}
                                    required
                                    placeholder={ct('placeholders.subnet')}
                                    onChange={(e) => updateGuideNetwork(networkIndex, 'subnet', e.target.value)}
                                  />
                                  <GuideErrors errors={guideFieldErrors[`network.${networkIndex}.subnet`]} />
                                </GuideField>
                                <GuideField label={ct('fields.gateway')}>
                                  <input
                                    className={inputBaseClass}
                                    value={network.gateway}
                                    required
                                    placeholder={ct('placeholders.gateway')}
                                    onChange={(e) => updateGuideNetwork(networkIndex, 'gateway', e.target.value)}
                                  />
                                  <GuideErrors errors={guideFieldErrors[`network.${networkIndex}.gateway`]} />
                                </GuideField>
                                <GuideField label={ct('fields.external')}>
                                  <select
                                    className={selectClass}
                                    value={network.external ? 'true' : 'false'}
                                    onChange={(e) =>
                                      updateGuideNetwork(networkIndex, 'external', e.target.value === 'true')
                                    }
                                  >
                                    <option value="true">true</option>
                                    <option value="false">false</option>
                                  </select>
                                </GuideField>
                                <IconButton onClick={() => removeGuideNetwork(networkIndex)} />
                              </div>
                            ))}
                          </div>

                          <div className="border border-neutral-300/30 rounded-md overflow-hidden bg-black/30">
                            <div className="px-3 py-2 border-b border-neutral-700 text-xs font-mono text-neutral-400">
                              docker-compose.yaml
                            </div>
                            <Suspense
                              fallback={
                                <div className="flex items-center justify-center h-[200px] text-neutral-400 font-mono text-sm">
                                  {ct('loadingEditor')}
                                </div>
                              }
                            >
                              <Editor
                                value={challenge.docker_compose || emptyComposeYaml}
                                onChange={updateDockerCompose}
                                language="yaml"
                                options={{
                                  readOnly: false,
                                  minimap: { enabled: false },
                                  scrollBeyondLastLine: false,
                                  scrollbar: {
                                    alwaysConsumeMouseWheel: false,
                                    vertical: 'auto',
                                    horizontal: 'auto',
                                  },
                                  lineNumbers: 'on',
                                  folding: true,
                                  wordWrap: 'on',
                                  fontSize: 14,
                                  fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
                                  tabSize: 2,
                                  insertSpaces: true,
                                  renderLineHighlight: 'line',
                                }}
                                height={`${19 * (challenge.docker_compose || emptyComposeYaml).split('\n').length}px`}
                                theme="vs-dark"
                              />
                            </Suspense>
                          </div>
                        </div>
                        <div className="mt-1 text-xs text-neutral-500 font-mono">
                          {yamlNoticeLines.map((line, index) => (
                            <span key={index}>
                              {line}
                              <br />
                            </span>
                          ))}
                        </div>
                        {rawValidation.list.length > 0 && (
                          <div className="mt-3 rounded-md border border-red-400/30 bg-red-400/10 p-3 text-xs font-mono text-red-200 space-y-1">
                            {rawValidation.list.map((error, index) => (
                              <div key={index}>{error}</div>
                            ))}
                          </div>
                        )}
                      </div>

                      {/* 网络策略 */}
                      <div className="mt-4">
                        <div className="flex justify-between items-center mb-2">
                          <label className="block text-sm font-mono text-neutral-400">
                            {t('admin.challengeModal.labels.networkPolicy')}
                          </label>
                          <Button
                            variant="ghost"
                            size="sm"
                            align="icon-left"
                            icon={<IconPlus size={14} />}
                            className="!bg-transparent !text-geek-400 hover:!text-geek-300"
                            onClick={addNetworkPolicy}
                          >
                            {t('admin.challengeModal.actions.addPolicy')}
                          </Button>
                        </div>

                        <div className="grid grid-cols-1 2xl:grid-cols-[minmax(0,1fr)_minmax(420px,0.95fr)] gap-4 items-start">
                          <div className="space-y-4">
                            {(challenge.network_policies || []).map((policy, policyIndex) => (
                              <div key={policyIndex} className="border border-neutral-700 rounded-md p-3 bg-black/20">
                                <div className="flex justify-between items-center mb-3">
                                  <h5 className="text-sm font-mono text-neutral-200">
                                    {t('admin.challengeModal.labels.policy', {
                                      index: policyIndex + 1,
                                    })}
                                  </h5>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    className="!bg-transparent !text-red-400 hover:!text-red-300"
                                    onClick={() => removeNetworkPolicy(policyIndex)}
                                  >
                                    <IconTrash size={16} />
                                  </Button>
                                </div>

                                {vpcMode && (
                                  <GuideField label={t('admin.challengeModal.labels.policyTarget')}>
                                    <select
                                      value={policy.service || ''}
                                      onChange={(e) => updatePolicyTarget(policyIndex, e.target.value)}
                                      className={`${selectClass} mb-3`}
                                    >
                                      <option value="">{t('admin.challengeModal.placeholders.policyTarget')}</option>
                                      {policyTargets.map((target) => (
                                        <option key={target.service} value={target.service}>
                                          {target.label}
                                        </option>
                                      ))}
                                    </select>
                                  </GuideField>
                                )}

                                {/* From 策略 */}
                                <div className="mb-3">
                                  <div className="flex justify-between items-center mb-2">
                                    <label className="text-xs font-mono text-neutral-400">
                                      {t('admin.challengeModal.labels.allowInbound')}
                                    </label>
                                    <Button
                                      variant="ghost"
                                      size="sm"
                                      align="icon-left"
                                      icon={<IconPlus size={12} />}
                                      className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                      onClick={() => addPolicyRule(policyIndex, 'from')}
                                    >
                                      {t('admin.challengeModal.actions.addRule')}
                                    </Button>
                                  </div>

                                  <div className="space-y-2">
                                    {(policy.from || []).map((fromRule, fromIndex) => (
                                      <div key={fromIndex} className="border border-neutral-700/50 p-2 rounded">
                                        <div className="flex justify-between items-center mb-2">
                                          <label className="text-xs font-mono text-neutral-400">CIDR</label>
                                          {policy.from.length > 1 && (
                                            <Button
                                              variant="ghost"
                                              size="icon"
                                              className="!bg-transparent !text-red-400 hover:!text-red-300"
                                              onClick={() => removePolicyRule(policyIndex, 'from', fromIndex)}
                                            >
                                              <IconTrash size={12} />
                                            </Button>
                                          )}
                                        </div>

                                        <input
                                          type="text"
                                          value={normalizePolicyRule(fromRule).cidr}
                                          onChange={(e) =>
                                            updatePolicyCidr(policyIndex, 'from', fromIndex, e.target.value)
                                          }
                                          className="w-full h-8 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                          placeholder="192.168.1.0/24"
                                        />

                                        <div className="mt-2">
                                          <div className="flex justify-between items-center mb-1">
                                            <label className="text-xs font-mono text-neutral-400">
                                              {t('admin.challengeModal.labels.excludeList')}
                                            </label>
                                            <Button
                                              variant="ghost"
                                              size="sm"
                                              align="icon-left"
                                              icon={<IconPlus size={10} />}
                                              className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                              onClick={() => addPolicyExcept(policyIndex, 'from', fromIndex)}
                                            >
                                              {t('admin.challengeModal.actions.addExclude')}
                                            </Button>
                                          </div>
                                          <div className="space-y-1">
                                            {normalizePolicyRule(fromRule).except.map((exceptItem, exceptIndex) => (
                                              <div key={exceptIndex} className="flex gap-1 items-center">
                                                <input
                                                  type="text"
                                                  value={exceptItem}
                                                  onChange={(e) =>
                                                    updatePolicyExcept(
                                                      policyIndex,
                                                      'from',
                                                      fromIndex,
                                                      exceptIndex,
                                                      e.target.value
                                                    )
                                                  }
                                                  className="flex-1 h-6 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                                  placeholder="192.168.1.1"
                                                />
                                                <Button
                                                  variant="ghost"
                                                  size="icon"
                                                  className="!bg-transparent !text-red-400 hover:!text-red-300 !w-6 !h-6"
                                                  onClick={() =>
                                                    removePolicyExcept(policyIndex, 'from', fromIndex, exceptIndex)
                                                  }
                                                >
                                                  <IconTrash size={10} />
                                                </Button>
                                              </div>
                                            ))}
                                          </div>
                                        </div>
                                      </div>
                                    ))}
                                  </div>
                                </div>

                                {/* To 策略 */}
                                <div>
                                  <div className="flex justify-between items-center mb-2">
                                    <label className="text-xs font-mono text-neutral-400">
                                      {t('admin.challengeModal.labels.allowOutbound')}
                                    </label>
                                    <Button
                                      variant="ghost"
                                      size="sm"
                                      align="icon-left"
                                      icon={<IconPlus size={12} />}
                                      className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                      onClick={() => addPolicyRule(policyIndex, 'to')}
                                    >
                                      {t('admin.challengeModal.actions.addRule')}
                                    </Button>
                                  </div>

                                  <div className="space-y-2">
                                    {(policy.to || []).map((toRule, toIndex) => (
                                      <div key={toIndex} className="border border-neutral-700/50 p-2 rounded">
                                        <div className="flex justify-between items-center mb-2">
                                          <label className="text-xs font-mono text-neutral-400">CIDR</label>
                                          {policy.to.length > 1 && (
                                            <Button
                                              variant="ghost"
                                              size="icon"
                                              className="!bg-transparent !text-red-400 hover:!text-red-300"
                                              onClick={() => removePolicyRule(policyIndex, 'to', toIndex)}
                                            >
                                              <IconTrash size={12} />
                                            </Button>
                                          )}
                                        </div>

                                        <input
                                          type="text"
                                          value={normalizePolicyRule(toRule).cidr}
                                          onChange={(e) => updatePolicyCidr(policyIndex, 'to', toIndex, e.target.value)}
                                          className="w-full h-8 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                          placeholder="0.0.0.0/0"
                                        />

                                        <div className="mt-2">
                                          <div className="flex justify-between items-center mb-1">
                                            <label className="text-xs font-mono text-neutral-400">
                                              {t('admin.challengeModal.labels.excludeList')}
                                            </label>
                                            <Button
                                              variant="ghost"
                                              size="sm"
                                              align="icon-left"
                                              icon={<IconPlus size={10} />}
                                              className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                              onClick={() => addPolicyExcept(policyIndex, 'to', toIndex)}
                                            >
                                              {t('admin.challengeModal.actions.addExclude')}
                                            </Button>
                                          </div>
                                          <div className="space-y-1">
                                            {normalizePolicyRule(toRule).except.map((exceptItem, exceptIndex) => (
                                              <div key={exceptIndex} className="flex gap-1 items-center">
                                                <input
                                                  type="text"
                                                  value={exceptItem}
                                                  onChange={(e) =>
                                                    updatePolicyExcept(
                                                      policyIndex,
                                                      'to',
                                                      toIndex,
                                                      exceptIndex,
                                                      e.target.value
                                                    )
                                                  }
                                                  className="flex-1 h-6 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                                  placeholder="10.0.0.0/8"
                                                />
                                                <Button
                                                  variant="ghost"
                                                  size="icon"
                                                  className="!bg-transparent !text-red-400 hover:!text-red-300 !w-6 !h-6"
                                                  onClick={() =>
                                                    removePolicyExcept(policyIndex, 'to', toIndex, exceptIndex)
                                                  }
                                                >
                                                  <IconTrash size={10} />
                                                </Button>
                                              </div>
                                            ))}
                                          </div>
                                        </div>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              </div>
                            ))}
                          </div>
                          <div className="2xl:sticky 2xl:top-0">
                            <NetworkTopologyPreview topology={networkTopology} t={t} />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* 底部按钮 */}
        <div className="flex shrink-0 justify-end gap-3 sm:gap-4 p-3 sm:p-4 border-t border-neutral-700 bg-neutral-950/50">
          <Button variant="ghost" size="sm" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button
            variant="primary"
            size="sm"
            disabled={challenge.type === 'pods' && (guideValidationErrors.length > 0 || rawValidation.list.length > 0)}
            onClick={() => onSubmit(challenge)}
          >
            {mode === 'add'
              ? t('admin.challengeModal.actions.add')
              : mode === 'edit'
                ? t('common.saveChanges')
                : t('admin.challengeModal.actions.confirmDelete')}
          </Button>
        </div>
      </motion.div>
    </div>
  );
}

export default AdminChallengeModal;
