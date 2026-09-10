import { useTranslation } from 'react-i18next';
import { createGuideNetwork } from './guideModel.js';
import { GuideListHeader, GuideInput, IconButton } from './GuideFields.jsx';

export default function NetworkDefinitionsEditor({ config, errors, onChange }) {
  const { t } = useTranslation();
  const ct = (key) => t(`admin.challengeModal.composeGuide.${key}`);
  const update = (index, field, value) =>
    onChange({
      ...config,
      networks: config.networks.map((network, i) => (i === index ? { ...network, [field]: value } : network)),
    });
  const remove = (index) => {
    const removedName = config.networks[index]?.name;
    onChange({
      networks: config.networks.filter((_, i) => i !== index),
      services: config.services.map((service) => ({
        ...service,
        networks: service.networks.filter((network) => network.name !== removedName),
      })),
    });
  };
  return (
    <>
      <GuideListHeader
        title={ct('sections.networks')}
        addLabel={ct('actions.addNetwork')}
        onAdd={() => onChange({ ...config, networks: [...config.networks, createGuideNetwork()] })}
      />
      {config.networks.map((network, index) => (
        <div key={index} className="grid grid-cols-[2fr_3fr_2fr_32px] gap-2">
          {['name', 'subnet', 'gateway'].map((field) => {
            const label = field === 'name' ? 'nameRequired' : field;
            return (
              <GuideInput
                key={field}
                label={ct(`fields.${label}`)}
                value={network[field]}
                required
                placeholder={ct(`placeholders.${label}`)}
                errors={errors[`network.${index}.${field}`]}
                onChange={(value) => update(index, field, value)}
              />
            );
          })}
          <IconButton onClick={() => remove(index)} />
        </div>
      ))}
    </>
  );
}
