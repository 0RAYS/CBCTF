import { useTranslation } from 'react-i18next';
import Checkbox from '@/components/common/Checkbox';
import { createGuideRange } from '@/components/features/Admin/challenges/editor/guideModel.js';
import {
  GuideErrors,
  GuideField,
  GuideInput,
  GuideListHeader,
  IconButton,
} from '@/components/features/Admin/challenges/editor/GuideFields.jsx';
import { selectClass } from '@/components/features/Admin/challenges/editor/editorStyles.js';

const toNumber = (value) => {
  const parsed = Number.parseInt(value, 10);
  return Number.isNaN(parsed) ? 0 : parsed;
};

export default function RangeConfigEditor({ config, errors, onChange }) {
  const { t } = useTranslation();
  const ct = (key, options) => t(`admin.challengeModal.composeGuide.${key}`, options);
  const range = config.range;
  const services = config.services || [];
  const networks = (config.networks || []).filter((network) => network.name.trim());

  const updateRange = (patch) => onChange({ ...config, range: { ...range, ...patch } });
  const toggle = (enabled) =>
    onChange({
      ...config,
      range: enabled ? { ...createGuideRange(), gateway: services[0]?.name.trim() || '' } : null,
      services: services.map((service) => ({ ...service, role: enabled ? service.role || '' : '' })),
    });
  const updateRoute = (index, field, value) =>
    updateRange({
      routes: (range.routes || []).map((route, i) => (i === index ? { ...route, [field]: value } : route)),
    });

  return (
    <div className="space-y-2 border-t border-neutral-700/60 pt-3">
      <Checkbox checked={Boolean(range)} onChange={(e) => toggle(e.target.checked)} label={ct('range.enable')} />
      {range ? (
        <>
          <p className="text-[11px] font-mono text-neutral-500">{ct('range.hint')}</p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
            {[
              { field: 'poolSize', label: 'poolSize', placeholder: 'poolSize' },
              { field: 'maxLeased', label: 'maxLeased', placeholder: 'maxLeased' },
              { field: 'leaseDuration', label: 'leaseDuration', placeholder: 'leaseDuration' },
            ].map((item) => (
              <GuideInput
                key={item.field}
                label={ct(`fields.${item.label}`)}
                value={String(range[item.field] ?? 0)}
                placeholder={ct(`placeholders.${item.placeholder}`)}
                errors={errors[`range.${item.field}`]}
                onChange={(value) => updateRange({ [item.field]: toNumber(value) })}
              />
            ))}
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
            <GuideField label={ct('fields.gatewayService')}>
              <select
                className={selectClass}
                value={range.gateway || ''}
                onChange={(e) => updateRange({ gateway: e.target.value })}
              >
                <option value="">{ct('placeholders.gatewayService')}</option>
                {services.map((service, index) => (
                  <option key={index} value={service.name.trim()}>
                    {service.name.trim() || ct('serviceIndexed', { index: index + 1 })}
                  </option>
                ))}
              </select>
              <GuideErrors errors={errors['range.gateway']} />
            </GuideField>
            <GuideField label={ct('fields.accessNetwork')}>
              <select
                className={selectClass}
                value={range.accessNetwork || ''}
                onChange={(e) => updateRange({ accessNetwork: e.target.value })}
              >
                <option value="">{ct('placeholders.accessNetwork')}</option>
                {networks.map((network) => (
                  <option key={network.name} value={network.name.trim()}>
                    {network.name.trim()}
                  </option>
                ))}
              </select>
              <GuideErrors errors={errors['range.accessNetwork']} />
            </GuideField>
            <GuideInput
              label={ct('fields.sshUser')}
              value={range.sshUser || ''}
              placeholder={ct('placeholders.sshUser')}
              onChange={(value) => updateRange({ sshUser: value })}
            />
          </div>
          <GuideListHeader
            title={ct('sections.rangeRoutes')}
            addLabel={ct('actions.addRoute')}
            disabled={networks.length < 2}
            onAdd={() => updateRange({ routes: [...(range.routes || []), { from: '', to: '' }] })}
          />
          {(range.routes || []).map((route, index) => (
            <div key={index} className="grid grid-cols-[1fr_1fr_32px] gap-2">
              {['from', 'to'].map((field) => (
                <GuideField key={field} label={ct(`fields.route${field === 'from' ? 'From' : 'To'}`)}>
                  <select
                    className={selectClass}
                    value={route[field] || ''}
                    onChange={(e) => updateRoute(index, field, e.target.value)}
                  >
                    <option value="">{ct(`placeholders.route${field === 'from' ? 'From' : 'To'}`)}</option>
                    {networks.map((network) => (
                      <option key={network.name} value={network.name.trim()}>
                        {network.name.trim()}
                      </option>
                    ))}
                  </select>
                  <GuideErrors errors={errors[`range.routes.${index}.${field}`]} />
                </GuideField>
              ))}
              <IconButton onClick={() => updateRange({ routes: (range.routes || []).filter((_, i) => i !== index) })} />
            </div>
          ))}
        </>
      ) : null}
    </div>
  );
}
