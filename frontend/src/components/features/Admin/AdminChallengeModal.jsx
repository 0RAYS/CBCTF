import { motion } from 'motion/react';
import {
  IconX,
  IconPlus,
  IconTrash,
  IconMaximize,
  IconMinimize,
  IconDeviceFloppy,
  IconCheck,
} from '@tabler/icons-react';
import { Button } from '../../../components/common';
import { lazy, Suspense, useState, useEffect } from 'react';
const Editor = lazy(() => import('@monaco-editor/react'));
import { useTranslation } from 'react-i18next';

/**
 * 题目管理弹窗组件
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示弹窗
 * @param {string} props.mode - 模式，'add'或'edit'
 * @param {Object} props.challenge - 当前编辑的题目对象
 * @param {Array} props.categories - 分类列表
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Function} props.onSubmit - 提交表单回调
 * @param {Function} props.onChange - 表单字段变更回调
 * @param {Function} props.onAddFlag - 添加Flag回调
 * @param {Function} props.onRemoveFlag - 移除Flag回调
 * @param {Function} props.onFlagChange - Flag变更回调
 * @param {Function} props.onAddOption - 添加选项回调
 * @param {Function} props.onRemoveOption - 移除选项回调
 * @param {Function} props.onOptionChange - 选项变更回调
 * @param {Function} props.onCorrectOptionChange - 正确选项变更回调
 */

const defaultYAML = `
services:
  service1:
    container_name: service1
    image: nicolaka/netshoot:latest
    ports:
      - "80:80/tcp"
    cpus: 0.5
    mem_limit: 500m
    environment:
      - FLAG=static{static_flag}
    volumes:
      - FLAG_0:/flag
    working_dir: /root
    command:
      - "sleep"
      - "infinity"
    networks:
      network1:
        ipv4_address: 192.168.0.2

  service2:
    container_name: service2
    image: nicolaka/netshoot:latest
    cpus: 0.5
    mem_limit: 500m
    environment:
      - FLAG=dynamic{dynamic_flag}
    volumes:
      - FLAG_1:/flag
    working_dir: /root
    command:
      - "sleep"
      - "infinity"
    networks:
      network1:
        ipv4_address: 192.168.0.3
      network2:
        ipv4_address: 192.168.1.3

  service3:
    container_name: service3
    image: nicolaka/netshoot:latest
    cpus: 0.5
    mem_limit: 500m
    environment:
      - FLAG=uuid{}
    volumes:
      - FLAG_2:/flag
    working_dir: /root
    command:
      - "sleep"
      - "infinity"
    networks:
      network2:
        ipv4_address: 192.168.1.4
      network3:
        ipv4_address: 192.168.2.4

volumes:
  FLAG_0:
    labels:
      - value=static{static_flag}
  FLAG_1:
    labels:
      - value=dynamic{dynamic_flag}
  FLAG_2:
    labels:
      - value=uuid{}

networks:
  network1:
    external: true
    ipam:
      config:
        - subnet: 192.168.0.0/24
          gateway: 192.168.0.1

  network2:
    ipam:
      config:
        - subnet: 192.168.1.0/24
          gateway: 192.168.1.1

  network3:
    ipam:
      config:
        - subnet: 192.168.2.0/24
          gateway: 192.168.2.1
`.trim();

