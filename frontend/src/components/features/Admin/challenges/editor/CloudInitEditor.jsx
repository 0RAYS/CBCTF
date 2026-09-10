import { useTranslation } from 'react-i18next';
import { createGuideCloudConfig, createGuideUser, createGuideGroup, createGuideWriteFile } from './guideModel.js';
import { GuideListHeader, GuideInput, GuideBoolean, GuideStringList, IconButton } from './GuideFields.jsx';

export default function CloudInitEditor({ value, onChange }) {
  const { t } = useTranslation();
  const ct = (key) => t(`admin.challengeModal.composeGuide.${key}`);
  const config = { ...createGuideCloudConfig(), ...value };
  const update = (field, index, itemField, nextValue) =>
    onChange({
      ...config,
      [field]: config[field].map((item, i) => (i === index ? { ...item, [itemField]: nextValue } : item)),
    });
  const add = (field, item) => onChange({ ...config, [field]: [...config[field], item] });
  const remove = (field, index) => onChange({ ...config, [field]: config[field].filter((_, i) => i !== index) });

  return (
    <>
      <GuideListHeader
        title={ct('sections.users')}
        addLabel={ct('actions.add')}
        onAdd={() => add('users', createGuideUser())}
      />
      {config.users.map((user, userIndex) => (
        <div key={userIndex} className="space-y-2 rounded border border-neutral-700/60 p-3">
          <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_1fr_32px] gap-2">
            {['name', 'gecos', 'shell'].map((field) => (
              <GuideInput
                key={field}
                label={ct(`fields.${field}`)}
                value={user[field]}
                placeholder={ct(`placeholders.${field}`)}
                onChange={(next) => update('users', userIndex, field, next)}
              />
            ))}
            <IconButton onClick={() => remove('users', userIndex)} />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
            {['homedir', 'passwd', 'plainTextPasswd'].map((field) => {
              const label = field === 'homedir' ? 'homeDir' : field;
              return (
                <GuideInput
                  key={field}
                  label={ct(`fields.${label}`)}
                  value={user[field]}
                  placeholder={ct(`placeholders.${label}`)}
                  onChange={(next) => update('users', userIndex, field, next)}
                />
              );
            })}
            <div className="grid grid-cols-3 gap-2">
              {['lockPasswd', 'noCreateHome', 'system'].map((field) => (
                <GuideBoolean
                  key={field}
                  label={ct(`fields.${field}`)}
                  value={user[field]}
                  onChange={(next) => update('users', userIndex, field, next)}
                />
              ))}
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
            {['groups', 'sudo', 'sshAuthorizedKeys'].map((field) => (
              <GuideStringList
                key={field}
                label={ct(`fields.${field}`)}
                placeholder={ct(`placeholders.${field}`)}
                values={user[field]}
                addLabel={ct('actions.add')}
                onChange={(next) => update('users', userIndex, field, next)}
              />
            ))}
          </div>
        </div>
      ))}

      <GuideListHeader
        title={ct('sections.groups')}
        addLabel={ct('actions.add')}
        onAdd={() => add('groups', createGuideGroup())}
      />
      {config.groups.map((group, groupIndex) => (
        <div key={groupIndex} className="space-y-2 rounded border border-neutral-700/60 p-3">
          <div className="grid grid-cols-[1fr_32px] gap-2">
            <GuideInput
              label={ct('fields.name')}
              value={group.name}
              placeholder={ct('placeholders.name')}
              onChange={(next) => update('groups', groupIndex, 'name', next)}
            />
            <IconButton onClick={() => remove('groups', groupIndex)} />
          </div>
          <GuideStringList
            label={ct('fields.members')}
            placeholder={ct('placeholders.members')}
            values={group.members}
            addLabel={ct('actions.add')}
            onChange={(next) => update('groups', groupIndex, 'members', next)}
          />
        </div>
      ))}

      <GuideStringList
        label={ct('fields.sshAuthorizedKeys')}
        placeholder={ct('placeholders.sshAuthorizedKeys')}
        values={config.sshAuthorizedKeys}
        addLabel={ct('actions.add')}
        onChange={(next) => onChange({ ...config, sshAuthorizedKeys: next })}
      />

      <GuideListHeader
        title={ct('sections.writeFiles')}
        addLabel={ct('actions.add')}
        onAdd={() => add('writeFiles', createGuideWriteFile())}
      />
      {config.writeFiles.map((file, fileIndex) => (
        <div key={fileIndex} className="space-y-2 rounded border border-neutral-700/60 p-3">
          <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_1fr_32px] gap-2">
            {['path', 'owner', 'permissions'].map((field) => (
              <GuideInput
                key={field}
                label={ct(`fields.${field}`)}
                value={file[field]}
                placeholder={ct(`placeholders.${field}`)}
                onChange={(next) => update('writeFiles', fileIndex, field, next)}
              />
            ))}
            <IconButton onClick={() => remove('writeFiles', fileIndex)} />
          </div>
          <div className="grid grid-cols-1 md:grid-cols-[1fr_120px_120px] gap-2">
            <GuideInput
              label={ct('fields.encoding')}
              value={file.encoding}
              placeholder={ct('placeholders.encoding')}
              onChange={(next) => update('writeFiles', fileIndex, 'encoding', next)}
            />
            {['append', 'defer'].map((field) => (
              <GuideBoolean
                key={field}
                label={ct(`fields.${field}`)}
                value={file[field]}
                onChange={(next) => update('writeFiles', fileIndex, field, next)}
              />
            ))}
          </div>
          <GuideInput
            multiline
            label={ct('fields.content')}
            value={file.content}
            placeholder={ct('placeholders.content')}
            onChange={(next) => update('writeFiles', fileIndex, 'content', next)}
          />
        </div>
      ))}
    </>
  );
}
