import { lazy, Suspense } from 'react';
import { useTranslation } from 'react-i18next';
import { createGuideService, emptyComposeYaml, getGuideServiceTargets, hasVpcNetworks } from './guideModel.js';
import { buildNetworkTopology } from './networkPolicy.js';
import { GuideListHeader } from './GuideFields.jsx';
import ServiceEditor from './ServiceEditor.jsx';
import NetworkDefinitionsEditor from './NetworkDefinitionsEditor.jsx';
import NetworkPolicyEditor from './NetworkPolicyEditor.jsx';
import TopologyPreview from './TopologyPreview.jsx';

const Editor = lazy(() => import('../../../../../lib/monacoSetup').then(() => import('@monaco-editor/react')));

export default function ComposeEditor({ challenge, onChange, compose }) {
  const { t } = useTranslation();
  const ct = (key) => t(`admin.challengeModal.composeGuide.${key}`);
  const { guideConfig, guideValidation, rawValidation, syncGuideConfig, updateDockerCompose } = compose;
  const vpcMode = hasVpcNetworks(guideConfig);
  const yamlNoticeLines = [
    t('admin.challengeModal.yamlNotice.flagFormat', { format: '`static{}`, `leet{}`, `uuid{}`' }),
    ...[
      'flagPrefix',
      'flagVolume',
      'dockerParams',
      'containerUnique',
      'networkIpRequired',
      'exposeNetwork',
      'noNetworkMerge',
      'containerNetwork',
    ].map((key) => t(`admin.challengeModal.yamlNotice.${key}`)),
  ];

  return (
    <div className="border-t border-neutral-700 pt-3 lg:pt-4">
      <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.containers')}</h3>
      <div className="border border-neutral-700 rounded-md p-3 lg:p-4 space-y-4 bg-black/10">
        <div>
          <div className="flex flex-wrap justify-between items-center gap-3 mb-3">
            <label className="block text-sm font-mono text-neutral-400">docker-compose.yaml</label>
          </div>
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
            <div className="space-y-4">
              <GuideListHeader
                title={ct('sections.services')}
                addLabel={ct('actions.addService')}
                onAdd={() =>
                  syncGuideConfig({ ...guideConfig, services: [...guideConfig.services, createGuideService()] })
                }
              />
              {guideConfig.services.map((service, serviceIndex) => (
                <ServiceEditor
                  key={serviceIndex}
                  service={service}
                  serviceIndex={serviceIndex}
                  networks={guideConfig.networks}
                  vpcMode={vpcMode}
                  errors={guideValidation.fields}
                  onChange={(next) =>
                    syncGuideConfig({
                      ...guideConfig,
                      services: guideConfig.services.map((item, i) => (i === serviceIndex ? next : item)),
                    })
                  }
                  onRemove={
                    guideConfig.services.length > 1
                      ? () =>
                          syncGuideConfig({
                            ...guideConfig,
                            services: guideConfig.services.filter((_, i) => i !== serviceIndex),
                          })
                      : undefined
                  }
                />
              ))}
              <NetworkDefinitionsEditor
                config={guideConfig}
                errors={guideValidation.fields}
                onChange={syncGuideConfig}
              />
            </div>
            <div className="border border-neutral-300/30 rounded-md overflow-hidden bg-black/30">
              <div className="px-3 py-2 border-b border-neutral-700 text-xs font-mono text-neutral-400">
                docker-compose.yaml
              </div>
              <Suspense
                fallback={
                  <div className="flex items-center justify-center h-[200px] text-neutral-400 font-mono text-sm">
                    {ct('loadingEditor')}
                  </div>
                }
              >
                <Editor
                  value={challenge.docker_compose || emptyComposeYaml}
                  onChange={updateDockerCompose}
                  language="yaml"
                  options={{
                    readOnly: false,
                    minimap: { enabled: false },
                    scrollBeyondLastLine: false,
                    scrollbar: { alwaysConsumeMouseWheel: false, vertical: 'auto', horizontal: 'auto' },
                    lineNumbers: 'on',
                    folding: true,
                    wordWrap: 'on',
                    fontSize: 14,
                    fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
                    tabSize: 2,
                    insertSpaces: true,
                    renderLineHighlight: 'line',
                  }}
                  height={`${19 * (challenge.docker_compose || emptyComposeYaml).split('\n').length}px`}
                  theme="vs-dark"
                />
              </Suspense>
            </div>
          </div>
          <div className="mt-1 text-xs text-neutral-500 font-mono">
            {yamlNoticeLines.map((line, index) => (
              <span key={index}>
                {line}
                <br />
              </span>
            ))}
          </div>
          {rawValidation.list.length > 0 && (
            <div className="mt-3 rounded-md border border-red-400/30 bg-red-400/10 p-3 text-xs font-mono text-red-200 space-y-1">
              {rawValidation.list.map((error, index) => (
                <div key={index}>{error}</div>
              ))}
            </div>
          )}
        </div>
        <div className="mt-4 grid grid-cols-1 2xl:grid-cols-[minmax(0,1fr)_minmax(420px,0.95fr)] gap-4 items-start">
          <NetworkPolicyEditor
            policies={challenge.network_policies || []}
            vpcMode={vpcMode}
            targets={getGuideServiceTargets(guideConfig)}
            onChange={(policies) => onChange({ ...challenge, network_policies: policies })}
          />
          <div className="2xl:sticky 2xl:top-0">
            <TopologyPreview topology={buildNetworkTopology(guideConfig, challenge.network_policies || [])} />
          </div>
        </div>
      </div>
    </div>
  );
}
