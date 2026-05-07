import { motion } from 'motion/react';
import { IconX, IconPlus, IconTrash } from '@tabler/icons-react';
import { Button } from '../../../components/common';
import { lazy, Suspense, useState, useEffect } from 'react';
const Editor = lazy(() => import('@monaco-editor/react'));
import { useTranslation } from 'react-i18next';

/**
 * 题目管理弹窗组件
 * @param {Object} props
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

const buildGuidedComposeYaml = (config) => {
  const lines = ['services:'];
  const flagVolumes = new Map();
  const networks = (config.networks || []).filter((network) => network.name.trim());

  (config.services || []).forEach((service, index) => {
    const serviceName = service.name.trim() || `service${index + 1}`;
    lines.push(`  ${serviceName}:`);
    if (service.containerName.trim()) lines.push(`    container_name: ${service.containerName.trim()}`);
    lines.push(`    image: ${service.image.trim()}`);
    if (service.cpus.trim()) lines.push(`    cpus: ${service.cpus.trim()}`);
    if (service.memLimit.trim()) lines.push(`    mem_limit: ${service.memLimit.trim()}`);

    const ports = (service.ports || []).filter((port) => port.target.trim());
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

    const envs = (service.environment || []).filter((env) => env.key.trim());
    if (envs.length > 0) {
      lines.push('    environment:');
      envs.forEach((env) => lines.push(`      - ${yamlQuote(`${env.key.trim()}=${env.value}`)}`));
    }

    const volumes = (service.volumes || []).filter((volume) => volume.source.trim() && volume.target.trim());
    if (volumes.length > 0) {
      lines.push('    volumes:');
      volumes.forEach((volume) => {
        const source = volume.source.trim();
        lines.push('      - type: volume', `        source: ${source}`, `        target: ${volume.target.trim()}`);
        if (source.startsWith('FLAG')) {
          flagVolumes.set(source, volume.value || 'uuid{}');
        }
      });
    }

    if (service.workingDir.trim()) lines.push(`    working_dir: ${service.workingDir.trim()}`);

    const command = String(service.command || '').split('\n');
    if (command.some((item) => item.trim())) {
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
        });
    }
  });

  if (flagVolumes.size > 0) {
    lines.push('', 'volumes:');
    flagVolumes.forEach((value, source) => {
      lines.push(`  ${source}:`, '    labels:', `      - value=${value}`);
    });
  }

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

const parsePortString = (value) => {
  const [main, protocol = 'tcp'] = stripYamlValue(value).split('/');
  const parts = main.split(':');
  if (parts.length === 1) return { published: '', target: parts[0] || '', protocol };
  return { published: parts[0] || '', target: parts[1] || '', protocol };
};

const parseVolumeString = (value) => {
  const parts = stripYamlValue(value).split(':');
  return { source: parts[0] || '', target: parts[1] || '', value: '' };
};

const parseComposeYamlToGuideConfig = (yaml, t = (key, values) => `${key}${values?.line ? ` ${values.line}` : ''}`) => {
  const lines = String(yaml || '')
    .replace(/\r\n/g, '\n')
    .split('\n')
    .map((raw, index) => ({ raw, index: index + 1, text: raw.replace(/#.*$/, '') }))
    .filter((line) => line.text.trim());
  const errors = [];
  const services = [];
  const networks = [];
  const volumeValues = new Map();
  let section = '';
  let service = null;
  let network = null;
  let volumeName = '';
  let serviceList = '';
  let serviceNetwork = null;
  let currentPort = null;
  let currentVolume = null;

  for (const line of lines) {
    if (/\t/.test(line.raw))
      errors.push(t('admin.challengeModal.composeGuide.validation.noTabs', { line: line.index }));
    const indent = line.text.match(/^ */)[0].length;
    const text = line.text.trim();

    if (indent === 0) {
      section = text.replace(/:$/, '');
      service = null;
      network = null;
      volumeName = '';
      serviceList = '';
      serviceNetwork = null;
      currentPort = null;
      currentVolume = null;
      if (!['services', 'volumes', 'networks'].includes(section)) {
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
        if (['ports', 'environment', 'volumes', 'command', 'networks'].includes(pair.key) && pair.value === '') {
          serviceList = pair.key;
          continue;
        }
        serviceList = '';
        if (pair.key === 'container_name') service.containerName = pair.value;
        else if (pair.key === 'image') service.image = pair.value;
        else if (pair.key === 'cpus') service.cpus = pair.value;
        else if (pair.key === 'mem_limit') service.memLimit = pair.value;
        else if (pair.key === 'working_dir') service.workingDir = pair.value;
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
      if (indent === 6 && serviceList === 'volumes' && text.startsWith('- ')) {
        const value = text.slice(2).trim();
        if (value.includes(':')) {
          const pair = parseKeyValueLine(value);
          if (pair && ['type', 'source', 'target'].includes(pair.key)) {
            currentVolume = { source: '', target: '', value: '' };
            service.volumes.push(currentVolume);
            if (pair.key === 'type' && pair.value !== 'volume') {
              errors.push(t('admin.challengeModal.composeGuide.validation.volumeTypeInvalid', { line: line.index }));
            } else if (pair.key === 'source') currentVolume.source = pair.value;
            else if (pair.key === 'target') currentVolume.target = pair.value;
          } else {
            service.volumes.push(parseVolumeString(value));
            currentVolume = null;
          }
        } else {
          service.volumes.push(parseVolumeString(value));
          currentVolume = null;
        }
        continue;
      }
      if (indent === 8 && serviceList === 'volumes' && currentVolume) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'type' && pair.value !== 'volume') {
          errors.push(t('admin.challengeModal.composeGuide.validation.volumeTypeInvalid', { line: line.index }));
        } else if (pair?.key === 'source') currentVolume.source = pair.value;
        else if (pair?.key === 'target') currentVolume.target = pair.value;
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
        serviceNetwork = { name: text.slice(0, -1), ipv4Address: '' };
        service.networks.push(serviceNetwork);
        continue;
      }
      if (indent === 8 && serviceList === 'networks' && serviceNetwork) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'ipv4_address') serviceNetwork.ipv4Address = pair.value;
        continue;
      }
      errors.push(t('admin.challengeModal.composeGuide.validation.unparseableLine', { line: line.index }));
    } else if (section === 'volumes') {
      if (indent === 2 && text.endsWith(':')) {
        volumeName = text.slice(0, -1);
      } else if (indent === 6 && text.startsWith('- value=')) {
        volumeValues.set(volumeName, text.slice('- value='.length));
      }
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

  services.forEach((item) => {
    item.volumes = item.volumes.map((volume) => ({ ...volume, value: volumeValues.get(volume.source) || '' }));
  });

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
  const fileFlagSources = new Set();
  const networkCidrs = new Map(definedNetworks.map((network) => [network.name.trim(), network.subnet.trim()]));
  const networkAssignedIps = new Map();

  definedNetworks.forEach((network) => {
    const name = network.name.trim();
    const gateway = network.gateway.trim();
    if (!name || !validateIp4(gateway)) return;
    networkAssignedIps.set(name, new Set([gateway]));
  });

  if (!config.services?.length) addError('services', t('admin.challengeModal.composeGuide.validation.serviceRequired'));
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
    const servicePortTargets = new Set();
    (service.ports || []).forEach((port, portIndex) => {
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
    (service.environment || []).forEach((env, envIndex) => {
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
    (service.volumes || []).forEach((volume, volumeIndex) => {
      const source = volume.source.trim();
      const target = volume.target.trim();
      if (!volume.source.trim())
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.source`,
          t('admin.challengeModal.composeGuide.validation.volumeSourceRequired', { label, index: volumeIndex + 1 })
        );
      if (!volume.target.trim())
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.target`,
          t('admin.challengeModal.composeGuide.validation.volumeTargetRequired', { label, index: volumeIndex + 1 })
        );
      if (source && !source.startsWith('FLAG')) {
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.source`,
          t('admin.challengeModal.composeGuide.validation.volumeSourcePrefix', { label, index: volumeIndex + 1 })
        );
      }
      if (source && fileFlagSources.has(source))
        addError(
          `service.${serviceIndex}.volumes.${volumeIndex}.source`,
          t('admin.challengeModal.composeGuide.validation.volumeSourceUnique', { label, index: volumeIndex + 1 })
        );
      if (source) fileFlagSources.add(source);
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

  const guideValidation = validateGuidedCompose(guideConfig, t);
  const guideValidationErrors = guideValidation.list;
  const guideFieldErrors = guideValidation.fields;
  const rawValidation = validateRawCompose(challenge.docker_compose || emptyComposeYaml, t);

  const syncGuideConfig = (nextConfig) => {
    setGuideConfig(nextConfig);
    onChange({ ...challenge, docker_compose: buildGuidedComposeYaml(nextConfig) });
  };

  useEffect(() => {
    if (challenge.type !== 'pods' || !challenge.docker_compose) return;
    const parsed = parseComposeYamlToGuideConfig(challenge.docker_compose, t);
    if (parsed.ok) setGuideConfig(parsed.config);
  }, [challenge.type, challenge.docker_compose]);

  // docker_compose 更新
  const updateDockerCompose = (value) => {
    const finalValue = value || '';
    const parsed = parseComposeYamlToGuideConfig(finalValue, t);
    if (parsed.ok) setGuideConfig(parsed.config);
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
      {
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
      },
    ];
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  const removeNetworkPolicy = (policyIndex) => {
    const newNetworkPolicies = challenge.network_policies.filter((_, i) => i !== policyIndex);
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 from/to 规则
  const addPolicyRule = (policyIndex, ruleType) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };

    policy[ruleType] = [
      ...policy[ruleType],
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
    const policy = { ...newNetworkPolicies[policyIndex] };

    policy[ruleType] = policy[ruleType].filter((_, i) => i !== ruleIndex);

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 cidr 值
  const updatePolicyCidr = (policyIndex, ruleType, ruleIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.cidr = value;

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 except 规则
  const addPolicyExcept = (policyIndex, ruleType, ruleIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = [...rule.except, ''];

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 移除 except 规则
  const removePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = rule.except.filter((_, i) => i !== exceptIndex);

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 except 值
  const updatePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = [...rule.except];
    rule.except[exceptIndex] = value;

    policy[ruleType] = [...policy[ruleType]];
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

  const podNoticeLines = [
    t('admin.challengeModal.podsNotice.flagFormat', {
      format: '`static{}`, `dynamic{}`, `uuid{}`',
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
      format: '`static{}`, `dynamic{}`, `uuid{}`',
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
        className="w-full max-w-7xl bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* 标题栏 */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
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
        <div className="p-6 max-h-[70vh] overflow-y-auto">
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
            <div className="space-y-4">
              {/* 基本信息 */}
              <div className="mb-4">
                <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.basic')}</h3>

                <div className="space-y-3">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
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
                <div className="border-t border-neutral-700 pt-4">
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
                <div className="border-t border-neutral-700 pt-4">
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
                      {challenge.type === 'static' ? 'static{}' : 'dynamic{} / uuid{}'}
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
                <div className="border-t border-neutral-700 pt-4">
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
                <div className="border-t border-neutral-700 pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">
                      {t('admin.challengeModal.sections.containers')}
                    </h3>
                  </div>

                  <div className="space-y-6">
                    <div className="border border-neutral-700 rounded-md p-4 space-y-4">
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
                                  <GuideField label={ct('fields.workingDir')}>
                                    <input
                                      className={inputBaseClass}
                                      value={service.workingDir}
                                      placeholder={ct('placeholders.workingDir')}
                                      onChange={(e) => updateGuideService(serviceIndex, 'workingDir', e.target.value)}
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
                                </div>

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
                                        errors={guideFieldErrors[`service.${serviceIndex}.ports.${portIndex}.target`]}
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
                                        errors={guideFieldErrors[`service.${serviceIndex}.environment.${envIndex}.key`]}
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
                                      onClick={() => removeGuideServiceListItem(serviceIndex, 'environment', envIndex)}
                                    />
                                  </div>
                                ))}

                                <GuideListHeader
                                  title={ct('sections.fileFlags')}
                                  addLabel={ct('actions.add')}
                                  onAdd={() =>
                                    addGuideServiceListItem(serviceIndex, 'volumes', {
                                      source: '',
                                      target: '',
                                      value: '',
                                    })
                                  }
                                />
                                {service.volumes.map((volume, volumeIndex) => (
                                  <div key={volumeIndex} className="grid grid-cols-[1fr_1fr_1fr_32px] gap-2">
                                    <GuideField label={ct('fields.source')}>
                                      <input
                                        className={inputBaseClass}
                                        value={volume.source}
                                        required
                                        placeholder={ct('placeholders.source')}
                                        onChange={(e) =>
                                          updateGuideServiceList(
                                            serviceIndex,
                                            'volumes',
                                            volumeIndex,
                                            'source',
                                            e.target.value
                                          )
                                        }
                                      />
                                      <GuideErrors
                                        errors={
                                          guideFieldErrors[`service.${serviceIndex}.volumes.${volumeIndex}.source`]
                                        }
                                      />
                                    </GuideField>
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
                                    <GuideField label={ct('fields.value')}>
                                      <input
                                        className={inputBaseClass}
                                        value={volume.value}
                                        placeholder="uuid{}"
                                        onChange={(e) =>
                                          updateGuideServiceList(
                                            serviceIndex,
                                            'volumes',
                                            volumeIndex,
                                            'value',
                                            e.target.value
                                          )
                                        }
                                      />
                                    </GuideField>
                                    <IconButton
                                      onClick={() => removeGuideServiceListItem(serviceIndex, 'volumes', volumeIndex)}
                                    />
                                  </div>
                                ))}

                                <GuideListHeader
                                  title={ct('sections.networks')}
                                  addLabel={ct('actions.add')}
                                  disabled={guideConfig.networks.length === 0}
                                  onAdd={() =>
                                    addGuideServiceListItem(serviceIndex, 'networks', { name: '', ipv4Address: '' })
                                  }
                                />
                                <GuideErrors errors={guideFieldErrors[`service.${serviceIndex}.networks`]} />
                                {service.networks.map((network, networkIndex) => (
                                  <div key={networkIndex} className="grid grid-cols-[1fr_1fr_32px] gap-2">
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
                              <div key={networkIndex} className="grid grid-cols-[1fr_1fr_1fr_88px_32px] gap-2">
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
                                        value={fromRule.cidr}
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
                                          {(fromRule.except || []).map((exceptItem, exceptIndex) => (
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
                                        value={toRule.cidr}
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
                                          {(toRule.except || []).map((exceptItem, exceptIndex) => (
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
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* 底部按钮 */}
        <div className="flex justify-end gap-4 p-4 border-t border-neutral-700">
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
