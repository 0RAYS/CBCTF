import { useTranslation } from 'react-i18next';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import Button from '../../../../common/Button';
import { GuideField } from './GuideFields.jsx';
import { selectClass } from './editorStyles.js';
import { defaultNetworkPolicy, normalizeNetworkPolicy, normalizePolicyRule } from './networkPolicy.js';

export default function NetworkPolicyEditor({ policies, vpcMode, targets, onChange }) {
  const { t } = useTranslation();
  const updatePolicy = (index, policy) => onChange(policies.map((item, i) => (i === index ? policy : item)));
  const updateRules = (index, direction, rules) =>
    updatePolicy(index, {
      ...normalizeNetworkPolicy(policies[index]),
      [direction]: rules,
    });
  const updateRule = (index, direction, ruleIndex, change) => {
    const policy = normalizeNetworkPolicy(policies[index]);
    updateRules(
      index,
      direction,
      policy[direction].map((rule, i) => (i === ruleIndex ? change(rule) : rule))
    );
  };
  return (
    <div className="space-y-4">
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
          onClick={() => onChange([...policies, defaultNetworkPolicy(vpcMode ? targets[0] : {})])}
        >
          {t('admin.challengeModal.actions.addPolicy')}
        </Button>
      </div>
      {policies.map((policy, policyIndex) => (
        <div key={policyIndex} className="border border-neutral-700 rounded-md p-3 bg-black/20">
          <div className="flex justify-between items-center mb-3">
            <h5 className="text-sm font-mono text-neutral-200">
              {t('admin.challengeModal.labels.policy', { index: policyIndex + 1 })}
            </h5>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-red-400 hover:!text-red-300"
              onClick={() => onChange(policies.filter((_, i) => i !== policyIndex))}
            >
              <IconTrash size={16} />
            </Button>
          </div>
          {vpcMode && (
            <GuideField label={t('admin.challengeModal.labels.policyTarget')}>
              <select
                value={policy.service || ''}
                className={`${selectClass} mb-3`}
                onChange={(e) => updatePolicy(policyIndex, { ...policy, service: e.target.value })}
              >
                <option value="">{t('admin.challengeModal.placeholders.policyTarget')}</option>
                {targets.map((target) => (
                  <option key={target.service} value={target.service}>
                    {target.label}
                  </option>
                ))}
              </select>
            </GuideField>
          )}
          {['from', 'to'].map((direction) => (
            <div key={direction} className={direction === 'from' ? 'mb-3' : ''}>
              <div className="flex justify-between items-center mb-2">
                <label className="text-xs font-mono text-neutral-400">
                  {t(`admin.challengeModal.labels.${direction === 'from' ? 'allowInbound' : 'allowOutbound'}`)}
                </label>
                <Button
                  variant="ghost"
                  size="sm"
                  align="icon-left"
                  icon={<IconPlus size={12} />}
                  className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                  onClick={() =>
                    updateRules(policyIndex, direction, [
                      ...normalizeNetworkPolicy(policy)[direction],
                      { cidr: '', except: [''] },
                    ])
                  }
                >
                  {t('admin.challengeModal.actions.addRule')}
                </Button>
              </div>
              <div className="space-y-2">
                {(policy[direction] || []).map((rawRule, ruleIndex) => {
                  const rule = normalizePolicyRule(rawRule);
                  return (
                    <div key={ruleIndex} className="border border-neutral-700/50 p-2 rounded">
                      <div className="flex justify-between items-center mb-2">
                        <label className="text-xs font-mono text-neutral-400">CIDR</label>
                        {policy[direction].length > 1 && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                            onClick={() =>
                              updateRules(
                                policyIndex,
                                direction,
                                normalizeNetworkPolicy(policy)[direction].filter((_, i) => i !== ruleIndex)
                              )
                            }
                          >
                            <IconTrash size={12} />
                          </Button>
                        )}
                      </div>
                      <input
                        type="text"
                        value={rule.cidr}
                        onChange={(e) =>
                          updateRule(policyIndex, direction, ruleIndex, (item) => ({ ...item, cidr: e.target.value }))
                        }
                        className="w-full h-8 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                        placeholder={direction === 'from' ? '192.168.1.0/24' : '0.0.0.0/0'}
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
                            onClick={() =>
                              updateRule(policyIndex, direction, ruleIndex, (item) => ({
                                ...item,
                                except: [...item.except, ''],
                              }))
                            }
                          >
                            {t('admin.challengeModal.actions.addExclude')}
                          </Button>
                        </div>
                        <div className="space-y-1">
                          {rule.except.map((except, exceptIndex) => (
                            <div key={exceptIndex} className="flex gap-1 items-center">
                              <input
                                type="text"
                                value={except}
                                onChange={(e) =>
                                  updateRule(policyIndex, direction, ruleIndex, (item) => ({
                                    ...item,
                                    except: item.except.map((value, i) => (i === exceptIndex ? e.target.value : value)),
                                  }))
                                }
                                className="flex-1 h-6 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                placeholder={direction === 'from' ? '192.168.1.1' : '10.0.0.0/8'}
                              />
                              <Button
                                variant="ghost"
                                size="icon"
                                className="!bg-transparent !text-red-400 hover:!text-red-300 !w-6 !h-6"
                                onClick={() =>
                                  updateRule(policyIndex, direction, ruleIndex, (item) => ({
                                    ...item,
                                    except: item.except.filter((_, i) => i !== exceptIndex),
                                  }))
                                }
                              >
                                <IconTrash size={10} />
                              </Button>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}
