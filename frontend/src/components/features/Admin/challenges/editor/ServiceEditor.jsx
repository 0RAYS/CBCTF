import { useTranslation } from 'react-i18next';
import { GuideField, GuideInput, GuideBoolean, GuideErrors, GuideListHeader, IconButton } from './GuideFields.jsx';
import { selectClass } from './editorStyles.js';
import CloudInitEditor from './CloudInitEditor.jsx';

export default function ServiceEditor({ service, serviceIndex, networks, vpcMode, errors, onChange, onRemove }) {
  const { t } = useTranslation();
  const ct = (key, options) => t(`admin.challengeModal.composeGuide.${key}`, options);
  const update = (field, value) => onChange({ ...service, [field]: value });
  const updateItem = (field, index, itemField, value) =>
    update(
      field,
      service[field].map((item, i) => (i === index ? { ...item, [itemField]: value } : item))
    );
  const addItem = (field, item) => update(field, [...service[field], item]);
  const removeItem = (field, index) =>
    update(
      field,
      service[field].filter((_, i) => i !== index)
    );
  const fieldErrors = (path) => errors[`service.${serviceIndex}.${path}`];

  return (
    <div className="border border-neutral-700 rounded-md p-3 bg-black/20 space-y-3">
      <div className="flex justify-between items-center">
        <span className="text-xs font-mono text-neutral-400">{ct('serviceIndexed', { index: serviceIndex + 1 })}</span>
        {onRemove && <IconButton onClick={onRemove} />}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {['name', 'containerName', 'image', 'cpus', 'memLimit'].map((field) => {
          const label = field === 'name' ? 'serviceName' : field;
          return (
            <GuideInput
              key={field}
              label={ct(`fields.${label}`)}
              value={service[field]}
              required={['name', 'image'].includes(field)}
              placeholder={ct(`placeholders.${label}`)}
              errors={fieldErrors(field)}
              onChange={(value) => update(field, value)}
            />
          );
        })}
        {!service.kubeVirt && (
          <>
            <GuideInput
              label={ct('fields.workingDir')}
              value={service.workingDir}
              placeholder={ct('placeholders.workingDir')}
              onChange={(value) => update('workingDir', value)}
            />
            <GuideInput
              multiline
              label={ct('fields.command')}
              value={service.command}
              placeholder={ct('placeholders.command')}
              onChange={(value) => update('command', value)}
            />
          </>
        )}
        {vpcMode && (
          <GuideBoolean
            label={ct('fields.kubeVirt')}
            value={service.kubeVirt}
            onChange={(value) => update('kubeVirt', value)}
          />
        )}
        {service.kubeVirt && (
          <>
            <GuideField label={ct('fields.bootloader')}>
              <select
                className={selectClass}
                value={service.bootloader}
                onChange={(e) => update('bootloader', e.target.value)}
              >
                <option value="">{ct('placeholders.bootloader')}</option>
                <option value="bios">bios</option>
                <option value="efi">efi</option>
              </select>
              <GuideErrors errors={fieldErrors('bootloader')} />
            </GuideField>
            <GuideBoolean
              label={ct('fields.secureBoot')}
              value={service.secureBoot}
              onChange={(value) => update('secureBoot', value)}
            />
          </>
        )}
      </div>

      {service.kubeVirt ? (
        <CloudInitEditor value={service.userData} onChange={(value) => update('userData', value)} />
      ) : (
        <>
          <GuideListHeader
            title={ct('sections.ports')}
            addLabel={ct('actions.add')}
            onAdd={() => addItem('ports', { published: '', target: '', protocol: '' })}
          />
          {service.ports.map((port, index) => (
            <div key={index} className="grid grid-cols-[1fr_1fr_90px_32px] gap-2">
              <GuideInput
                label={ct('fields.name')}
                value={port.published}
                placeholder={ct('placeholders.name')}
                errors={fieldErrors(`ports.${index}.published`)}
                onChange={(value) => updateItem('ports', index, 'published', value)}
              />
              <GuideInput
                label={ct('fields.target')}
                value={port.target}
                required
                placeholder={ct('placeholders.target')}
                errors={fieldErrors(`ports.${index}.target`)}
                onChange={(value) => updateItem('ports', index, 'target', value)}
              />
              <GuideField label={ct('fields.protocol')}>
                <select
                  className={selectClass}
                  value={port.protocol}
                  onChange={(e) => updateItem('ports', index, 'protocol', e.target.value)}
                >
                  <option value="">{ct('fields.protocol')}</option>
                  <option value="tcp">tcp</option>
                  <option value="udp">udp</option>
                </select>
              </GuideField>
              <IconButton onClick={() => removeItem('ports', index)} />
            </div>
          ))}

          <GuideListHeader
            title={ct('sections.environment')}
            addLabel={ct('actions.add')}
            onAdd={() => addItem('environment', { key: '', value: '' })}
          />
          {service.environment.map((env, index) => (
            <div key={index} className="grid grid-cols-[1fr_1fr_32px] gap-2">
              <GuideInput
                label={ct('fields.key')}
                value={env.key}
                required
                placeholder={ct('placeholders.key')}
                errors={fieldErrors(`environment.${index}.key`)}
                onChange={(value) => updateItem('environment', index, 'key', value)}
              />
              <GuideInput
                label={ct('fields.value')}
                value={env.value}
                placeholder={ct('placeholders.value')}
                onChange={(value) => updateItem('environment', index, 'value', value)}
              />
              <IconButton onClick={() => removeItem('environment', index)} />
            </div>
          ))}

          <GuideListHeader
            title={ct('sections.fileFlags')}
            addLabel={ct('actions.add')}
            onAdd={() => addItem('volumes', { target: '', content: '' })}
          />
          {service.volumes.map((volume, index) => (
            <div key={index} className="grid grid-cols-[1fr_1fr_32px] gap-2">
              <GuideInput
                label={ct('fields.target')}
                value={volume.target}
                required
                placeholder={ct('placeholders.flagTarget')}
                errors={fieldErrors(`volumes.${index}.target`)}
                onChange={(value) => updateItem('volumes', index, 'target', value)}
              />
              <GuideInput
                label={ct('fields.content')}
                value={volume.content || ''}
                placeholder="uuid{}"
                onChange={(value) => updateItem('volumes', index, 'content', value)}
              />
              <IconButton onClick={() => removeItem('volumes', index)} />
            </div>
          ))}
        </>
      )}

      <GuideListHeader
        title={ct('sections.networks')}
        addLabel={ct('actions.add')}
        disabled={networks.length === 0}
        onAdd={() => addItem('networks', { name: '', ipv4Address: '', macAddress: '' })}
      />
      <GuideErrors errors={fieldErrors('networks')} />
      {service.networks.map((network, index) => (
        <div key={index} className="grid grid-cols-[1fr_1fr_1fr_32px] gap-2">
          <GuideField label={ct('fields.network')}>
            <select
              className={selectClass}
              value={network.name}
              required
              onChange={(e) => updateItem('networks', index, 'name', e.target.value)}
            >
              <option value="">{ct('placeholders.network')}</option>
              {networks.map((item, i) => (
                <option key={i} value={item.name}>
                  {item.name}
                </option>
              ))}
            </select>
            <GuideErrors errors={fieldErrors(`networks.${index}.name`)} />
          </GuideField>
          <GuideInput
            label={ct('fields.ipv4Address')}
            value={network.ipv4Address}
            required
            placeholder={ct('placeholders.ipv4Address')}
            errors={fieldErrors(`networks.${index}.ipv4Address`)}
            onChange={(value) => updateItem('networks', index, 'ipv4Address', value)}
          />
          <GuideInput
            label={ct('fields.macAddress')}
            value={network.macAddress || ''}
            required={service.kubeVirt}
            placeholder={ct('placeholders.macAddress')}
            errors={fieldErrors(`networks.${index}.macAddress`)}
            onChange={(value) => updateItem('networks', index, 'macAddress', value)}
          />
          <IconButton onClick={() => removeItem('networks', index)} />
        </div>
      ))}
    </div>
  );
}
