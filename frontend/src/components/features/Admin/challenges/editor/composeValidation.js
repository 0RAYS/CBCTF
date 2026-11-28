import { hasOpenPort, RANGE_ROLES } from '@/components/features/Admin/challenges/editor/guideModel.js';
import { parseComposeYamlToGuideConfig } from '@/components/features/Admin/challenges/editor/composeParser.js';

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

export const validateRawCompose = (yaml, t) => {
  const parsed = parseComposeYamlToGuideConfig(yaml, t);
  if (!parsed.ok) return { list: parsed.errors, fields: {}, config: null };
  const validation = validateGuidedCompose(parsed.config, t);
  return { ...validation, config: parsed.config };
};

const validationKey = (key) => `admin.challengeModal.composeGuide.validation.${key}`;

// 靶场拓扑校验，规则与后端 parseRangeTopology 一致：网关、接入网段与互通声明都必须在题目里成立。
const validateRangeTopology = (config, t, networkNames) => {
  const errors = [];
  const fields = {};
  const addError = (path, message) => {
    errors.push(message);
    fields[path] = [...(fields[path] || []), message];
  };
  const services = config.services || [];
  const range = config.range;

  services.forEach((service, serviceIndex) => {
    const role = (service.role || '').trim();
    if (!role) return;
    const label =
      service.name.trim() || t('admin.challengeModal.composeGuide.serviceIndexed', { index: serviceIndex + 1 });
    if (!range) {
      addError(`service.${serviceIndex}.role`, t(validationKey('rangeRoleRequiresRange'), { label }));
      return;
    }
    if (!RANGE_ROLES.includes(role)) {
      addError(`service.${serviceIndex}.role`, t(validationKey('rangeRoleUnknown'), { label, role }));
      return;
    }
    if (role === 'gateway' && service.name.trim() !== (range.gateway || '').trim()) {
      addError(
        `service.${serviceIndex}.role`,
        t(validationKey('rangeRoleGatewayMismatch'), { label, gateway: (range.gateway || '').trim() })
      );
    }
  });
  if (!range) return { errors, fields };

  if (!Number.isFinite(range.poolSize) || range.poolSize < 0) {
    addError('range.poolSize', t(validationKey('rangePoolSizeInvalid')));
  }
  if (!Number.isFinite(range.maxLeased) || range.maxLeased < 0) {
    addError('range.maxLeased', t(validationKey('rangeMaxLeasedInvalid')));
  }
  if (!Number.isFinite(range.leaseDuration) || range.leaseDuration <= 0) {
    addError('range.leaseDuration', t(validationKey('rangeLeaseDurationInvalid')));
  }

  const gatewayName = (range.gateway || '').trim();
  const gateway = services.find((service) => service.name.trim() === gatewayName);
  if (!gatewayName) {
    addError('range.gateway', t(validationKey('rangeGatewayRequired')));
  } else if (!gateway) {
    addError('range.gateway', t(validationKey('rangeGatewayUndefined'), { gateway: gatewayName }));
  } else if (gateway.kubeVirt) {
    addError('range.gateway', t(validationKey('rangeGatewayVirtual'), { gateway: gatewayName }));
  }

  const accessNetwork = (range.accessNetwork || '').trim();
  const gatewayNetworks = gateway ? (gateway.networks || []).map((network) => network.name.trim()) : [];
  if (!accessNetwork) {
    if (gatewayNetworks.length !== 1) {
      addError('range.accessNetwork', t(validationKey('rangeAccessNetworkRequired')));
    }
  } else {
    if (gateway && !gatewayNetworks.includes(accessNetwork)) {
      addError(
        'range.accessNetwork',
        t(validationKey('rangeAccessNetworkNotAttached'), { gateway: gatewayName, network: accessNetwork })
      );
    }
    if (!networkNames.has(accessNetwork)) {
      addError('range.accessNetwork', t(validationKey('rangeAccessNetworkUndefined'), { network: accessNetwork }));
    }
  }

  const seenRoutes = new Set();
  (range.routes || []).forEach((route, routeIndex) => {
    const ends = [
      { field: 'from', name: (route.from || '').trim() },
      { field: 'to', name: (route.to || '').trim() },
    ];
    ends.forEach((end) => {
      if (!end.name) {
        addError(
          `range.routes.${routeIndex}.${end.field}`,
          t(validationKey('rangeRouteNetworkRequired'), { index: routeIndex + 1 })
        );
        return;
      }
      if (!networkNames.has(end.name)) {
        addError(
          `range.routes.${routeIndex}.${end.field}`,
          t(validationKey('rangeRouteNetworkUndefined'), { network: end.name })
        );
      }
    });
    const [from, to] = ends.map((end) => end.name);
    if (from && to && from === to) {
      addError(`range.routes.${routeIndex}.to`, t(validationKey('rangeRouteSelfLoop'), { from, to }));
    }
    if (from && to) {
      const key = [from, to].sort().join('|');
      if (seenRoutes.has(key)) {
        addError(`range.routes.${routeIndex}.to`, t(validationKey('rangeRouteDuplicate'), { from, to }));
      }
      seenRoutes.add(key);
    }
  });

  // 选手从接入网段进拓扑,接入网段不与任何网段互通时选手够不到别的网段上的节点。
  const accessRouted = (range.routes || []).some((route) => {
    const from = (route.from || '').trim();
    const to = (route.to || '').trim();
    return from === accessNetwork || to === accessNetwork;
  });
  if (accessNetwork && !accessRouted) {
    const nodeOnOtherNetwork = (config.services || []).some(
      (service) =>
        service.name.trim() !== gatewayName &&
        (service.networks || []).some((network) => network.name.trim() !== accessNetwork)
    );
    if (nodeOnOtherNetwork) {
      addError('range.accessNetwork', t(validationKey('rangeAccessNetworkNotRouted'), { network: accessNetwork }));
    }
  }

  return { errors, fields };
};