function AdminChallengeModal({
  isOpen = false,
  mode = 'add',
  challenge = {
    name: '',
    description: '',
    category: '',
    type: 'static',
    generator_image: '',
    flags: [''],
    options: [],
    docker_compose: '',
    network_policies: [
      {
        from: [
          {
            cidr: '',
            except: [''],
          },
        ],
        to: [
          {
            cidr: '',
            except: [''],
          },
        ],
      },
    ],
  },
  categories = [],
  onClose,
  onSubmit,
  onChange,
  onAddFlag,
  onRemoveFlag,
  onFlagChange,
  onAddOption,
  onRemoveOption,
  onOptionChange,
  onCorrectOptionChange,
}) {
  const { t } = useTranslation();

  // 全屏编辑器状态
  const [fullscreenEditor, setFullscreenEditor] = useState({
    isOpen: false,
    value: '',
  });

  // 全屏编辑器控制
  const openFullscreenEditor = () => {
    setFullscreenEditor({
      isOpen: true,
      value: challenge.docker_compose || defaultYAML,
    });
  };

  // 初始化时确保docker_compose有默认值
  useEffect(() => {
    if (challenge.type === 'pods' && !challenge.docker_compose) {
      onChange({ ...challenge, docker_compose: defaultYAML });
    }
  }, [challenge.type, challenge.docker_compose, onChange]);

  const closeFullscreenEditor = () => {
    setFullscreenEditor({
      isOpen: false,
      value: '',
    });
  };

  const saveFullscreenEditor = () => {
    // 如果值为空，使用默认模板
    const finalValue = fullscreenEditor.value || defaultYAML;
    onChange({ ...challenge, docker_compose: finalValue });
    closeFullscreenEditor();
  };

  const updateFullscreenValue = (value) => {
    setFullscreenEditor({
      ...fullscreenEditor,
      value,
    });
  };

  // docker_compose 更新
  const updateDockerCompose = (value) => {
    // 如果值为空，使用默认模板
    const finalValue = value || defaultYAML;
    onChange({ ...challenge, docker_compose: finalValue });
  };

  // 网络策略操作
  const addNetworkPolicy = () => {
    const newNetworkPolicies = [
      ...(challenge.network_policies || []),
      {
        from: [
          {
            cidr: '',
            except: [''],
          },
        ],
        to: [
          {
            cidr: '0.0.0.0/0',
            except: ['10.0.0.0/8', '172.16.0.0/12', '192.168.0.0/16', '100.64.0.0/10'],
          },
        ],
      },
    ];
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  const removeNetworkPolicy = (policyIndex) => {
    const newNetworkPolicies = challenge.network_policies.filter((_, i) => i !== policyIndex);
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 from/to 规则
  const addPolicyRule = (policyIndex, ruleType) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };

    policy[ruleType] = [
      ...policy[ruleType],
      {
        cidr: '',
        except: [''],
      },
    ];

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 移除 from/to 规则
  const removePolicyRule = (policyIndex, ruleType, ruleIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };

    policy[ruleType] = policy[ruleType].filter((_, i) => i !== ruleIndex);

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 cidr 值
  const updatePolicyCidr = (policyIndex, ruleType, ruleIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.cidr = value;

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加 except 规则
  const addPolicyExcept = (policyIndex, ruleType, ruleIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = [...rule.except, ''];

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 移除 except 规则
  const removePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = rule.except.filter((_, i) => i !== exceptIndex);

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 更新 except 值
  const updatePolicyExcept = (policyIndex, ruleType, ruleIndex, exceptIndex, value) => {
    const newNetworkPolicies = [...challenge.network_policies];
    const policy = { ...newNetworkPolicies[policyIndex] };
    const rule = { ...policy[ruleType][ruleIndex] };

    rule.except = [...rule.except];
    rule.except[exceptIndex] = value;

    policy[ruleType] = [...policy[ruleType]];
    policy[ruleType][ruleIndex] = rule;

    newNetworkPolicies[policyIndex] = policy;
    onChange({ ...challenge, network_policies: newNetworkPolicies });
  };

  // 添加全屏编辑器的键盘事件监听
  useEffect(() => {
    const handleKeyDown = (event) => {
      if (fullscreenEditor.isOpen) {
        if (event.key === 'Escape') {
          closeFullscreenEditor();
        } else if (event.ctrlKey && event.key === 's') {
          event.preventDefault();
          saveFullscreenEditor();
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [fullscreenEditor.isOpen, fullscreenEditor.value]);

  // 常用样式类
  const inputBaseClass =
    'w-full h-10 bg-black/20 border border-neutral-300/30 rounded-md px-4 text-neutral-50 focus:outline-none focus:border-geek-400';
  const selectClass = 'select-custom select-custom-md';
  const textareaClass =
    'w-full h-20 bg-black/20 border border-neutral-300/30 rounded-md px-4 py-2 text-neutral-50 focus:outline-none focus:border-geek-400 resize-none';

  const podNoticeLines = [
    t('admin.challengeModal.podsNotice.flagFormat', { format: '`static{}`, `dynamic{}`, `uuid{}`' }),
    t('admin.challengeModal.podsNotice.flagPrefix'),
    t('admin.challengeModal.podsNotice.flagVolume'),
    '',
    t('admin.challengeModal.podsNotice.services'),
    '',
    t('admin.challengeModal.podsNotice.networks'),
    t('admin.challengeModal.podsNotice.networksWithIp'),
    t('admin.challengeModal.podsNotice.networksShared'),
    t('admin.challengeModal.podsNotice.portUnique'),
    '',
    t('admin.challengeModal.podsNotice.exampleIndent'),
  ];

  const yamlNoticeLines = [
    t('admin.challengeModal.yamlNotice.flagFormat', { format: '`static{}`, `dynamic{}`, `uuid{}`' }),
    t('admin.challengeModal.yamlNotice.flagPrefix'),
    t('admin.challengeModal.yamlNotice.flagVolume'),
    t('admin.challengeModal.yamlNotice.dockerParams'),
    t('admin.challengeModal.yamlNotice.containerUnique'),
    t('admin.challengeModal.yamlNotice.containerNetwork'),
  ];

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <motion.div
        className="w-full max-w-4xl bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* 标题栏 */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
          <h2 className="text-xl font-mono text-neutral-50">
            {mode === 'add'
              ? t('admin.challengeModal.title.add')
              : mode === 'edit'
                ? t('admin.challengeModal.title.edit')
                : t('admin.challengeModal.title.delete')}
          </h2>
          <Button
            variant="ghost"
            size="icon"
            className="!bg-transparent !text-neutral-400 hover:!text-neutral-300"
            onClick={onClose}
          >
            <IconX size={18} />
          </Button>
        </div>

        {/* 表单内容 */}
        <div className="p-6 max-h-[70vh] overflow-y-auto">
          {mode === 'delete' ? (
            <div className="text-center py-12">
              <div className="mb-8">
                <div className="mx-auto w-20 h-20 bg-red-400/20 rounded-full flex items-center justify-center mb-6">
                  <IconTrash size={40} className="text-red-400" />
                </div>
                <h3 className="text-2xl font-mono text-neutral-50 mb-4">{t('admin.challengeModal.delete.title')}</h3>
                <p className="text-neutral-400 font-mono text-lg mb-2">{t('admin.challengeModal.delete.prompt')}</p>
                <p className="text-red-400 font-mono text-sm">{t('admin.challengeModal.delete.warning')}</p>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              {/* 基本信息 */}
              <div className="mb-4">
                <h3 className="text-lg font-mono text-neutral-50 mb-3">{t('admin.challengeModal.sections.basic')}</h3>

                <div className="space-y-3">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div className="md:col-span-2">
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.name')}
                      </label>
                      <input
                        type="text"
                        value={challenge.name}
                        onChange={(e) => onChange({ ...challenge, name: e.target.value })}
                        className={inputBaseClass}
                        placeholder={t('admin.challengeModal.placeholders.name')}
                        required
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.category')}
                      </label>
                      <input
                        type="text"
                        value={challenge.category}
                        onChange={(e) => onChange({ ...challenge, category: e.target.value })}
                        className={selectClass}
                        placeholder={t('admin.challengeModal.placeholders.category')}
                        list="category-options"
                      />
                      <datalist id="category-options">
                        {categories.map((category) => (
                          <option key={category} value={category} />
                        ))}
                      </datalist>
                    </div>

                    <div>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">
                        {t('admin.challengeModal.labels.type')}
                      </label>
                      <select
                        value={challenge.type}
                        onChange={(e) => onChange({ ...challenge, type: e.target.value })}
                        className={selectClass}
                        required
                      >
                        <option value="static">{t('admin.challengeModal.types.static')}</option>
                        <option value="question">{t('admin.challengeModal.types.question')}</option>
                        <option value="dynamic">{t('admin.challengeModal.types.dynamic')}</option>
                        <option value="pods">{t('admin.challengeModal.types.pods')}</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.description')}
                    </label>
                    <textarea
                      value={challenge.description}
                      onChange={(e) => onChange({ ...challenge, description: e.target.value })}
                      className={textareaClass}
                      placeholder={t('admin.challengeModal.placeholders.description')}
                    />
                  </div>
                </div>
              </div>

              {/* 如果是动态类型题目，显示输入 generator */}
              {challenge.type === 'dynamic' && (
                <div className="border-t border-neutral-700 pt-4">
                  <h3 className="text-lg font-mono text-neutral-50 mb-3">
                    {t('admin.challengeModal.sections.generator')}
                  </h3>
                  <div>
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {t('admin.challengeModal.labels.generatorImage')}
                    </label>
                    <input
                      type="text"
                      value={challenge.generator_image}
                      onChange={(e) => onChange({ ...challenge, generator_image: e.target.value })}
                      className={inputBaseClass}
                      placeholder={t('admin.challengeModal.placeholders.generatorImage')}
                      required
                    />
                  </div>
                </div>
              )}

              {/* flag 设置 - 非pods, question类型才显示 */}
              {challenge.type !== 'pods' && challenge.type !== 'question' && (
                <div className="border-t border-neutral-700 pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.flags')}</h3>
                    <Button
                      variant="primary"
                      size="sm"
                      align="icon-left"
                      icon={<IconPlus size={16} />}
                      onClick={onAddFlag}
                    >
                      {t('admin.challengeModal.actions.addFlag')}
                    </Button>
                  </div>
                  <div className="space-y-3">
                    <label className="block text-sm font-mono text-neutral-400 mb-1">
                      {challenge.type === 'static' ? 'static{}' : 'dynamic{} / uuid{}'}
                    </label>
                    {challenge.flags.map((flag, index) => (
                      <div key={index} className="flex gap-2 items-center">
                        <input
                          type="text"
                          value={
                            typeof flag === 'string'
                              ? flag
                              : flag.value || (challenge.type === 'static' ? 'static{}' : 'uuid{}')
                          }
                          onChange={(e) => onFlagChange(index, e.target.value)}
                          className={inputBaseClass}
                          placeholder={t('admin.challengeModal.placeholders.flag', { index: index + 1 })}
                        />
                        {challenge.flags.length > 1 && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                            onClick={() => onRemoveFlag(index)}
                          >
                            <IconTrash size={18} />
                          </Button>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* 选项设置 - question类型才显示 */}
              {challenge.type === 'question' && (
                <div className="border-t border-neutral-700 pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.sections.options')}</h3>
                    <Button
                      variant="primary"
                      size="sm"
                      align="icon-left"
                      icon={<IconPlus size={16} />}
                      onClick={onAddOption}
                    >
                      {t('admin.challengeModal.actions.addOption')}
                    </Button>
                  </div>
                  <div className="space-y-3">
                    {(challenge.options || []).map((option, index) => (
                      <div key={option.rand_id || index} className="flex gap-2 items-center">
                        <div className="flex-1 flex gap-2 items-center">
                          <input
                            type="text"
                            value={option.content || ''}
                            onChange={(e) => onOptionChange(index, e.target.value)}
                            className={inputBaseClass}
                            placeholder={t('admin.challengeModal.placeholders.option', { index: index + 1 })}
                          />
                          <Button
                            variant={option.correct ? 'primary' : 'ghost'}
                            size="icon"
                            className={`!bg-transparent ${option.correct ? '!text-green-400 hover:!text-green-300' : '!text-red-400 hover:!text-red-300'}`}
                            onClick={() => onCorrectOptionChange(index)}
                            title={
                              option.correct
                                ? t('admin.challengeModal.actions.unsetCorrect')
                                : t('admin.challengeModal.actions.setCorrect')
                            }
                          >
                            {option.correct ? <IconCheck size={18} /> : <IconX size={18} />}
                          </Button>
                        </div>
                        {(challenge.options || []).length > 1 && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                            onClick={() => onRemoveOption(index)}
                          >
                            <IconTrash size={18} />
                          </Button>
                        )}
                      </div>
                    ))}
                    {(challenge.options || []).length === 0 && (
                      <div className="text-center py-4 text-neutral-500 font-mono text-sm">
                        {t('admin.challengeModal.empty.options')}
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* pods类型的flag说明 */}
              {challenge.type === 'pods' && (
                <div className="border-t border-neutral-700 pt-4">
                  <h3 className="text-lg font-mono text-neutral-50 mb-3">
                    {t('admin.challengeModal.sections.podsNotice')}
                  </h3>
                  <div className="p-3 bg-geek-400/10 border border-geek-400/20 rounded-md">
                    <p className="text-sm text-geek-400 font-mono">
                      {podNoticeLines.map((line, index) =>
                        line ? (
                          <span key={index}>
                            {line}
                            <br />
                          </span>
                        ) : (
                          <br key={index} />
                        )
                      )}
                    </p>
                  </div>
                </div>
              )}

              {/* 容器设置 */}
              {challenge.type === 'pods' && (
                <div className="border-t border-neutral-700 pt-4">
                  <div className="flex justify-between items-center mb-3">
                    <h3 className="text-lg font-mono text-neutral-50">
                      {t('admin.challengeModal.sections.containers')}
                    </h3>
                  </div>

                  <div className="space-y-6">
                    <div className="border border-neutral-700 rounded-md p-4 space-y-4">
                      {/* YAML 配置 */}
                      <div>
                        <div className="flex justify-between items-center mb-1">
                          <label className="block text-sm font-mono text-neutral-400">docker-compose.yaml</label>
                          <Button
                            variant="ghost"
                            size="sm"
                            align="icon-left"
                            icon={<IconMaximize size={14} />}
                            className="!bg-transparent !text-neutral-400 hover:!text-neutral-300 !text-xs"
                            onClick={openFullscreenEditor}
                          >
                            {t('admin.challengeModal.actions.fullscreen')}
                          </Button>
                        </div>
                        <div className="border border-neutral-300/30 rounded-md overflow-hidden">
                          <Suspense
                            fallback={
                              <div className="flex items-center justify-center h-[200px] text-neutral-400 font-mono text-sm">
                                Loading editor…
                              </div>
                            }
                          >
                            <Editor
                              value={challenge.docker_compose || defaultYAML}
                              onChange={updateDockerCompose}
                              language="yaml"
                              options={{
                                readOnly: false,
                                minimap: { enabled: false },
                                scrollBeyondLastLine: false,
                                scrollbar: {
                                  vertical: 'auto',
                                  horizontal: 'auto',
                                },
                                lineNumbers: 'on',
                                folding: true,
                                wordWrap: 'on',
                                fontSize: 14,
                                fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
                                tabSize: 2,
                                insertSpaces: true,
                                renderLineHighlight: 'line',
                              }}
                              height={`${19 * (challenge.docker_compose || defaultYAML).split('\n').length}px`}
                              theme="vs-dark"
                            />
                          </Suspense>
                        </div>
                        <div className="mt-1 text-xs text-neutral-500 font-mono">
                          {yamlNoticeLines.map((line, index) => (
                            <span key={index}>
                              {line}
                              <br />
                            </span>
                          ))}
                        </div>
                      </div>

                      {/* 网络策略 */}
                      <div className="mt-4">
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
                            onClick={addNetworkPolicy}
                          >
                            {t('admin.challengeModal.actions.addPolicy')}
                          </Button>
                        </div>

                        <div className="space-y-4">
                          {(challenge.network_policies || []).map((policy, policyIndex) => (
                            <div key={policyIndex} className="border border-neutral-700 rounded-md p-3 bg-black/20">
                              <div className="flex justify-between items-center mb-3">
                                <h5 className="text-sm font-mono text-neutral-200">
                                  {t('admin.challengeModal.labels.policy', { index: policyIndex + 1 })}
                                </h5>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  className="!bg-transparent !text-red-400 hover:!text-red-300"
                                  onClick={() => removeNetworkPolicy(policyIndex)}
                                >
                                  <IconTrash size={16} />
                                </Button>
                              </div>

                              {/* From 策略 */}
                              <div className="mb-3">
                                <div className="flex justify-between items-center mb-2">
                                  <label className="text-xs font-mono text-neutral-400">
                                    {t('admin.challengeModal.labels.allowInbound')}
                                  </label>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    align="icon-left"
                                    icon={<IconPlus size={12} />}
                                    className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                    onClick={() => addPolicyRule(policyIndex, 'from')}
                                  >
                                    {t('admin.challengeModal.actions.addRule')}
                                  </Button>
                                </div>

                                <div className="space-y-2">
                                  {(policy.from || []).map((fromRule, fromIndex) => (
                                    <div key={fromIndex} className="border border-neutral-700/50 p-2 rounded">
                                      <div className="flex justify-between items-center mb-2">
                                        <label className="text-xs font-mono text-neutral-400">CIDR</label>
                                        {policy.from.length > 1 && (
                                          <Button
                                            variant="ghost"
                                            size="icon"
                                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                                            onClick={() => removePolicyRule(policyIndex, 'from', fromIndex)}
                                          >
                                            <IconTrash size={12} />
                                          </Button>
                                        )}
                                      </div>

                                      <input
                                        type="text"
                                        value={fromRule.cidr}
                                        onChange={(e) =>
                                          updatePolicyCidr(policyIndex, 'from', fromIndex, e.target.value)
                                        }
                                        className="w-full h-8 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                        placeholder="192.168.1.0/24"
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
                                            onClick={() => addPolicyExcept(policyIndex, 'from', fromIndex)}
                                          >
                                            {t('admin.challengeModal.actions.addExclude')}
                                          </Button>
                                        </div>
                                        <div className="space-y-1">
                                          {(fromRule.except || []).map((exceptItem, exceptIndex) => (
                                            <div key={exceptIndex} className="flex gap-1 items-center">
                                              <input
                                                type="text"
                                                value={exceptItem}
                                                onChange={(e) =>
                                                  updatePolicyExcept(
                                                    policyIndex,
                                                    'from',
                                                    fromIndex,
                                                    exceptIndex,
                                                    e.target.value
                                                  )
                                                }
                                                className="flex-1 h-6 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                                placeholder="192.168.1.1"
                                              />
                                              <Button
                                                variant="ghost"
                                                size="icon"
                                                className="!bg-transparent !text-red-400 hover:!text-red-300 !w-6 !h-6"
                                                onClick={() =>
                                                  removePolicyExcept(policyIndex, 'from', fromIndex, exceptIndex)
                                                }
                                              >
                                                <IconTrash size={10} />
                                              </Button>
                                            </div>
                                          ))}
                                        </div>
                                      </div>
                                    </div>
                                  ))}
                                </div>
                              </div>

                              {/* To 策略 */}
                              <div>
                                <div className="flex justify-between items-center mb-2">
                                  <label className="text-xs font-mono text-neutral-400">
                                    {t('admin.challengeModal.labels.allowOutbound')}
                                  </label>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    align="icon-left"
                                    icon={<IconPlus size={12} />}
                                    className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs"
                                    onClick={() => addPolicyRule(policyIndex, 'to')}
                                  >
                                    {t('admin.challengeModal.actions.addRule')}
                                  </Button>
                                </div>

                                <div className="space-y-2">
                                  {(policy.to || []).map((toRule, toIndex) => (
                                    <div key={toIndex} className="border border-neutral-700/50 p-2 rounded">
                                      <div className="flex justify-between items-center mb-2">
                                        <label className="text-xs font-mono text-neutral-400">CIDR</label>
                                        {policy.to.length > 1 && (
                                          <Button
                                            variant="ghost"
                                            size="icon"
                                            className="!bg-transparent !text-red-400 hover:!text-red-300"
                                            onClick={() => removePolicyRule(policyIndex, 'to', toIndex)}
                                          >
                                            <IconTrash size={12} />
                                          </Button>
                                        )}
                                      </div>

                                      <input
                                        type="text"
                                        value={toRule.cidr}
                                        onChange={(e) => updatePolicyCidr(policyIndex, 'to', toIndex, e.target.value)}
                                        className="w-full h-8 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                        placeholder="0.0.0.0/0"
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
                                            onClick={() => addPolicyExcept(policyIndex, 'to', toIndex)}
                                          >
                                            {t('admin.challengeModal.actions.addExclude')}
                                          </Button>
                                        </div>
                                        <div className="space-y-1">
                                          {(toRule.except || []).map((exceptItem, exceptIndex) => (
                                            <div key={exceptIndex} className="flex gap-1 items-center">
                                              <input
                                                type="text"
                                                value={exceptItem}
                                                onChange={(e) =>
                                                  updatePolicyExcept(
                                                    policyIndex,
                                                    'to',
                                                    toIndex,
                                                    exceptIndex,
                                                    e.target.value
                                                  )
                                                }
                                                className="flex-1 h-6 bg-black/30 border border-neutral-700 rounded px-2 text-neutral-50 text-xs"
                                                placeholder="10.0.0.0/8"
                                              />
                                              <Button
                                                variant="ghost"
                                                size="icon"
                                                className="!bg-transparent !text-red-400 hover:!text-red-300 !w-6 !h-6"
                                                onClick={() =>
                                                  removePolicyExcept(policyIndex, 'to', toIndex, exceptIndex)
                                                }
                                              >
                                                <IconTrash size={10} />
                                              </Button>
                                            </div>
                                          ))}
                                        </div>
                                      </div>
                                    </div>
                                  ))}
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* 底部按钮 */}
        <div className="flex justify-end gap-4 p-4 border-t border-neutral-700">
          <Button variant="ghost" size="sm" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button variant="primary" size="sm" onClick={() => onSubmit(challenge)}>
            {mode === 'add'
              ? t('admin.challengeModal.actions.add')
              : mode === 'edit'
                ? t('common.saveChanges')
                : t('admin.challengeModal.actions.confirmDelete')}
          </Button>
        </div>
      </motion.div>

      {/* 全屏YAML编辑器 */}
      {fullscreenEditor.isOpen && (
        <div className="fixed inset-0 z-[100] bg-black/80 backdrop-blur-sm flex items-center justify-center">
          <motion.div
            className="w-[95vw] h-[95vh] bg-neutral-900 border border-neutral-700 rounded-lg overflow-hidden"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.2 }}
          >
            <div className="flex justify-between items-center p-4 border-b border-neutral-700">
              <h3 className="text-lg font-mono text-neutral-50">{t('admin.challengeModal.fullscreen.title')}</h3>
              <div className="flex gap-2">
                <Button variant="primary" size="sm" onClick={saveFullscreenEditor}>
                  <IconDeviceFloppy size={16} className="mr-2" />
                  {t('admin.challengeModal.fullscreen.save')}
                </Button>
                <Button variant="ghost" size="sm" onClick={closeFullscreenEditor}>
                  <IconMinimize size={16} className="mr-2" />
                  {t('admin.challengeModal.fullscreen.exit')}
                </Button>
              </div>
            </div>
            <div className="flex-1 h-full">
              <Suspense
                fallback={
                  <div className="flex items-center justify-center h-full text-neutral-400 font-mono text-sm">
                    Loading editor…
                  </div>
                }
              >
                <Editor
                  value={fullscreenEditor.value}
                  onChange={updateFullscreenValue}
                  language="yaml"
                  options={{
                    readOnly: false,
                    minimap: { enabled: true },
                    scrollBeyondLastLine: false,
                    lineNumbers: 'on',
                    folding: true,
                    wordWrap: 'on',
                    fontSize: 16,
                    fontFamily: 'Consolas, "Courier New", monospace',
                    tabSize: 2,
                    insertSpaces: true,
                    renderLineHighlight: 'line',
                  }}
                  height="calc(95vh - 70px)"
                  theme="vs-dark"
                />
              </Suspense>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}

export default AdminChallengeModal;
