export const createGuideWriteFile = () => ({
  path: '',
  content: '',
  owner: '',
  permissions: '',
  encoding: '',
  append: false,
  defer: false,
});

export const createGuideUser = () => ({
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

export const createGuideGroup = () => ({ name: '', members: [] });

export const createGuideCloudConfig = () => ({ users: [], groups: [], writeFiles: [], sshAuthorizedKeys: [] });

export const createGuideService = () => ({
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

export const createGuideNetwork = () => ({ name: '', subnet: '', gateway: '' });
export const createGuideConfig = () => ({ services: [], networks: [] });
export const emptyComposeYaml = 'services:';

export const hasVpcNetworks = (config) => (config.networks || []).some((network) => network.name.trim());

export const normalizeGuideConfigForMode = (config) => {
  if (hasVpcNetworks(config)) return config;
  return { ...config, services: (config.services || []).map((service) => ({ ...service, kubeVirt: false })) };
};

export const hasOpenPort = (config) =>
  (config.services || []).some(
    (service) => !service.kubeVirt && (service.ports || []).some((port) => port.target.trim())
  );

export const getGuideServiceTargets = (config) =>
  (config.services || []).map((service, index) => {
    const name = service.containerName.trim() || service.name.trim() || `service${index + 1}`;
    return { label: name, service: name };
  });
