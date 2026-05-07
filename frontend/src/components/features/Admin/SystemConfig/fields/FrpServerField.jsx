import { IconPlus, IconX } from '@tabler/icons-react';
import { Input, Button } from '../../../../common';
import { useTranslation } from 'react-i18next';

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

/**
 * FrpServerField - Specialized component for FRP server configuration
 * Handles complex nested structure: frps -> allowed -> exclude
 * @param {Array} frpsList - Array of FRP server configs
 * @param {Function} updateConfig - Immer draft updater function
 */
export function FrpServerField({ frpsList = [], updateConfig }) {
  const { t } = useTranslation();

  const addFrpsServer = () => {
    updateConfig((draft) => {
      if (!Array.isArray(draft.k8s.frp.frps)) {
        draft.k8s.frp.frps = [];
      }
      draft.k8s.frp.frps.push({ host: '', port: 0, token: '', allowed: [] });
    });
  };

  const removeFrpsServer = (frpsIndex) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps.splice(frpsIndex, 1);
    });
  };

  const updateFrpsField = (frpsIndex, field, value) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps[frpsIndex][field] = value;
    });
  };

  const addAllowedRange = (frpsIndex) => {
    updateConfig((draft) => {
      if (!Array.isArray(draft.k8s.frp.frps[frpsIndex].allowed)) {
        draft.k8s.frp.frps[frpsIndex].allowed = [];
      }
      draft.k8s.frp.frps[frpsIndex].allowed.push({ from: 0, to: 0, exclude: [] });
    });
  };

  const removeAllowedRange = (frpsIndex, allowedIndex) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps[frpsIndex].allowed.splice(allowedIndex, 1);
    });
  };

  const updateAllowedField = (frpsIndex, allowedIndex, field, value) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex][field] = value;
    });
  };

  const addExcludePort = (frpsIndex, allowedIndex) => {
    updateConfig((draft) => {
      if (!Array.isArray(draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex].exclude)) {
        draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex].exclude = [];
      }
      draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex].exclude.push(0);
    });
  };

  const removeExcludePort = (frpsIndex, allowedIndex, excludeIndex) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex].exclude.splice(excludeIndex, 1);
    });
  };

  const updateExcludePort = (frpsIndex, allowedIndex, excludeIndex, value) => {
    updateConfig((draft) => {
      draft.k8s.frp.frps[frpsIndex].allowed[allowedIndex].exclude[excludeIndex] = value;
    });
  };

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between">
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.frps')}</span>
        <Button size="icon" variant="ghost" aria-label={t('common.add')} onClick={addFrpsServer}>
          <IconPlus size={14} />
        </Button>
      </div>
      <div className="space-y-2">
        {frpsList.map((frps, frpsIndex) => (
          <div key={`frps-${frpsIndex}`} className="space-y-1 border border-neutral-300/10 rounded p-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-mono text-neutral-400">
                {t('admin.system.k8s.host')} #{frpsIndex + 1} {t('admin.system.labels.frpsHostname')}
              </span>
              <Button
                size="icon"
                variant="ghost"
                aria-label={t('common.remove')}
                onClick={() => removeFrpsServer(frpsIndex)}
              >
                <IconX size={14} />
              </Button>
            </div>
            <div className="grid gap-2">
              <Input
                size="sm"
                value={frps.host || ''}
                placeholder={t('admin.system.k8s.host')}
                onChange={(event) => updateFrpsField(frpsIndex, 'host', event.target.value)}
              />
              <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.frpsPort')}</span>
              <Input
                size="sm"
                type="number"
                value={frps.port || 0}
                placeholder={t('admin.system.k8s.port')}
                onChange={(event) =>
                  updateFrpsField(frpsIndex, 'port', sanitizeNumber(event.target.value, frps.port || 0))
                }
              />
              <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.frpsToken')}</span>
              <Input
                size="sm"
                value={frps.token || ''}
                placeholder={t('common.leaveBlankToKeep')}
                onChange={(event) => updateFrpsField(frpsIndex, 'token', event.target.value)}
              />
            </div>
            <div className="space-y-1">
              <div className="flex items-center justify-between">
                <span className="text-xs font-mono text-neutral-400">{t('admin.system.k8s.allowedPorts')}</span>
                <Button
                  size="icon"
                  variant="ghost"
                  aria-label={t('common.add')}
                  onClick={() => addAllowedRange(frpsIndex)}
                >
                  <IconPlus size={14} />
                </Button>
              </div>
              <div className="space-y-1">
                {(frps.allowed || []).map((allowed, allowedIndex) => (
                  <div key={`allowed-${frpsIndex}-${allowedIndex}`} className="flex flex-col gap-2">
                    <span className="text-xs font-mono text-neutral-400">From → To</span>
                    <div className="flex items-center gap-2">
                      <Input
                        size="sm"
                        type="number"
                        value={allowed.from ?? 0}
                        placeholder="from"
                        onChange={(event) =>
                          updateAllowedField(
                            frpsIndex,
                            allowedIndex,
                            'from',
                            sanitizeNumber(event.target.value, allowed.from ?? 0)
                          )
                        }
                      />
                      <Input
                        size="sm"
                        type="number"
                        value={allowed.to ?? 0}
                        placeholder="to"
                        onChange={(event) =>
                          updateAllowedField(
                            frpsIndex,
                            allowedIndex,
                            'to',
                            sanitizeNumber(event.target.value, allowed.to ?? 0)
                          )
                        }
                      />
                      <Button
                        size="icon"
                        variant="ghost"
                        aria-label={t('common.remove')}
                        onClick={() => removeAllowedRange(frpsIndex, allowedIndex)}
                      >
                        <IconX size={14} />
                      </Button>
                    </div>
                    <div className="space-y-1">
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-mono text-neutral-400">Exclude</span>
                        <Button
                          size="icon"
                          variant="ghost"
                          aria-label={t('common.add')}
                          onClick={() => addExcludePort(frpsIndex, allowedIndex)}
                        >
                          <IconPlus size={14} />
                        </Button>
                      </div>
                      {(allowed.exclude || []).map((excludeItem, excludeIndex) => (
                        <div key={`exclude-${excludeIndex}`} className="flex items-center gap-2">
                          <Input
                            size="sm"
                            type="number"
                            value={excludeItem ?? 0}
                            onChange={(event) =>
                              updateExcludePort(
                                frpsIndex,
                                allowedIndex,
                                excludeIndex,
                                sanitizeNumber(event.target.value, excludeItem ?? 0)
                              )
                            }
                          />
                          <Button
                            size="icon"
                            variant="ghost"
                            aria-label={t('common.remove')}
                            onClick={() => removeExcludePort(frpsIndex, allowedIndex, excludeIndex)}
                          >
                            <IconX size={14} />
                          </Button>
                        </div>
                      ))}
                    </div>
                    <div className="border-b border-neutral-300/20" />
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
