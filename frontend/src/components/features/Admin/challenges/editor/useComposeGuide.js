import { useState, useEffect, useRef } from 'react';
import { createGuideConfig, emptyComposeYaml, normalizeGuideConfigForMode } from './guideModel.js';
import { buildGuidedComposeYaml } from './composeSerializer.js';
import { parseComposeYamlToGuideConfig } from './composeParser.js';
import { validateGuidedCompose, validateRawCompose } from './composeValidation.js';

// Own only the bidirectional text/guide boundary. Domain edits live in their editors.
export default function useComposeGuide({ isOpen, challenge, onChange, t }) {
  const [guideConfig, setGuideConfig] = useState(createGuideConfig);
  const guideGeneratedComposeRef = useRef('');
  const guideValidation = validateGuidedCompose(guideConfig, t);
  const rawValidation = validateRawCompose(challenge.docker_compose || emptyComposeYaml, t);

  const syncGuideConfig = (nextConfig) => {
    const normalizedConfig = normalizeGuideConfigForMode(nextConfig);
    const dockerCompose = buildGuidedComposeYaml(normalizedConfig);
    guideGeneratedComposeRef.current = dockerCompose;
    setGuideConfig(normalizedConfig);
    onChange({ ...challenge, docker_compose: dockerCompose });
  };

  useEffect(() => {
    if (!isOpen || challenge.type !== 'pods') {
      guideGeneratedComposeRef.current = '';
      setGuideConfig(createGuideConfig());
      return;
    }
    const dockerCompose = challenge.docker_compose || '';
    if (!dockerCompose) {
      guideGeneratedComposeRef.current = '';
      setGuideConfig(createGuideConfig());
      return;
    }
    if (guideGeneratedComposeRef.current === dockerCompose) return;
    guideGeneratedComposeRef.current = '';
    const parsed = parseComposeYamlToGuideConfig(dockerCompose, t);
    if (parsed.ok) setGuideConfig(normalizeGuideConfigForMode(parsed.config));
  }, [isOpen, challenge.id, challenge.type, challenge.docker_compose, t]);

  const updateDockerCompose = (value) => {
    const finalValue = value || '';
    guideGeneratedComposeRef.current = '';
    const parsed = parseComposeYamlToGuideConfig(finalValue, t);
    if (!parsed.ok) {
      onChange({ ...challenge, docker_compose: finalValue });
      return;
    }
    setGuideConfig(normalizeGuideConfigForMode(parsed.config));
    onChange({ ...challenge, docker_compose: finalValue });
  };

  return { guideConfig, guideValidation, rawValidation, syncGuideConfig, updateDockerCompose };
}
