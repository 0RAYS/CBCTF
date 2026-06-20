import { IconPlus, IconTrash } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '../../../../common';

const emptyServer = () => ({ host: '', port: 7000, token: '', allowed: [] });
const emptyAllowed = () => ({ from: 10000, to: 30000, exclude: [] });

const sanitizeNumber = (value, fallbackValue = 0) => {
  if (value === '' || value === null || value === undefined) {
    return fallbackValue;
  }
  const numeric = Number(value);
  return Number.isNaN(numeric) ? fallbackValue : numeric;
};

const parseExclude = (value) =>
  value
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
    .map(Number)
    .filter((item) => !Number.isNaN(item));

export function FrpServerField({ config, updateConfig }) {
  const { t } = useTranslation();
  const servers = config.k8s.frp.frps || [];

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <span className="text-xs font-mono text-neutral-400">{t('admin.system.labels.frps')}</span>
        <Button
          variant="outline"
          size="xs"
          icon={<IconPlus size={14} />}
          onClick={() =>
            updateConfig((draft) => {
              draft.k8s.frp.frps.push(emptyServer());
            })
          }
        >
          {t('common.add')}
        </Button>
      </div>
      <div className="space-y-3">
        {servers.map((server, serverIndex) => (
          <div key={serverIndex} className="rounded-lg border border-neutral-800 bg-neutral-950/60 p-3 space-y-3">
            <div className="grid gap-2 md:grid-cols-4">
              <Input
                size="sm"
                value={server.host || ''}
                placeholder={t('admin.system.labels.frpsHostname')}
                onChange={(event) =>
                  updateConfig((draft) => {
                    draft.k8s.frp.frps[serverIndex].host = event.target.value;
                  })
                }
              />
              <Input
                size="sm"
                type="number"
                value={server.port ?? 0}
                placeholder={t('admin.system.labels.frpsPort')}
                onChange={(event) =>
                  updateConfig((draft) => {
                    draft.k8s.frp.frps[serverIndex].port = sanitizeNumber(event.target.value, server.port);
                  })
                }
              />
              <Input
                size="sm"
                type="password"
                value={server.token || ''}
                placeholder={t('admin.system.labels.frpsToken')}
                onChange={(event) =>
                  updateConfig((draft) => {
                    draft.k8s.frp.frps[serverIndex].token = event.target.value;
                  })
                }
              />
              <Button
                variant="danger"
                size="sm"
                icon={<IconTrash size={14} />}
                onClick={() =>
                  updateConfig((draft) => {
                    draft.k8s.frp.frps.splice(serverIndex, 1);
                  })
                }
              >
                {t('common.delete')}
              </Button>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-xs text-neutral-500">{t('admin.system.k8s.allowedPorts')}</span>
                <Button
                  variant="outline"
                  size="xs"
                  icon={<IconPlus size={14} />}
                  onClick={() =>
                    updateConfig((draft) => {
                      draft.k8s.frp.frps[serverIndex].allowed.push(emptyAllowed());
                    })
                  }
                >
                  {t('common.add')}
                </Button>
              </div>
              {(server.allowed || []).map((allowed, allowedIndex) => (
                <div key={allowedIndex} className="grid gap-2 md:grid-cols-[1fr_1fr_2fr_auto]">
                  <Input
                    size="sm"
                    type="number"
                    value={allowed.from ?? 0}
                    onChange={(event) =>
                      updateConfig((draft) => {
                        draft.k8s.frp.frps[serverIndex].allowed[allowedIndex].from = sanitizeNumber(
                          event.target.value,
                          allowed.from
                        );
                      })
                    }
                  />
                  <Input
                    size="sm"
                    type="number"
                    value={allowed.to ?? 0}
                    onChange={(event) =>
                      updateConfig((draft) => {
                        draft.k8s.frp.frps[serverIndex].allowed[allowedIndex].to = sanitizeNumber(
                          event.target.value,
                          allowed.to
                        );
                      })
                    }
                  />
                  <Input
                    size="sm"
                    value={(allowed.exclude || []).join(',')}
                    placeholder="20000,20001"
                    onChange={(event) =>
                      updateConfig((draft) => {
                        draft.k8s.frp.frps[serverIndex].allowed[allowedIndex].exclude = parseExclude(
                          event.target.value
                        );
                      })
                    }
                  />
                  <Button
                    variant="danger"
                    size="sm"
                    icon={<IconTrash size={14} />}
                    onClick={() =>
                      updateConfig((draft) => {
                        draft.k8s.frp.frps[serverIndex].allowed.splice(allowedIndex, 1);
                      })
                    }
                  />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