export const validateGuidedCompose = (config, t = (key) => key) => {
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
  const rangeGatewayName = (config.range?.gateway || '').trim();
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
  // 靶场题的对外入口由玩家网关承担，节点不必暴露端口。
  if (config.services?.some((service) => !service.kubeVirt) && !hasOpenPort(config) && !config.range) {
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
    // 玩家网关是平台固定的跳板容器,镜像由平台提供,题目里不用写。
    if (!(rangeGatewayName && name === rangeGatewayName) && !service.image.trim())
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
    const fixedFlagNames = new Set();
    (service.fixedFlags || []).forEach((flag, flagIndex) => {
      const flagName = (flag.name || '').trim();
      const flagValue = (flag.value || '').trim();
      if (!flagName)
        addError(
          `service.${serviceIndex}.fixedFlags.${flagIndex}.name`,
          t(validationKey('fixedFlagNameRequired'), { label, index: flagIndex + 1 })
        );
      if (flagName && fixedFlagNames.has(flagName))
        addError(
          `service.${serviceIndex}.fixedFlags.${flagIndex}.name`,
          t(validationKey('fixedFlagNameUnique'), { label, index: flagIndex + 1 })
        );
      if (flagName) fixedFlagNames.add(flagName);
      if (!flagValue)
        addError(
          `service.${serviceIndex}.fixedFlags.${flagIndex}.value`,
          t(validationKey('fixedFlagValueRequired'), { label, index: flagIndex + 1 })
        );
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
  const rangeValidation = validateRangeTopology(config, t, networkNames);
  rangeValidation.errors.forEach((message) => errors.push(message));
  Object.entries(rangeValidation.fields).forEach(([path, messages]) => {
    fields[path] = [...(fields[path] || []), ...messages];
  });
  return { list: errors, fields };
};
