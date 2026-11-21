import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconChevronDown, IconChevronRight, IconPlus, IconTrash } from '@tabler/icons-react';
import Button from '@/components/common/Button';
import Checkbox from '@/components/common/Checkbox';
import { GuideField } from '@/components/features/Admin/challenges/editor/GuideFields.jsx';
import { inputBaseClass } from '@/components/features/Admin/challenges/editor/editorStyles.js';
import {
  defaultNetworkPolicy,
  normalizeNetworkPolicy,
  resolvePolicies,
  resolveSharedPolicy,
} from '@/components/features/Admin/challenges/editor/networkPolicy.js';

const nt = (key) => `admin.challengeModal.networkPolicy.${key}`;

const serialize = (policy) => ({
  service: policy.service || '',
  allow_internet: policy.allow_internet === true,
  egress_cidrs: policy.egress_cidrs || [],
  allow_external: policy.allow_external === true,
  peers: policy.peers || [],
  from: policy.from || [],
  to: policy.to || [],
});

export default function NetworkPolicyEditor({ policies, vpcMode, targets, context, onChange }) {
  const { t } = useTranslation();
  const shared = !vpcMode;
  const resolved = resolvePolicies(policies, targets, context);
  const sharedPolicy = shared ? resolveSharedPolicy(policies) : null;

  const write = (entries) => onChange(entries.map(serialize));

  // applyEntries 基于当前策略列表做一次写入，避免多次 onChange 相互覆盖。
  const applyEntries = (mutate) => write(mutate(policies.map(normalizeNetworkPolicy)));

  const updateService = (service, changes) =>
    applyEntries((entries) => {
      const effective = resolved.find((item) => item.service === service) || defaultNetworkPolicy(service);
      const next = { ...effective, ...changes };
      const index = entries.findIndex((entry) => entry.service === service);
      if (index >= 0) entries[index] = next;
      else entries.push(next);
      return entries;
    });

  const updateShared = (changes) => write([{ ...sharedPolicy, ...changes }]);

  const togglePeer = (service, peer, enabled) =>
    applyEntries((entries) => {
      const peersOf = (name) => {
        const peers = new Set(resolved.find((item) => item.service === name)?.peers || []);
        const other = name === service ? peer : service;
        if (enabled) peers.add(other);
        else peers.delete(other);
        return [...peers].filter((item) => item !== name).sort();
      };
      [service, peer].forEach((name) => {
        const next = {
          ...(resolved.find((item) => item.service === name) || defaultNetworkPolicy(name)),
          peers: peersOf(name),
        };
        const index = entries.findIndex((entry) => entry.service === name);
        if (index >= 0) entries[index] = next;
        else entries.push(next);
      });
      return entries;
    });

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-2">
        <label className="block text-sm font-mono text-neutral-400">{t(nt('title'))}</label>
        <span className="text-[11px] font-mono text-neutral-500">
          {t(shared ? nt('sharedMode.tag') : nt('vpcMode.tag'))}
        </span>
      </div>
      <div className="rounded-md border border-neutral-700/60 bg-black/20 p-3 text-[11px] font-mono text-neutral-400 space-y-1">
        <p>{t(shared ? nt('sharedMode.hint') : nt('vpcMode.hint'))}</p>
        <p className="text-neutral-500">{t(nt('egress.scopeHint'))}</p>
      </div>
      {shared ? (
        <PolicyCard
          title={t(nt('sharedMode.cardTitle'))}
          subtitle={targets.map((target) => target.label).join(', ')}
          policy={sharedPolicy}
          showPeers={false}
          peerOptions={[]}
          onChange={(changes) => updateShared(changes)}
          onPeerToggle={() => {}}
        />
      ) : (
        targets.map((target) => {
          const policy = resolved.find((item) => item.service === target.service);
          if (!policy) return null;
          return (
            <PolicyCard
              key={target.service}
              title={target.label}
              policy={policy}
              showPeers={targets.length > 1}
              peerOptions={targets
                .filter((item) => item.service !== target.service)
                .map((item) => ({
                  service: item.service,
                  label: item.label,
                  checked: policy.peers.includes(item.service),
                  locked: policy.subnet_peers.includes(item.service),
                }))}
              onChange={(changes) => updateService(target.service, changes)}
              onPeerToggle={(peer, enabled) => togglePeer(target.service, peer, enabled)}
            />
          );
        })
      )}
    </div>
  );
}

