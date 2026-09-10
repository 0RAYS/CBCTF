import { createGuideService, createGuideUser, createGuideGroup, createGuideWriteFile } from './guideModel.js';

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
  return { key: line.slice(0, index).trim(), value: stripYamlValue(line.slice(index + 1)) };
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

// This is the guide's supported Compose subset, not a general YAML parser.
// Raw editor text is retained separately and must never be replaced merely by parsing it.
export const parseComposeYamlToGuideConfig = (
  yaml,
  t = (key, values) => `${key}${values?.line ? ` ${values.line}` : ''}`
) => {
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
        service = { ...createGuideService(), name: text.slice(0, -1) };
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
        network = { name: text.slice(0, -1), subnet: '', gateway: '' };
        networks.push(network);
      } else if (network && indent === 4) {
        const pair = parseKeyValueLine(text);
        if (pair?.key === 'external' && !['true', 'false'].includes(pair.value))
          errors.push(t('admin.challengeModal.composeGuide.validation.networkExternalInvalid', { line: line.index }));
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
