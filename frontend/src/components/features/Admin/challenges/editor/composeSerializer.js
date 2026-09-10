import { createGuideCloudConfig } from './guideModel.js';

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
const hasGuideWriteFile = (file) =>
  file.path.trim() ||
  file.content.trim() ||
  file.owner.trim() ||
  file.permissions.trim() ||
  file.encoding.trim() ||
  file.append ||
  file.defer;

const isGuideCloudConfigEmpty = (cloudConfig = {}) =>
  !(cloudConfig.users || []).some(hasGuideUser) &&
  !(cloudConfig.groups || []).some(hasGuideGroup) &&
  !(cloudConfig.writeFiles || []).some(hasGuideWriteFile) &&
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
  const writeFiles = (cloudConfig.writeFiles || []).filter(hasGuideWriteFile);
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

export const buildGuidedComposeYaml = (config) => {
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
    if (kubeVirt) lines.push('    x-kubevirt: true');
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