function PolicyCard({ title, subtitle, policy, showPeers, peerOptions, onChange, onPeerToggle }) {
  const { t } = useTranslation();
  const [advancedOpen, setAdvancedOpen] = useState(false);
  return (
    <div className="border border-neutral-700 rounded-md p-3 bg-black/20 space-y-3">
      <div>
        <h5 className="text-sm font-mono text-neutral-200">{title}</h5>
        {subtitle && <div className="truncate text-[11px] font-mono text-neutral-500">{subtitle}</div>}
      </div>
      <div className="space-y-2">
        <Checkbox
          checked={policy.allow_internet}
          onChange={(event) => onChange({ allow_internet: event.target.checked })}
          label={t(nt('egress.internet'))}
        />
        {policy.allow_internet && (
          <div className="pl-6 text-[11px] font-mono text-neutral-500">{t(nt('egress.publicScope'))}</div>
        )}
      </div>
      <div className="space-y-2">
        <div className="flex justify-between items-center">
          <label className="text-[11px] font-mono text-neutral-500">{t(nt('egress.ranges'))}</label>
          <Button
            variant="ghost"
            size="sm"
            align="icon-left"
            icon={<IconPlus size={12} />}
            className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs !h-7"
            onClick={() => onChange({ egress_cidrs: [...policy.egress_cidrs, ''] })}
          >
            {t(nt('actions.addRange'))}
          </Button>
        </div>
        {policy.egress_cidrs.length === 0 ? (
          <div className="pl-1 text-[11px] font-mono text-neutral-600">{t(nt('egress.noRanges'))}</div>
        ) : (
          <div className="space-y-1">
            {policy.egress_cidrs.map((cidr, index) => (
              <div key={index} className="flex gap-1 items-center">
                <input
                  className={`${inputBaseClass} !h-8 flex-1 px-2 text-xs font-mono`}
                  value={cidr}
                  placeholder="10.0.0.0/16"
                  onChange={(event) =>
                    onChange({
                      egress_cidrs: policy.egress_cidrs.map((item, i) => (i === index ? event.target.value : item)),
                    })
                  }
                />
                <Button
                  variant="ghost"
                  size="icon"
                  className="!bg-transparent !text-red-400 hover:!text-red-300 !w-7 !h-7"
                  onClick={() => onChange({ egress_cidrs: policy.egress_cidrs.filter((_, i) => i !== index) })}
                >
                  <IconTrash size={12} />
                </Button>
              </div>
            ))}
          </div>
        )}
      </div>
      {showPeers && (
        <div className="space-y-2">
          <label className="text-[11px] font-mono text-neutral-500">{t(nt('peers.label'))}</label>
          <div className="flex flex-wrap gap-x-4 gap-y-1">
            {peerOptions.map((option) => (
              <span key={option.service} className="flex items-center gap-1">
                <Checkbox
                  checked={option.checked || option.locked}
                  disabled={option.locked}
                  onChange={(event) => onPeerToggle(option.service, event.target.checked)}
                  label={option.label}
                  className="text-xs font-mono"
                />
                {option.locked && (
                  <span className="text-[10px] font-mono text-neutral-600">{t(nt('peers.sameSubnet'))}</span>
                )}
              </span>
            ))}
          </div>
        </div>
      )}
      <Checkbox
        checked={policy.allow_external}
        onChange={(event) => onChange({ allow_external: event.target.checked })}
        label={t(nt('ingress.external'))}
      />
      {policy.allow_external && (
        <div className="pl-6 text-[11px] font-mono text-neutral-500">{t(nt('ingress.externalHint'))}</div>
      )}
      <div>
        <button
          type="button"
          className="flex items-center gap-1 text-[11px] font-mono text-neutral-500 hover:text-neutral-300"
          onClick={() => setAdvancedOpen((open) => !open)}
        >
          {advancedOpen ? <IconChevronDown size={12} /> : <IconChevronRight size={12} />}
          {t(nt('advanced.title'))}
        </button>
        {advancedOpen && (
          <div className="mt-2 space-y-3 border-l border-neutral-700 pl-3">
            {['from', 'to'].map((direction) => (
              <LegacyRules
                key={direction}
                direction={direction}
                rules={policy[direction]}
                onChange={(rules) => onChange({ [direction]: rules })}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function LegacyRules({ direction, rules, onChange }) {
  const { t } = useTranslation();
  const updateRule = (index, change) => onChange(rules.map((rule, i) => (i === index ? { ...rule, ...change } : rule)));
  return (
    <div className="space-y-2">
      <div className="flex justify-between items-center">
        <label className="text-[11px] font-mono text-neutral-500">
          {t(nt(direction === 'from' ? 'advanced.inbound' : 'advanced.outbound'))}
        </label>
        <Button
          variant="ghost"
          size="sm"
          align="icon-left"
          icon={<IconPlus size={12} />}
          className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs !h-7"
          onClick={() => onChange([...rules, { cidr: '', except: [] }])}
        >
          {t(nt('actions.addRule'))}
        </Button>
      </div>
      {rules.map((rule, index) => (
        <div key={index} className="space-y-1">
          <div className="flex gap-1 items-center">
            <input
              className={`${inputBaseClass} !h-8 flex-1 px-2 text-xs font-mono`}
              value={rule.cidr}
              placeholder={direction === 'from' ? '192.168.1.0/24' : '0.0.0.0/0'}
              onChange={(event) => updateRule(index, { cidr: event.target.value })}
            />
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-red-400 hover:!text-red-300 !w-7 !h-7"
              onClick={() => onChange(rules.filter((_, i) => i !== index))}
            >
              <IconTrash size={12} />
            </Button>
          </div>
          <GuideField label={t(nt('advanced.except'))}>
            <input
              className={`${inputBaseClass} !h-7 px-2 text-xs font-mono`}
              value={rule.except.join(', ')}
              placeholder="10.0.0.0/8"
              onChange={(event) =>
                updateRule(index, {
                  except: event.target.value
                    .split(',')
                    .map((item) => item.trim())
                    .filter(Boolean),
                })
              }
            />
          </GuideField>
        </div>
      ))}
    </div>
  );
}
