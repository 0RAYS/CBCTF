import TechStack from '../components/features/Support/TechStack';
import { useTranslation } from 'react-i18next';

function TechStackPage() {
  const { t } = useTranslation();
  const techStackConfig = {
    technologies: [
      {
        category: t('support.techStack.categories.frontend'),
        items: [
          { name: 'React', icon: '⚛️', description: t('support.techStack.items.react') },
          { name: 'TailwindCSS', icon: '🎨', description: t('support.techStack.items.tailwind') },
          { name: 'Motion', icon: '✨', description: t('support.techStack.items.motion') },
        ],
      },
      {
        category: t('support.techStack.categories.backend'),
        items: [
          { name: 'Go', icon: '🔵', description: t('support.techStack.items.go') },
          { name: 'MySQL', icon: '🎲', description: t('support.techStack.items.mysql') },
          { name: 'Redis', icon: '⚡', description: t('support.techStack.items.redis') },
        ],
      },
      {
        category: t('support.techStack.categories.devops'),
        items: [
          { name: 'Docker', icon: '🐳', description: t('support.techStack.items.docker') },
          { name: 'K8s', icon: '⚓️', description: t('support.techStack.items.k8s') },
          { name: 'Jenkins', icon: '🔄', description: t('support.techStack.items.jenkins') },
        ],
      },
    ],
    developers: [
      {
        name: '1manity',
        role: t('support.techStack.roles.frontendLead'),
        picture: 'https://1manity.top/10.jpg',
        github: 'https://github.com/1manity',
      },
      {
        name: 'JBNRZ',
        role: t('support.techStack.roles.backendLead'),
        picture: 'https://q.qlogo.cn/headimg_dl?dst_uin=3537659915&spec=640&img_type=jpg',
        github: 'https://github.com/JBNRZ',
      },
    ],
  };

  return <TechStack {...techStackConfig} />;
}

export default TechStackPage;
